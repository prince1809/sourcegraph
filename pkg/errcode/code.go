// Package errcode maps Go errors to HTTP status code as well as other useful
// functions for inspecting errors.
package errcode

import "net/http"

// HTTP returns the most appropriate HTTP status code that describes
// err. It contains a hard-coded list of error types and error values
// (such as mapping store.RepoNotFoundError to NotFound) and
// heuristic (such as mapping os.IsNotExist-satisfying errors to
// NotFound). All other errors are mapped to HTTP 500 Internal Server
// Error.
func HTTP(err error) int {

	return http.StatusInternalServerError
}

type HTTPErr struct {
	Status int   // HTTP status code
	Err    error // Optional reason for the HTTP error.
}
