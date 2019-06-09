package app

import (
	"github.com/prince1809/sourcegraph/cmd/frontend/internal/app/router"
	"github.com/prince1809/sourcegraph/cmd/frontend/internal/app/ui"
	"net/http"
)

// NewHandler returns a new app handler that uses the app router.
//
// 🚨 SECURITY: The caller MUST wrap the returned handler in middleware that checks authentication
// and sets the actor in the request context.
func NewHandler() http.Handler {
	r := router.Router()

	m := http.NewServeMux()

	m.Handle("/", r)

	r.Get(router.UI).Handler(ui.Router())

	return m
}
