package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/inconshreveable/log15.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {

}

const reposTTL = time.Hour * 24 * 45

var reposRemoved = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "src",
	Subsystem: "gitserver",
	Name:      "repos_removed",
	Help:      "number of repos removed during cleanup",
})

var reposRecloned = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "src",
	Subsystem: "gitserver",
	Name:      "repos_recloned",
	Help:      "number of repos removed and recloned due to age",
})

// cleanupRepos walks the repos directory and performs maintenance tasks:
//
// 1. Remove corrupt repos.
// 2. Remove stale lock files.
// 3. Remove inactive repos on sourcegraph.com
// 4. Reclone repos after a while. (simulate git gc)
func (s *Server) cleanupRepos() {
	//bCtx, bCancel := s.serverContext()
	//defer bCancel()
	//
	//maybeRemoveCorrupt := func(gitDir string) (done bool, err error) {
	//	// We treat repositories missing HEAD to be corrupt. Both our cloning
	//	// and fetching ensure there is  a HEAD file.
	//	_, err := os.Stat(filepath.Join(gitDir, "HEAD"))
	//	if !os.IsNotExist(err) {
	//		return false, err
	//	}
	//
	//	log15.Info("removing corrupt repo", "repo", gitDir)
	//	if err := s.removeReposDirectory(gitDir); err != nil {
	//		return true, err
	//	}
	//	reposRemoved.Inc()
	//	return true, nil
	//}
	//
	//ensureGitAttributes := func(gitDir string) (done bool, err error) {
	//	return false, setGitAttributes(gitDir)
	//}

}

// removeReposDirectory atomically removes a directory from s.ReposDir
//
// It first moves the directory to a temporary location to avoid leaving
// partial state in the event of server restart or concurrent notifications to
// the directory.
//
// Additionally it removes parent empty directories up until s.ReposDir.
func (s *Server) removeReposDirectory(dir string) error {
	// Rename out of the location so we can atomically stop using the repo.
	tmp, err := s.tempDir("delete-repo")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	if err := os.Rename(dir, filepath.Join(tmp, "repo")); err != nil {
		return err
	}

	// Everything after this point is just cleanup, so any error that occurs
	// should not be returned, just logged.

	// Cleanup empty parent directories. We just attempt to remove and if we
	// have a failure we assume it's due to the directory having other
	// children. If we checked first we could race with someone else adding a
	// new clone.
	rootInfo, err := os.Stat(s.ReposDir)
	if err != nil {
		log15.Warn("Failed to stat ReposDir", "error", err)
		return nil
	}
	current := dir
	for {
		parent := filepath.Dir(current)
		if parent == current {
			// This shouldn't happen, but protecting against escaping
			// reposDir
			break
		}
		current = parent
		info, err := os.Stat(current)
		if os.IsNotExist(err) {
			// Someone else beat us to it.
			break
		}
		if err != nil {
			log15.Warn("failed to stat parent directory", "dir", current, "error", err)
			return nil
		}
		if os.SameFile(rootInfo, info) {
			// Stop, we are at the parent.
			break
		}

		if err := os.Remove(current); err != nil {
			// Stop, we assume remove failed due to current not being empty.
			break
		}
	}

	// Delete the atomically renamed dir. We do this last since if it fails we
	// will rely on a janitor job to clean up for us.
	if err := os.RemoveAll(filepath.Join(tmp, "repo")); err != nil {
		log15.Warn("failed to cleanup after removing dir", "dir", dir, "error", err)
	}
	return nil
}

// SetupAndClearTmp sets up the tempdir for ReposDir as well as clearing it
// out. It returns the temporary directory location.
func (s *Server) SetupAndClearTmp() (string, error) {
	// Additionally we create directories with the prefix .tmp-old which are
	// asynchronously removed. We do not remove in place since it may be a
	// slow operation to block on. Our tmp dir will be ${s.ReposDir}/.tmp
	dir := filepath.Join(s.ReposDir, tempDirName) // .tmp
	oldPrefix := tempDirName + "-old"
	if _, err := os.Stat(dir); err == nil {
		// Rename the current tmp file so we can asynchronously remove it. Use
		// a consistent pattern so if we get interrupted, so we can clean it
		// another name.
		oldTmp, err := ioutil.TempDir(s.ReposDir, oldPrefix)
		if err != nil {
			return "", err
		}
		// oldTmp dir exists, so we need to use a child of oldTmp as the
		// rename target.
		if err := os.Rename(dir, filepath.Join(oldTmp, tempDirName)); err != nil {
			return "", err
		}
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	// Asynchronously remove old temporary directories
	files, err := ioutil.ReadDir(s.ReposDir)
	if err != nil {
		log15.Error("failed to do tmp cleanup", "error", err)
	} else {
		for _, f := range files {
			// Remove older .tmp directories as well as our older tmp-
			// directories we would place into ReposDir.
			if !strings.HasPrefix(f.Name(), oldPrefix) && !strings.HasPrefix(f.Name(), "tmp-") {
				continue
			}
			go func(path string) {
				if err := os.RemoveAll(path); err != nil {
					log15.Error("cleanup: failed to remove old temporary directory", "path", "error", err)
				}
			}(filepath.Join(s.ReposDir, f.Name()))
		}
	}
	return dir, nil
}
