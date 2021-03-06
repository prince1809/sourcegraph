package conf

import (
	"github.com/prince1809/sourcegraph/pkg/conf/conftypes"
	"github.com/prince1809/sourcegraph/schema"
	"log"
	"os"
	"path/filepath"
)

// Unified represents the overall global Sourcegraph configuration from various
// sources:
//
// - The critical configuration, from the database (from the management console).
// - The site configuration, from the database (from the site-admin panel).
// - Service connections, from the frontend (e.g. which gitservers to talk to).
type Unified struct {
	schema.SiteConfiguration
	Critical           schema.CriticalConfiguration
	ServiceConnections conftypes.ServiceConnections
}

type configurationMode int

const (
	// The user of pkg/conf reads and writes to the configuration file.
	// This should only ever be used by frontend.
	modeServer configurationMode = iota

	// The user of pkg/conf only reads the configuration file.
	modeClient

	// The user pkg/conf is a test case.
	modeTest
)

func getMode() configurationMode {
	mode := os.Getenv("CONFIGURATION_MODE")

	switch mode {
	case "server":
		return modeServer
	case "client":
		return modeClient
	default:
		// Detect 'go test' and default to test mode in that case.
		p, err := os.Executable()
		if err == nil && filepath.Ext(p) == ".test" {
			return modeTest
		}

		// Otherwise we default to client mode, so that most services need not
		// specify CONFIGURATION_MODE=client explicitly
		return modeClient

	}
}

var (
	configurationServerFrontendOnly            *Server
	configurationServerFrontendOnlyInitialized = make(chan struct{})
)

func init() {
	clientStore := NewStore()
	defaultClient = &client{store: clientStore}

	mode := getMode()

	// Don't kickoff the background updaters for the client/server
	// when running test cases.
	if mode == modeTest {
		close(configurationServerFrontendOnlyInitialized)

		// Seed the client store with a dummy configuration for test cases.
		_, err := clientStore.MaybeUpdate(conftypes.RawUnified{
			Critical:           "{}",
			Site:               "{}",
			ServiceConnections: conftypes.ServiceConnections{},
		})
		if err != nil {
			log.Fatalf("received error when setting up the store for the default client during test, err :%s", err)
		}
		return
	}

	// The default client is started in InitConfigurationServerFrontendOnly in
	// the case of server mode.
	if mode == modeClient {
		go defaultClient.continuouslyUpdate(nil)
		close(configurationServerFrontendOnlyInitialized)
	}
}

// InitConfigurationServerFrontendOnly creates and returns a configuration
// server. This should only be invoked by frontend, or else a panic will
// occur. This function should only ever be called once.
func InitConfigurationServerFrontendOnly(source ConfigurationSource) *Server {
	mode := getMode()

	if mode == modeTest {
		return nil
	}

	if mode == modeClient {
		panic("cannot call this function while in client mode")
	}

	server := NewServer(source)
	server.Start()

	// Install the passthrough configuration source for defaultClient. This is
	// so that the frontend does not request configuration from itself via HTTP
	// and instead only relies on the DB.
	defaultClient.passthrough = source

	go defaultClient.continuouslyUpdate(nil)
	configurationServerFrontendOnly = server
	close(configurationServerFrontendOnlyInitialized)

	return server
}
