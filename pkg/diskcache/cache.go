package diskcache

import (
	"context"
	"io"
	"os"
	"time"
)

type Store struct {
	Dir string

	Component string

	BackgroundTimeout time.Duration

	BeforeEvict func(string)
}

type File struct {
	*os.File

	Path string
}

type Fetcher func(context.Context) (io.ReadCloser, error)
