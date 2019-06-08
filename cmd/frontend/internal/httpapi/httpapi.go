package httpapi

import (
	"github.com/gorilla/mux"
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

	m.Path("/ping").Methods("GET").Name("ping").HandlerFunc(handlePing)

	m.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("API no route: %s %s from %s", r.Method, r.URL, r.Referer())
		http.Error(w, "no route", http.StatusNotFound)
	})

	return m
}
