package cli

import (
	"context"
	"net/http"
)

// newExternalHTTPHandler creates and returns the HTTP handler that serves the app and API pages to
// external clients.
func newExternalHTTPHandler(ctx context.Context) (http.Handler, error) {
	// Each auth middleware determines on  a per-request basis whether it should be enabled (if not, it
	// immediately delegates the request to the next middleware in the chain).
	//authMiddlewares := auth.AuthMiddleware()
	return nil, nil
}
