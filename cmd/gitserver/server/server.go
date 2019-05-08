package server

import (
	"context"
	"github.com/prince1809/sourcegraph/pkg/api"
	"github.com/prince1809/sourcegraph/pkg/mutablelimiter"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"net/http"
	"sync"
	"time"
)

// tempDirName is the name used for the temporary directory under ReposDir.
const tempDirName = ".tmp"

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
	return mux
}
