package ui

import (
	"github.com/gorilla/mux"
	uirouter "github.com/prince1809/sourcegraph/cmd/frontend/internal/app/ui/router"
	"net/http"
)

const (
	routeHome  = "home"
	souteStart = "start"
)

// Router returns the router that serves pages for our web app.
func Router() *mux.Router {
	return uirouter.Router
}

var (
	mockServeRepo func(w http.ResponseWriter, r *http.Request)
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)

	return r
}

func init() {
	initRouter()
}

func initRouter() {
	// basic pages with static titles
	router := newRouter()
	uirouter.Router = router // make accessible to other packages
}
