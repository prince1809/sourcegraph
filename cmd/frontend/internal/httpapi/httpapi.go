package httpapi

import (
	"github.com/gorilla/mux"
	"net/http"
)

// NewInternalHTTPHandler returns a new API handler for internal endpoints that uses
// the provided API router, which must have been created by httpapi/router.NewInternal.
//
// 🚨 SECURITY: This handler should not be served on a publicly exposed port. 🚨
// This handler is not guaranteed to provide the same authorization checks as
// public API handlers.
func NewInternalHandler(m *mux.Router) http.Handler {
	if m == nil {

	}

	return m
}
