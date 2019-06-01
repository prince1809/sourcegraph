package confdb

import (
	"github.com/prince1809/sourcegraph/pkg/db/dbtesting"
	"testing"
)

func TestCriticalGetLatest(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := dbtesting.TestContext(t)

	latest, err := CriticalGetLatest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if latest == nil {
		t.Errorf("expected non-nil latest config since default config should be created, got: %+v", latest)
	}
}
