package httpapi

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/prince1809/sourcegraph/cmd/frontend/db"
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


func serveReposList(w http.ResponseWriter, r *http.Request) error{
	var opt db.ReposListOptions
	err := json.NewDecoder(r.Body).Decode(&opt)
	if err != nil {
		return err
	}
	//res, err := backen
	w.WriteHeader(http.StatusOK)
	//w.Write()
	return nil
}
func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
