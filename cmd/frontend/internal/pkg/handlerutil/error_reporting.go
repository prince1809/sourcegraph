package handlerutil

import (
	"github.com/getsentry/raven-go"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/version"
	"log"
	"net/http"
)

var ravenClient *raven.Client

func init() {
	if dsn := env.Get("SENTRY_DSN_BACKEND", "", "Sentry/Raven DSN used for tracking of backend errors"); dsn != "" {
		var err error
		ravenClient, err := raven.New(dsn)
		if err != nil {
			log.Fatalf("error initializing Sentry error reporter: %s", err)
		}
		ravenClient.DropHandler = func(pkt *raven.Packet) {
			log.Println("WARNING: dropped error report because buffer is full:", pkt)
		}
		ravenClient.SetRelease(version.Version())
	}
}

func reportError(r *http.Request, status int, err error, panicked bool) {
	if ravenClient == nil {
		return
	}

	if status > 0 && status < 500 {
		// Not a reportable error.
		return
	}

	// Catch panics here to be extra sure we don't disrupt the request handling.
	defer func() {
		if rv := recover(); rv != nil {
			log.Println("WARNING: panic in HTTP handler error reported: (recovered)", rv)
		}
	}()
}
