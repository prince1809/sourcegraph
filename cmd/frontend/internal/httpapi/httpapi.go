package httpapi

import (
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/prince1809/sourcegraph/cmd/frontend/internal/pkg/handlerutil"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/trace"
	"gopkg.in/inconshreveable/log15.v2"
	"log"
	"net/http"

	apirouter "github.com/prince1809/sourcegraph/cmd/frontend/internal/httpapi/router"
)

// NewInternalHTTPHandler returns a new API handler for internal endpoints that uses
// the provided API router, which must have been created by httpapi/router.NewInternal.
//
// 🚨 SECURITY: This handler should not be served on a publicly exposed port. 🚨
// This handler is not guaranteed to provide the same authorization checks as
// public API handlers.
func NewInternalHandler(m *mux.Router) http.Handler {
	if m == nil {
		m = apirouter.New(nil)
	}

	m.StrictSlash(true)

	m.Get(apirouter.Configuration).Handler(trace.TraceRoute(handler(serveConfiguration)))
	m.Path("/ping").Methods("GET").Name("ping").HandlerFunc(handlePing)

	m.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API no route: %s %s from %s", r.Method, r.URL, r.Referer())
		http.Error(w, "no route", http.StatusNotFound)
	})

	return m
}

// handler is a wrapper func for API handlers.
func handler(h func(http.ResponseWriter, *http.Request) error) http.Handler {
	return handlerutil.HandlerWithErrorReturn{
		Handler: func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("Content-Type", "application/json")
			return h(w, r)
		},
		Error: handleError,
	}
}

func handleError(w http.ResponseWriter, r *http.Request, status int, err error) {
	// Handle custom errors
	if ee, ok := err.(*handlerutil.URLMovedError); ok {
		err := handlerutil.RedirectToNewRepoName(w, r, ee.NewRepo)
		if err != nil {
			log15.Error("error redirecting to new URI", "err", err, "new_url", ee.NewRepo)
		}
		return
	}

	// Never cache error responses.
	w.Header().Set("cache-control", "no-cache, max-age=0")

	errBody := err.Error()

	var displayErrBody string
	if env.InsecureDev {
		// Only display error message to admins when in debug mode, since it may
		// contain sensitive info (like API keys in net/http error messages).
		displayErrBody = string(errBody)
	}
	http.Error(w, displayErrBody, status)
	traceSpan := opentracing.SpanFromContext(r.Context())
	var spanURL string
	if traceSpan != nil {
		spanURL = trace.SpanURL(traceSpan)
	}
	if status < 200 || status >= 500 {
		log15.Error("API HTTP handler error response", "method", r.Method, "request_uri", r.URL.RequestURI(), "status_code", status, "error", err, "trace", spanURL)
	}
}
