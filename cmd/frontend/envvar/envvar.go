package envvar

import (
	"github.com/prince1809/sourcegraph/pkg/env"
	"strconv"
)

var sourcegraphDotComMode, _ = strconv.ParseBool(env.Get("SOURCEGRAPHDOTCOM_MODE", "false", "run as Sourcegraph.com with additional marketing and redirects"))

// SourcegraphDotComMode is true if this server is running Sourcegraph.com (solely by checking the
// SOURCEGRAPHDOTCOM_MODE env var). Sourcegraph.com shows add'l marketing and set up some add'l
// redirects.
func SourcegraphDotComMode() bool {
	return sourcegraphDotComMode
}

// MockSourcegraphDotComMode is used by tests to mock the results of SourcegraphDotComMode.
func MockSourcegraphDotComMode(value bool) {
	sourcegraphDotComMode = value

}
