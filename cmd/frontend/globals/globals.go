package globals

import (
	"github.com/prince1809/sourcegraph/pkg/conf"
	"net/url"
)

// ExternalURL is the fully-resolved, externally accessible frontend URL.
var ExternalURL = &url.URL{Scheme: "http", Host: "example.com"}

// ConfigurationServerFrontendOnly provides the contents of the site configuration
// to other services and manages modifications to it.
//
// Any another service that attempts to use this variable will panic.
var ConfigurationServerFrontendOnly *conf.Server
