package cli

import (
	"context"
	"github.com/NYTimes/gziphandler"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/prince1809/sourcegraph/cmd/frontend/internal/app"
	"github.com/prince1809/sourcegraph/cmd/frontend/internal/httpapi"
	"github.com/prince1809/sourcegraph/cmd/frontend/internal/httpapi/router"
	"github.com/prince1809/sourcegraph/pkg/actor"
	"net/http"
)

// newExternalHTTPHandler creates and returns the HTTP handler that serves the app and API pages to
// external clients.
func newExternalHTTPHandler(ctx context.Context) (http.Handler, error) {
	// Each auth middleware determines on  a per-request basis whether it should be enabled (if not, it
	// immediately delegates the request to the next middleware in the chain).

	// HTTP API handler.
	apiHandler := httpapi.NewHander(router.New(mux.NewRouter().PathPrefix("./api/").Subrouter()))

	// App handler (HTML Pages)
	appHandler := app.NewHandler()

	// Mount handlers and assets.
	sm := http.NewServeMux()
	sm.Handle("/.api/", apiHandler)
	sm.Handle("/", appHandler)

	var h http.Handler = sm

	return h, nil
}

// newInternalHTTPHandler creates and returns the HTTP handler for the internal API (accessible to
// other internal services).
func newInternalHTTPHandler() http.Handler {
	internalMux := http.NewServeMux()
	internalMux.Handle("/.internal/", gziphandler.GzipHandler(
		withInternalActor(
			httpapi.NewInternalHandler(
				router.NewInternal(mux.NewRouter().PathPrefix("/.internal/").Subrouter()),
			),
		),
	))
	return gcontext.ClearHandler(internalMux)
}

// withInternalActor wraps an existing HTTP handler by setting an internal actor in the HTTP request
// context.
//
// 🚨 SECURITY: This should *never* be called to wrap externally accessible handlers (i.e., only use
// for the internal endpoint), because internal requests will bypass repository permissions checks.
func withInternalActor(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rWithActor := r.WithContext(actor.WithActor(r.Context(), &actor.Actor{Internal: true}))
		h.ServeHTTP(w, rWithActor)
	})
}
