package handlerutil

import (
	"github.com/gorilla/mux"
	"github.com/prince1809/sourcegraph/pkg/api"
	"net/http"
)

// RedirectToNewRepoName writes an HTTP redirect response with a
// Location that matches the request's location except with the
// Repo route var updated to refer to newRepoName (instead of the
// originally requested repo name).
func RedirectToNewRepoName(w http.ResponseWriter, r *http.Request, newRepoName api.RepoName) error {
	origVars := mux.Vars(r)
	origVars["Repo"] = string(newRepoName)

	var pairs []string
	for k, v := range origVars {
		pairs = append(pairs, k, v)
	}
	destURL, err := mux.CurrentRoute(r).URLPath(pairs...)
	if err != nil {
		return err
	}

	http.Redirect(w, r, destURL.String(), http.StatusMovedPermanently)
	return nil
}
