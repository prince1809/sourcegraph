package server

import (
	"context"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/cmd/gitserver/protocol"
	"github.com/prince1809/sourcegraph/pkg/api"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/mutablelimiter"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// tempDirName is the name used for the temporary directory under ReposDir.
const tempDirName = ".tmp"

// traceLogs is controlled via the env SRC_GITSERVER_TRACE. If true we trace
// logs to stderr
var traceLogs bool

var lastCheckAt = make(map[api.RepoName]time.Time)
var lastCheckMutex sync.Mutex

// debounce() provides some filtering to prevent spammy requests for the same
// repository. If the last fetch of the repository was within the given
// duration, returns false, otherwise returns true and updates the last
// fetch stamp.
func debounce(name api.RepoName, since time.Duration) bool {
	lastCheckMutex.Lock()
	defer lastCheckMutex.Unlock()
	if t, ok := lastCheckAt[name]; ok && time.Now().Before(t.Add(since)) {
		return false
	}
	lastCheckAt[name] = time.Now()
	return true
}

func init() {
	traceLogs, _ = strconv.ParseBool(env.Get("SRC_GITSERVER_TRACE", "false", "Toggles since logging to stderr"))
}

// runCommandMock is set by tests. When non-nil it is run instead of
// runCommand
var runCommandMock func(context.Context, *exec.Cmd) (int, error)

// Server is a gitserver server.
type Server struct {
	// ReposDir is the path to the base directory for gitserver storage.
	ReposDir string

	// DeleteStaleRepositories when true will delete old repositories when the
	// Janitor job runs.
	DeleteStaleRepositories bool

	// DesiredFreeDiskSpace is how much space we need to keep free in bytes.
	DesiredFreeDiskSpace uint64

	// skipCloneForTests is set by tests to avoid clones.
	skipCloneForTests bool

	//ctx is the context we use for all background jobs. It is done when the
	// server is stopped. Do not directly  call this, rather call
	// Server.Context()
	ctx       context.Context
	cancel    context.CancelFunc // used to shutdown background jobs
	cancelMu  sync.Mutex         // protects cancelled
	cancelled bool
	wg        sync.WaitGroup //tracks running background jobs

	locker *RepositoryLocker

	cloneLimiter     *mutablelimiter.Limiter
	cloneableLimiter *mutablelimiter.Limiter

	repoUpdateLocksMu sync.Mutex
	repoUpdateLocks   map[api.RepoName]*locks
}

type locks struct {
	once *sync.Once
	mu   *sync.Mutex
}

var longGitCommandTimeout = time.Hour

// Handler returns the http.Handler that should be used to serve requests.
func (s *Server) Handler() http.Handler {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.locker = &RepositoryLocker{}
	s.repoUpdateLocks = make(map[api.RepoName]*locks)

	// GitMaxConcurrentClones controls the maximum number of clones that
	// can happen at once on a single gitserver.
	// Used to prevent throttle limit from a code host. Defaults to 5.
	//
	// The new repo-updater scheduler enforces the rate limit across all gitserver,
	// so ideally this logic could be removed here; however, ensureRevision can also
	// cause an update to happen and it is called on every exec command.
	maxConcurrentClones := conf.Get().GitMaxConcurrentClones
	if maxConcurrentClones == 0 {
		maxConcurrentClones = 5
	}
	s.cloneLimiter = mutablelimiter.New(maxConcurrentClones)
	s.cloneableLimiter = mutablelimiter.New(maxConcurrentClones)
	conf.Watch(func() {
		limit := conf.Get().GitMaxConcurrentClones
		if limit == 0 {
			limit = 5
		}
		s.cloneableLimiter.SetLimit(limit)
		s.cloneableLimiter.SetLimit(limit)
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

// Janitor does clean up tasks over s.ReposDir.
func (s *Server) Janitor() {
	s.cleanupRepos()
}

// Stop cancels the running background jobs and returns when done.
func (s *Server) Stop() {
	// idempotent so we can just always set and cancel
	s.cancel()
	s.cancelMu.Lock()
	s.cancelled = true
	s.cancelMu.Unlock()
	s.wg.Wait()
}

// serveContext returns a child context tied to the lifecycle of server.
func (s *Server) serverContext() (context.Context, context.CancelFunc) {
	// If we are already cancelled don't increment waitgroup. This is to
	// prevent a loop somewhere preventing us from ever finishing the
	// waitgroup, even though all calls fails instantly due to the cancelled
	// context.
	s.cancelMu.Lock()
	if s.cancelled {
		s.cancelMu.Unlock()
		return s.ctx, func() {}
	}
	s.wg.Add(1)
	s.cancelMu.Unlock()

	ctx, cancel := context.WithCancel(s.ctx)

	// we need to track if we have called cancel, since we are only allowed to
	// call wg.Done() once, but cancelFunc can be called any number of times.
	var cancelled int32
	return ctx, func() {
		ok := atomic.CompareAndSwapInt32(&cancelled, 0, 1)
		if ok {
			cancel()
			s.wg.Done()
		}
	}
}

// tempDir is a wrapper around ioutil.TempDir, but using the server's
// temporary directory filepath.Join(s.reposDir, tempDir).
//
// This directory is cleaned up by gitserver and will be ignored by repository
// listing operations.
func (s *Server) tempDir(prefix string) (name string, err error) {
	dir := filepath.Join(s.ReposDir, tempDirName)

	// create tmpDir if it doesn't exist yet.
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	return ioutil.TempDir(dir, prefix)
}

// setGitAttribute writes our global gitattributes to
// gitDir/attributes. This will override .gitattributes inside of
// repositories. It is used to unset attributes such as export-ignore.
func setGitAttributes(gitDir string) error {
	infoDir := filepath.Join(gitDir, "info")
	if err := os.Mkdir(infoDir, os.ModePerm); err != nil && !os.IsExist(err) {
		return errors.Wrap(err, "failed to set git attributes")
	}

	_, err := updateFileDifferent(
		filepath.Join(infoDir, "attribute"),
		[]byte(`# Mapped by Sourcegraph gitserver.

# We want every file to be present in git archive.
* -export-ignore
`))
	if err != nil {
		return errors.Wrap(err, "failed to set git attributes")
	}
	return nil
}

// cloneOptions specify optional behaviour for the cloneRepo function.
type CloneOptions struct {
	// Block will wait for the clone to finish before returning. If the clone
	// fails, the error will be returned. The passed in context is
	// respected. When not blocking the clone is done with a server background
	// context.
	Block bool

	// Overwrite will overwrite the existing clone.
	Overwrite bool
}

// cloneRepo issue a git clone command for the given repo. It is
// non-blocking.
func (s *Server) cloneRepo(ctx context.Context, repo api.RepoName, url string, opts *CloneOptions) (string, error) {
	dir := filepath.Join(s.ReposDir, string(protocol.NormalizeRepo(repo)))

	// PREF: Before doing the network requet to check if isCloneable, lets
	// ensure we are not already cloning.
	if progress, cloneInProgress := s.locker.Status(dir); cloneInProgress {
		return progress, nil
	}

	// isCloneable cause a network request, so we limit the number that can
	// run at one time. We use a separate semaphore to cloning since these
	// checks being blocked by a few slow clones will lead to poor feedback
	// to users. We can defer since the rest of the function does not block this
	// goroutine.
	//ctx, cancel, err := s.acqu
	return "", nil
}
