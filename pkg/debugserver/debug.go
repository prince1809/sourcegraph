package debugserver

import (
	"encoding/json"
	"github.com/prince1809/sourcegraph/pkg/env"
	"net/http"
	"os"
	"path/filepath"
)

var (
	addr = env.Get("SRC_PROF_HTTP", ":6060", "net/http/pprof http bind address.")
)

func init() {
	err := json.Unmarshal([]byte(env.Get("SRC_PROF_SERVICES", "[]", "List of net/http/pprof  http bind address.")), &Services)
	if err != nil {
		panic("failed to JSON unmarshal SRC_PROF_SERVICES: " + err.Error())
	}

	if addr == "" {
		// Look for our binname in the services list
		name := filepath.Base(os.Args[0])
		for _, svc := range Services {
			if svc.Name == name {
				addr = svc.Host
				break
			}
		}
	}
}

// Endpoint is a handler for the debug server. It will be displayed on the
// debug index page.
type Endpoint struct {
	// Name is the name shown on the index page for the endpoint.
	Name string
	// Path is passed to http.Mux.Handle as the pattern.
	Path string
	// Handler is the debug handler
	Handler http.Handler
}

// Services is the list of registered services debug addresses. Populated
// from SRC_PROF_MAP
var Services []Service

// Services is a service's debug addr (host:port).
type Service struct {
	// Name of the service. Always the binary name. example: "gitserver"
	Name string

	// Host is the host:port for the services SRC_PROF_HTTP. example:
	// "127.0.0.1:6060"
	Host string

	// DefaultPath is the path to the service we should link to.
	DefaultPath string
}

// Start runs a debug server (pprof, prometheus, etc) if it is configured ( via
// SRC_PROF_HTTP environment variable). It is blocking.
func Start(extra ...Endpoint) {
	if addr == "" {
		return
	}
}
