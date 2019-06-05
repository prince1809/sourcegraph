package db

import (
	"context"
	"database/sql"
)

// queryable allows us to reuse the same logic for certain operations both
// inside and outside a explicit transactions.
type queryable interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}
