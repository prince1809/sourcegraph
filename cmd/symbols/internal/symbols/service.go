package symbols

import (
	"context"
	"github.com/google/zoekt/ctags"
	"github.com/prince1809/sourcegraph/pkg/api"
	"github.com/prince1809/sourcegraph/pkg/diskcache"
	"github.com/prince1809/sourcegraph/pkg/gitserver"
	"io"
	"time"
)

// Service is the symbols service.
type Service struct {
	FetchTar func(context.Context, gitserver.Repo, api.CommitID) (io.ReadCloser, error)

	MaxConcurrentFetchTar int

	NewParser func() (ctags.Parser, error)

	NumParserProcesses int

	Path string

	MaxCacheSizeBytes int64

	cache *diskcache.Store

	fetchSem chan int

	parsers chan ctags.Parser
}

func (s *Service) Start() error {
	if err := s.startParsers(); err != nil {
		return nil
	}

	if s.MaxConcurrentFetchTar == 0 {
		s.MaxConcurrentFetchTar = 15
	}

	s.fetchSem = make(chan int, s.MaxConcurrentFetchTar)

	s.cache = &diskcache.Store{
		Dir:               s.Path,
		Component:         "symbols",
		BackgroundTimeout: 20 * time.Minute,
	}

	go s.watchAndEvict()
	return nil
}


func (s *Service) watchAndEvict() {
	if s.MaxCacheSizeBytes == 0  {
		return
	}

	for {
		time.Sleep(10 * time.Second)
		//stat, err := s.cache.Ev
	}
}
