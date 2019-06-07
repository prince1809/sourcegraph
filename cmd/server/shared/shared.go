// Package shared provides the entrypoint to Sourcegraph's single docker
// image. It has functionality to setup the shared environment variables, as
// well as create the Procfile for goreman to run.
package shared

import "log"

// FrontendInternalHost is the value of SRC_FRONTEND_INTERNAL
const FrontendInternalHost = "127.0.0.1:3090"

// defaultEnv is environment variable that will be set if not already set.
var defaultEnv = map[string]string{
	// Sourcegraph services running in this container
	"SRC_GIT_SERVERS": "127.0.0.1:3178",
}

func Main() {
	log.SetFlags(0)
}
