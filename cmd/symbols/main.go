package main

import (
	"context"
	"fmt"
	"github.com/prince1809/sourcegraph/cmd/symbols/internal/symbols"
	"github.com/prince1809/sourcegraph/pkg/api"
	"github.com/prince1809/sourcegraph/pkg/debugserver"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/gitserver"
	"github.com/prince1809/sourcegraph/pkg/tracer"
	"github.com/prince1809/sourcegraph/pkg/vcs/git"
	"io"
	"log"
	"runtime"
	"strconv"
)

var (
	cacheDir       = env.Get("CACHE_DIR", "/tmp/symbols-cache", "directory to store cached symbols")
	cacheSizeMB    = env.Get("SYMBOLS_CACHE_SIZE_MB", "100000", "maximum size of the disk in megabytes ")
	ctagsProcesses = env.Get("CTAGS_PROCESSES", strconv.Itoa(runtime.NumCPU()), "number of ctags child processes to run")
)

const port = "3184"

func main() {
	env.Lock()
	env.HandleHelpFlag()
	log.SetFlags(0)
	tracer.Init()

	symbols.MustRegisterSqlite3WithPcre()

	go debugserver.Start()

	service := symbols.Service{
		FetchTar: func(ctx context.Context, repo gitserver.Repo, commit api.CommitID) (io.ReadCloser, error) {
			return git.Archive(ctx, repo, git.ArchiveOptions{Treeish: string(commit), Format: "tar"})
		},
	}

	if mb, err := strconv.ParseInt(cacheSizeMB, 10, 64); err != nil {
		log.Fatalf("Invalid SYMBOLS_CACHE_SIZE_MB: %s", err)
	} else {
		service.MaxCacheSizeBytes = mb * 1000 * 1000
	}

	var err error
	service.NumParserProcesses, err = strconv.Atoi(ctagsProcesses)
	if err != nil {
		log.Fatalf("Invalid CTAGS_PROCESSES: %s", err)
	}
	if err := service.Start(); err != nil {
		log.Fatalln("Start:", err)
	}
	fmt.Println("symbols")
}
