package server

import (
	"path/filepath"
	"sync"
)

// RepositoryLocker provides locks for doing operations to a repository
// directory. When a repository locked, only the owner of the lock is
// allowed to run command against it.
//
// Repositories are identified by the absolute path to their directories. Note
// the directories are the parent of the $GIT_DIR (ie excluding the /.git
// suffix). All operations affect the %GIT_DIR, but for legacy reasons we
// identify repositories by the parent.

type RepositoryLocker struct {
	// mu protects status
	mu sync.Mutex
	// status tracks directories that are looked. The value is the status. If
	// a directory is in status, the directory is locked.
	status map[string]string
}

// Status returns the status of the locked directory dir. If dir is not
// locked, then locked is false.
func (rl *RepositoryLocker) Status(dir string) (status string, locked bool) {
	dir = rl.normalize(dir)

	rl.mu.Lock()
	status, locked = rl.status[dir]
	rl.mu.Unlock()
	return
}

// normalize cleans dir and ensures dir is not pointing to the GIT_DIR, but
// rather the parent. ie it will translate
// /data/repos/example.com/foo/bar/.git to /data/repos/example.com/foo/bar
func (rl *RepositoryLocker) normalize(dir string) string {
	dir = filepath.Clean(dir)

	// Use parent if we are passed a $GIT_DIR
	if name := filepath.Base(dir); name == ".git" {
		return filepath.Dir(dir)
	}

	return dir
}

// RepositoryLock is returned by RepositoryLocker.TryAcquire. It allows
// updating the status of a directory lock, as well as releasing the lock.
type RepositoryLock struct {
	locker *RepositoryLocker
	dir    string

	// done is protected by locker.mu
	done bool
}


