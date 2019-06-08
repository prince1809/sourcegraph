package handlerutil

import "github.com/prince1809/sourcegraph/pkg/api"

type URLMovedError struct {
	NewRepo api.RepoName `json:"RedirectTo"`
}

func (e *URLMovedError) Error() string {
	return "URL moved to " + string(e.NewRepo)
}
