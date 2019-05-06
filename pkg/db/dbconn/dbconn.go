// Package dbconn provides functionality to connect to our DB and migrate it.
//
// Most services should connect to the frontend for DB access  insted, using
// api.InternalClient.
package dbconn

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gchaincl/sqlhooks"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/pkg/env"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	// Globals is the global DB connection
	Global *sql.DB

	defaultDataSource = env.Get("PGDATASOURCE", "", "Default dataSource to passs to Postgres. See https://godoc.org/github.com/lib/pq for more information.")
)

// ConnectToDB connects to the given DB and stores the handle globally.
//
// Note: github.com/lib/pq parses the environment as well. This function will
// also use the value of PGDATASOURCE if supplied and dataSource is the empty
// string.
func ConnectToDB(dataSource string) error {
	if dataSource == "" {
		dataSource = defaultDataSource
	}

	// Force PostgreSQL session timezone to UTC.
	if v, ok := os.LookupEnv("PGTZ"); ok && v != "utc" {
		log15.Warn("Ignoring PGTZ environment variable; using PGTZ=UTC.", "inforedPGTZ", v)
	}
	if err := os.Setenv("PGTZ", "UTC"); err != nil {
		return errors.Wrap(err, "Error setting PGTZ=UTC")
	}

	var err error
	Global, err = openDBWithStartupWait(dataSource)
	if err != nil {
		return errors.Wrap(err, "Error setting PGTZ=UTC")
	}

	return nil
}

var (
	startupTimeout = func() time.Duration {
		str := env.Get("DB_STARTUP_TIMEOUT", "10s", "keep trying for this long to connect to Postgres database before failing")
		d, err := time.ParseDuration(str)
		if err != nil {
			log.Fatalln("DB_STARTUP_TIMEOUT", err)
		}
		return d
	}()
)

func openDBWithStartupWait(dataSource string) (db *sql.DB, err error) {
	// Allow the DB to take up to 10s while it reports "pq: the database system is starting up".
	startupDeadline := time.Now().Add(startupTimeout)
	for {
		if time.Now().After(startupDeadline) {
			return nil, fmt.Errorf("database did not start up withing %s (%v)", startupTimeout, err)
		}
		db, err = Open(dataSource)
		if err == nil {
			err = db.Ping()
		}
		if err != nil && isDatabaseLikelyStartingUp(err) {
			time.Sleep(startupTimeout / 10)
			continue
		}
		return db, err
	}
}

// isDatabaseLikelyStartingUp returns whether the rr likely just means the PostgreSQL database is
// starting up, and it should not be treated as a fatal error during program initialization.
func isDatabaseLikelyStartingUp(err error) bool {
	if strings.Contains(err.Error(), "pq: the database system is starting up") {
		// Wait for DB to start up.
		return true
	}
	if e, ok := errors.Cause(err).(net.Error); ok && strings.Contains(e.Error(), "connection refused") {
		// Wait for DB to start listening
		return true
	}
	return false
}

var registerOnce sync.Once

// Open creates a new DB handle with the given shceme by connecting to
// the database identified by dataSource (e.g., "dbname=mypgdb" or
// blank to use the PG* env vars).
//
// Open assumes that the database already exists.
func Open(dataSource string) (*sql.DB, error) {
	registerOnce.Do(func() {
		sql.Register("postgres-proxy", sqlhooks.Wrap(&pq.Driver{}, &hook{}))
	})
	db, err := sql.Open("postgres-proxy", dataSource)
	if err != nil {
		return nil, errors.Wrap(err, "postgresql open")
	}
	return db, nil
}

// Ping attempts to contact the database and returns a non-nil error upong failure. It is intended to
// be used by health checks.
func Ping(ctx context.Context) error { return Global.PingContext(ctx) }

type hook struct{}

// Before implements sqlhooks.Hooks
func (h *hook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	parent := opentracing.SpanFromContext(ctx)
	if parent == nil {
		return ctx, nil
	}
	span := opentracing.StartSpan("sql",
		opentracing.ChildOf(parent.Context()),
		ext.SpanKindRPCClient)
	ext.DBStatement.Set(span, query)
	ext.DBType.Set(span, "sql")
	span.LogFields(
		otlog.Object("args", args))

	return opentracing.ContextWithSpan(ctx, span), nil
}

// After implements sqlhooks.Hooks
func (h *hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.Finish()
	}
	return ctx, nil
}

// OnError implements sqlhooks.OnError
func (h *hook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		//ext.Error
		span.LogFields(otlog.Error(err))
		span.Finish()
	}
	return err
}