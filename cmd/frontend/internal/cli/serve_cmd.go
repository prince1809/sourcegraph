package cli

import (
	"fmt"
	"github.com/prince1809/sourcegraph/cmd/frontend/globals"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"github.com/prince1809/sourcegraph/pkg/db/dbconn"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/vfsutil"
	"log"
	"path/filepath"
	"strconv"
)

var (
	trace          = env.Get("SRC_LOG_TRACE", "HTTP", "space separated list of trace logs to show. Options: all,")
	traceThreshold = env.Get("SRC_LOG_TRACE_THRESHOLD", "", "show traces that take longer than this")

	printLogo, _ = strconv.ParseBool(env.Get("LOGO", "false", "print Sourcegraph logo upon startup"))

	httpAddr         = env.Get("SRC_HTTP_ADDR", ":3080", "HTTP listen address for app and HTTP API")
	httpAddrInternal = env.Get("SRC_HTTP_ADDR_INTERNAL", ":3090", "HTTP listen address for internal HTTP API. This should never be exposed externally, as it lacks certain authz checks.")

	nginxAddr = env.Get("SRC_NGINX_HTTP_ADDR", "", "HTTP listems address for nginx reverse proxy to SRC_HTTP_ADDR. Has preference over SRC_HTTP_ADDR for ExternalURL.")

	// dev browser browser extension ID. You can find this by going to chrome://extensions
	devExtension = "chrome-extension://bmfbcejdknlknpncfpeloejonjoledha"
	// production browser extension ID. This is found by viewing our extension in the chrome store.
	prodExtension = "chrome-extension://dgjhfomjieaadpoljlnidmbgkdffpack"
)

func init() {
	// If CACHE_DIR is specified, use that
	cacheDir := env.Get("CACHE_DIR", "/tmp", "directory to store cached archives.")
	vfsutil.ArchiveCacheDir = filepath.Join(cacheDir, "frontend-archive-cache")

}

// Main is the main entrypoint for the frontend server progral¥m.
func Main() error {
	log.SetFlags(0)
	log.SetPrefix("")

	// Connect to the database for the frontend server program.
	if err := dbconn.ConnectToDB(""); err != nil {
		log.Fatal(err)
	}
	globals.ConfigurationServerFrontendOnly = conf.InitConfigurationServerFrontendOnly(&configurationSource{})
	if printLogo {
		fmt.Println(" ")
		fmt.Println(logoColor)
		fmt.Println(" ")
	}
	fmt.Printf("* Sourcegraph is ready at: %s \n", globals.ExternalURL)
	return nil
}
