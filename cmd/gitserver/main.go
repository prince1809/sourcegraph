// gitserver is the gitserver server.
package main

import (
	"fmt"
	"github.com/prince1809/sourcegraph/pkg/env"
	"log"
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

	if reposDir == "" {
		log.Fatal("git-server: SRC_REPOS_DIR is required")
	}
	if err := os.MkdirAll(reposDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create SRC_REPOS_DIR: %s", err)
	}

	wantFreeG, err := strconv.Atoi(wantFreeG)
	if err != nil {
		log.Fatalf("parsing $SRC_REPOS_DESIRED_FREE_GB: %v", err)
	}

	fmt.Println("gitserver started")
}
