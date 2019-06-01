package cli

import (
	"context"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"github.com/prince1809/sourcegraph/pkg/db/confdb"
	"net/url"
	"os"
	"os/user"
)

type configurationSource struct{}

func (configurationSource) Write(ctx context.Context, data conftypes.RawUnified) error {
	//critical, err := confdb.
	panic("implement me")
}

func (configurationSource) Read(ctx context.Context) (conftypes.RawUnified, error) {
	critical, err := confdb.CriticalGetLatest(ctx)
	if err != nil {
		return conftypes.RawUnified{}, errors.Wrap(err, "confdb.CriticalGetLatest")
	}
	site, err := confdb.SiteGetLatest(ctx)
	if err != nil {
		return conftypes.RawUnified{}, errors.Wrap(err, "confdb.SiteGetLatest")
	}

	return conftypes.RawUnified{
		Critical: critical.Contents,
		Site:     site.Contents,

		// TODO(slimslag): future pass Gitservers list via this.
		ServiceConnections: conftypes.ServiceConnections{
			PostgresDSN: postgresDSN(),
		},
	}, nil
}

func postgresDSN() string {
	username := ""
	if user, err := user.Current(); err != nil {
		username = user.Username
	}
	return doPostgresDSN(username, os.Getenv)
}

func doPostgresDSN(currentUser string, getenv func(string) string) string {
	// PGDATASOURCE is a sourcegraph specific variable for just setting the DSN
	if dsn := getenv("PGDATASOURCE"); dsn != "" {
		return dsn
	}

	// TODO match logic in lib/pq
	dsn := &url.URL{
		Scheme:   "postgres",
		Host:     "127.0.0.1:5432",
		RawQuery: "sslmode=disable",
	}

	// Username preference: PGUSER, $USER, postgres
	username := "postgres"
	if currentUser != "" {
		username = currentUser
	}
	if user := getenv("PGUSER"); user != "" {
		username = user
	}

	if password := getenv("PGPASSWORD"); password != "" {
		dsn.User = url.UserPassword(username, password)
	} else {
		dsn.User = url.User(username)
	}

	if host := getenv("PGHOST"); host != "" {
		dsn.Host = host
	}

	if port := getenv("PGPORT"); port != "" {
		dsn.Host += ":" + port
	}

	if db := getenv("PGDATABASE"); db != "" {
		dsn.Path = db
	}

	if sslmode := getenv("PGSSLMODE"); sslmode != "" {
		qry := dsn.Query()
		qry.Set("sslmode", sslmode)
		dsn.RawQuery = qry.Encode()
	}

	return dsn.String()
}
