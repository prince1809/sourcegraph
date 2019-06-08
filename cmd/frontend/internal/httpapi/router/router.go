package router

import "github.com/gorilla/mux"

const (
	GraphQL = "graphql"

	Registry    = "registry"
	RepoRefresh = "repo.refresh"
	Telemetry   = "telemetry"

	SavedQueriesListAll = "internal.saved-queries.list-all"
)

// New creates a new API router with router URL pattern definitions but
// no handlers attached to the routes.
func New(base *mux.Router) *mux.Router {
	if base == nil {
		panic("base == nil")
	}

	base.StrictSlash(true)

	addRegistryRoute(base)
	addGraphQLRoute(base)
	addTelemetryRoute(base)
	
	return base
}

// NewInternal creates a new API router for internal endpoints.
func NewInternal(base *mux.Router) *mux.Router {
	if base == nil {
		panic("base == nil")
	}

	base.StrictSlash(true)
	// Internal API endpoints should only be served on the internal Handler
	base.Path("/saved-queries/list-all").Methods("POST").Name(SavedQueriesListAll)

	addRegistryRoute(base)
	addGraphQLRoute(base)
	addTelemetryRoute(base)

	return base
}

func addRegistryRoute(m *mux.Router) {
	m.PathPrefix("/registry").Methods("GET").Name(Registry)
}

func addTelemetryRoute(m *mux.Router) {
	m.Path("/telemetry/{TelemetryPath:*}").Methods("POST").Name(Telemetry)
}

func addGraphQLRoute(m *mux.Router) {
	m.Path("/graphql").Methods("POST").Name(GraphQL)
}
