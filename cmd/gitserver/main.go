// gitserver is the gitserver server.
package main

import (
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/prince1809/sourcegraph/cmd/gitserver/server"
	"github.com/prince1809/sourcegraph/pkg/debugserver"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/tracer"
	"gopkg.in/inconshreveable/log15.v2"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

var (
	reposDir          = env.Get("SRC_REPOS_DIR", "/data/repos", "Root dir containing repos.")
	runRepoCleanup, _ = strconv.ParseBool(env.Get("SRC_RUN_REPO_CLEANUP", "", "Periodically remove inactive repositories."))
	wantFreeG         = env.Get("SRC_REPOS_DESIRED_FREE_GB", "10", "How many gigabytes of space to keep free on the disk with the repos")
	janitorInterval   = env.Get("SRC_REPOS_JANITOR_INTERVAL", "1m", "Interval between cleanup runs")
)

func main() {
	env.Lock()
	env.HandleHelpFlag()
	tracer.Init()

	log.Println("server starting")

	if reposDir == "" {
		log.Fatal("git-server: SRC_REPOS_DIR is required")
	}
	if err := os.MkdirAll(reposDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create SRC_REPOS_DIR: %s", err)
	}

	wantFreeG2, err := strconv.Atoi(wantFreeG)
	if err != nil {
		log.Fatalf("parsing $SRC_REPOS_DESIRED_FREE_GB: %v", err)
	}
	gitserver := server.Server{
		ReposDir:                reposDir,
		DeleteStaleRepositories: runRepoCleanup,
		DesiredFreeDiskSpace:    uint64(wantFreeG2 * 1024 * 1024 * 1024),
	}
	gitserver.RegisterMetrics()

	if tmpDir, err := gitserver.SetupAndClearTmp(); err != nil {
		log.Fatalf("failed to setup temporary directory: %s", err)
	} else {
		os.Setenv("TMP_DIR", tmpDir)
	}

	// Create Handler now since it also initializes state
	handler := nethttp.Middleware(opentracing.GlobalTracer(), gitserver.Handler())

	go debugserver.Start()

	port := "3178"
	host := ""
	if env.InsecureDev {
		host = "127.0.0.1"
	}
	addr := net.JoinHostPort(host, port)
	srv := &http.Server{Addr: addr, Handler: handler}
	log15.Info("git-server: listening", "addr", srv.Addr)

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Stop accepting requests. In the future we should use graceful shutdown.
	srv.Close()

	// The most important thins this does is kill all our clones. if we just
	// shutdown they will be orphaned and continue running.
	gitserver.Stop()
}
