package router

import "github.com/gorilla/mux"

const (
	GraphQL = "graphql"
)

// New creates a new API router with router URL pattern definitions but
// no handlers attached to the routes.
func New(base *mux.Router) *mux.Router {

	return base
}
