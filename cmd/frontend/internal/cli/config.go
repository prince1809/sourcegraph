package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/cmd/frontend/db"
	"github.com/prince1809/sourcegraph/cmd/frontend/types"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"github.com/prince1809/sourcegraph/pkg/db/confdb"
	"github.com/prince1809/sourcegraph/pkg/jsonc"
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

		devOverrideSiteConfig := os.Getenv("DEV_OVERRIDE_SITE_CONFIG")
		if devOverrideSiteConfig != "" {
			site, err := ioutil.ReadFile(devOverrideSiteConfig)
			if err != nil {
				log.Fatal(err)
			}
			raw.Site = string(site)
		}

		if devOverrideCriticalConfig != "" || devOverrideSiteConfig != "" {
			err := (&configurationSource{}).Write(context.Background(), raw)
			if err != nil {
				log.Fatal(err)
			}
		}

		devOverrideExtSvcConfig := os.Getenv("DEV_OVERRIDE_EXTSVC_CONFIG")
		if devOverrideExtSvcConfig != "" {
			existing, err := db.ExternalServices.List(context.Background(), db.ExternalServicesListOptions{})
			if err != nil {
				log.Fatal(err)
			}
			if len(existing) > 0 {
				return
			}

			extsvc, err := ioutil.ReadFile(devOverrideExtSvcConfig)
			if err != nil {
				log.Fatal(err)
			}
			var configs map[string][]*json.RawMessage
			if err := jsonc.Unmarshal(string(extsvc), &configs); err != nil {
				log.Fatal(err)
			}
			for key, cfgs := range configs {
				for i, cfg := range cfgs {
					marshaledCfg, err := json.MarshalIndent(cfg, "", "  ")
					if err != nil {
						log.Fatal(err)
					}
					if err := db.ExternalServices.Create(context.Background(), &types.ExternalService{
						Kind:        key,
						DisplayName: fmt.Sprintf("Dev %s #%d", key, i+1),
						Config:      string(marshaledCfg),
					}); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}

}

type configurationSource struct{}

func (c configurationSource) Read(ctx context.Context) (conftypes.RawUnified, error) {
	critical, err := confdb.CriticalGetLatest(ctx)
	fmt.Println("CriticalGetLatest:","CriticalGetLatest:",  critical)
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

func (c configurationSource) Write(ctx context.Context, data conftypes.RawUnified) error {
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
