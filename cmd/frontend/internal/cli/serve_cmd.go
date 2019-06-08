package cli

import (
	"context"
	"fmt"
	"github.com/keegancsmith/tmpfriend"
	"github.com/prince1809/sourcegraph/cmd/frontend/globals"
	"github.com/prince1809/sourcegraph/cmd/frontend/internal/cli/loghandlers"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"github.com/prince1809/sourcegraph/pkg/db/dbconn"
	"github.com/prince1809/sourcegraph/pkg/debugserver"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/sysreq"
	"github.com/prince1809/sourcegraph/pkg/tracer"
	"github.com/prince1809/sourcegraph/pkg/version"
	"github.com/prince1809/sourcegraph/pkg/vfsutil"
	"gopkg.in/inconshreveable/log15.v2"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
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

// Main is the main entrypoint for the frontend server program.
func Main() error {
	log.SetFlags(0)
	log.SetPrefix("")

	// Connect to the database for the frontend server program.
	if err := dbconn.ConnectToDB(""); err != nil {
		log.Fatal(err)
	}
	globals.ConfigurationServerFrontendOnly = conf.InitConfigurationServerFrontendOnly(&configurationSource{})
	fmt.Println("validate config")
	conf.MustValidateDefaults()
	handleConfigOverrides()

	// Filter trace logs
	d, _ := time.ParseDuration(traceThreshold)
	tracer.Init(tracer.Filter(loghandlers.Trace(strings.Fields(trace), d)))

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "help", "-h", "--help":
			log.Printf("Version: %s", version.Version())
			log.Print()

			env.PrintHelp()

			log.Print()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			for _, st := range sysreq.Check(ctx, skippedSysReqs()) {
				log.Printf("%s:", st.Name)
				if st.OK() {
					log.Print("\tOK")
					continue
				}
				if st.Skipped {
					log.Print("\tSkipped")
					continue
				}
				if st.Problem != "" {
					log.Print("\t" + st.Problem)
				}
				if st.Err != nil {
					log.Printf("\tError: %s", st.Err)
				}
				if st.Fix != "" {
					log.Printf("\tPossible fix: %s", st.Fix)
				}
			}
			return nil
		}
	}

	fmt.Println("###### Validate config ######")
	printConfigValidation()

	cleanup := tmpfriend.SetupOrNOOP()
	defer cleanup()

	// Don't proceed if system requirements are missing, to avoid
	// presenting uses with a half-working experience.
	if err := checkSysReqs(context.Background(), os.Stderr); err != nil {
		return err
	}

	go debugserver.Start()

	//siteid.init()


	if printLogo {
		fmt.Println(" ")
		fmt.Println(logoColor)
		fmt.Println(" ")
	}

	// Create the external HTTP handler
	externalHandler, err := newExternalHTTPHandler(context.Background())
	if err != nil {
		return err
	}

	// serve will serve externalHanlder on l, It additionally handles graceful restarts.
	srv := &httpServers{}

	// Start the http server.
	l, err := net.Listen("tcp", httpAddr)
	if err != nil {
		return err
	}
	log15.Debug("HTTP running", "on", httpAddr)
	srv.GoServe(l, &http.Server{
		Handler:      externalHandler,
		ReadTimeout:  75 * time.Second,
		WriteTimeout: 60 * time.Second,
	})

	fmt.Printf("* Sourcegraph is ready at: %s \n", globals.ExternalURL)

	srv.Wait()

	return nil
}

type httpServers struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	servers []*http.Server
	wrapper func(http.Handler) http.Handler
}

// GoServe serves srv in a new goroutine. If serve returns an error other than
// http.ErrServerClosed it will fatal.
func (s *httpServers) GoServe(l net.Listener, srv *http.Server) {
	s.addServer(srv)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := srv.Serve(l); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func (s *httpServers) addServer(srv *http.Server) *http.Server {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.wrapper != nil {
		srv.Handler = s.wrapper(srv.Handler)
	}
	s.servers = append(s.servers, srv)
	return srv
}

// Wait waits until all servers are closed.
func (s *httpServers) Wait() {
	s.wg.Wait()
}

func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	return false
}
