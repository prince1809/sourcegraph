package git

import (
	"context"
	"github.com/prince1809/sourcegraph/pkg/gitserver"
	"io"
)

type ArchiveOptions struct {
	Treeish string
	Format  string
	Paths   []string
}

func Archive(ctx context.Context, repo gitserver.Repo, opt ArchiveOptions) (_ io.ReadCloser, err error) {
	panic("implement me")
}
