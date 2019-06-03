package cli

import (
	"context"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"github.com/prince1809/sourcegraph/pkg/db/confdb"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/user"
)

// handleConfigOverrides handles allowing dev environments to forcibly override
// the configuration in the database upon startup. This is used to e.g. ensure
// dev environments have a consistent configuration and to load secrets from
// a separate private repository
func handleConfigOverrides() {
	if conf.IsDev(conf.DeployType()) {
		raw := conf.Raw()

		devOverrideCriticalConfig := os.Getenv("DEV_OVERRIDE_CRITICAL_CONFIG")
		if devOverrideCriticalConfig != "" {
			critical, err := ioutil.ReadFile(devOverrideCriticalConfig)
			if err != nil {
				log.Fatal(err)
			}
			raw.Critical = string(critical)
		}
	}
}

type configurationSource struct{}

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

func (configurationSource) Write(ctx context.Context, data conftypes.RawUnified) error {
	//critical, err := confdb.
	panic("implement me")
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
