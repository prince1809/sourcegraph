package httpapi

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/cmd/frontend/globals"
	"net/http"
)

func serveConfiguration(w http.ResponseWriter, r *http.Request) error {
	raw, err := globals.ConfigurationServerFrontendOnly.Source.Read(r.Context())
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(raw)
	if err != nil {
		return errors.Wrap(err, "Encode")
	}
	return nil
}
func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
