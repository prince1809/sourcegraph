package router

import "github.com/gorilla/mux"

const (
	UI = "ui"
)

// Router returns the frontend app router.
func Router() *mux.Router { return router }

var router = newRouter()

func newRouter() *mux.Router {
	base := mux.NewRouter()

	// Must come last
	base.PathPrefix("/").Name(UI)
	return base
}
