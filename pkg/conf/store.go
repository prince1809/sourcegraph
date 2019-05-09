package conf

import (
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"runtime/debug"
	"sync"
	"time"
)

// Store manages the in-memory storage, access,
// and updating of the site configuration in threadsafe manner.
type Store struct {
	configMu  sync.RWMutex
	lastValid *Unified
	mock      *Unified

	rawMu sync.RWMutex
	raw   conftypes.RawUnified

	ready chan struct{}
	once  sync.Once
}

// NewStore returns a new configuration store.
func NewStore() *Store {
	return &Store{
		ready: make(chan struct{}),
	}
}

// LastValid returns the last valid site configuration that this
// store was updated with.
func (s *Store) LastValid() *Unified {
	s.WaitUntilInitialized()

	s.configMu.RLock()
	defer s.configMu.RUnlock()

	if s.mock != nil {
		return s.mock
	}

	return s.lastValid
}

// Raw returns the last raw configuration that this store was updated with.
func (s *Store) Raw() conftypes.RawUnified {
	s.WaitUntilInitialized()

	s.rawMu.RLock()
	defer s.rawMu.RUnlock()
	return s.raw
}

type UpdateResult struct {
	Changed bool
	Old     *Unified
	New     *Unified
}

func (s *Store) MaybeUpdate(rawConfig conftypes.RawUnified) (UpdateResult, error) {
	s.rawMu.Lock()
	defer s.rawMu.Unlock()

	s.configMu.Lock()
	defer s.configMu.Unlock()

	result := UpdateResult{
		Changed: false,
		Old:     s.lastValid,
		New:     s.lastValid,
	}

	if rawConfig.Critical == "" {
		return result, errors.New("invalid critical configuration (empty string)")
	}

	if rawConfig.Site == "" {
		return result, errors.New("invalid site configuration (empty string)")
	}

	if s.raw.Equal(rawConfig) {
		return result, nil
	}
	s.raw = rawConfig

	newConfig, err := ParseConfig(rawConfig)
	if err != nil {
		return result, errors.Wrap(err, "when parsing rawConfig during update")
	}

	result.Changed = true
	result.New = newConfig
	s.lastValid = newConfig

	s.initialize()

	return result, nil
}

// WaitUntilInitialized blocks and only returns to the caller once the store
// has initialized with a syntactically valid configuration file (via MaybeUpdate() or Mock()).
func (s *Store) WaitUntilInitialized() {
	mode := getMode()
	if mode == modeServer {
		select {
		case <-configurationServerFrontendOnlyInitialized:
		// We assume that we're in an unrecoverable deadlock if frontend hasn't
		// started its configuration server after 30 seconds.
		case <-time.After(30 * time.Second):
			// The running goroutine is not necessarily the cause of the
			// deadlock, so ask Go to dump all goroutine stack traces.
			debug.SetTraceback("all")
			panic("deadlock detected: you have called conf.Get or conf.Watch before the frontend has been initialized (you may need to use a goroutine)")
		}
	}
	<-s.ready
}

func (s *Store) initialize() {
	s.once.Do(func() {
		close(s.ready)
	})
}
