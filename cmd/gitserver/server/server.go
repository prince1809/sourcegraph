package server

import (
	"context"
	"sync"
)

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
}
