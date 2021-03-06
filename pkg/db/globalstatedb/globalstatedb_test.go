package globalstatedb

import (
	"github.com/prince1809/sourcegraph/pkg/db/dbtesting"
	"testing"
)

func TestGet(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := dbtesting.TestContext(t)
	config, err := Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if config.SiteID == "" {
		t.Fatal("expected site_id to be set")
	}
}
