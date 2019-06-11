package symbols

import (
	"github.com/google/zoekt/ctags"
	"github.com/pkg/errors"
	"runtime"
)

func (s *Service) startParsers() error {
	n := s.NumParserProcesses
	if n == 0 {
		n = runtime.GOMAXPROCS(0)
	}

	s.parsers = make(chan ctags.Parser, n)
	for i := 0; i < n; i++ {
		parser, err := s.NewParser()
		if err != nil {
			return errors.Wrap(err, "NewParser")
		}
		s.parsers <- parser
	}
	return nil
}
