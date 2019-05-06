package cli

import (
	"github.com/prince1809/sourcegraph/cmd/frontend/globals"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"github.com/prince1809/sourcegraph/pkg/db/dbconn"
	"log"
)

// Main is the main entrypoint for the frontend server progral¥m.
func Main() error {
	log.SetFlags(0)
	log.SetPrefix("")

	// Connect to the database for the frontend server program.
	if err := dbconn.ConnectToDB(""); err != nil {
		log.Fatal(err)
	}
	globals.ConfigurationServerFrontendOnly = conf.InitConfigurationServerFrontendOnly(&configurationSource{})



	println("frontend server started")
	return nil
}
