package server

import "sync"

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
