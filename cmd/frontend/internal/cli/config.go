package cli

import (
	"context"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
)

type configurationSource struct{}

func (configurationSource) Write(ctx context.Context, data conftypes.RawUnified) error {
	critical, err := confdb.
	panic("implement me")
}

func (configurationSource) Read(ctx context.Context) (conftypes.RawUnified, error) {
	panic("implement me")
}
