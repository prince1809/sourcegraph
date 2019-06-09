package ui

import (
	"github.com/gorilla/mux"
	uirouter "github.com/prince1809/sourcegraph/cmd/frontend/internal/app/ui/router"
)

const (
	routeHome  = "home"
	souteStart = "start"
)

// Router returns the router that serves pages for our web app.
func Router() *mux.Router {
	return uirouter.Router
}
