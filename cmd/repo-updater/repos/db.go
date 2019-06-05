package repos

import (
	"context"
	"database/sql"
)

// A DB captures the essential methods of a sql.DB.
type DB interface {
	QueryContext(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error)
}
