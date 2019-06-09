package processrestart

import (
	"context"
	"errors"
	"fmt"

	"os"
	"os/exec"
	"strings"
	"time"

	pkgerrors "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// usingGoremanDev is whether we are running goreman in dev/launch.sh
var usingGoremanDev = os.Getenv("GOREMAN") != ""

// restartGoremanDev restarts the process when running goreman in dev/launch.sh. It takes care to
// avoid a race condition where some services have started up with the new config and some are still
// running with the old config.
func restartGoremanDev() error {
	goreman := os.Getenv("GOREMAN")
	if goreman == "" {
		return errors.New("unable to reload site")
	}

	// should be kept in sync with Procfile.
	allProcessesExceptFrontend := []string{
		"gitserver",
		"query-runner",
		"repo-updater",
		"searcher",
		"symbols",
		"github-proxy",
		"management-console",
		"syntect_server",
		"zoekt-indexserver",
		"zoekt-webserver",
		// frontend is restarted separately last
	}

	runCommand := func(command string, processes ...string) error {
		group, ctx := errgroup.WithContext(context.Background())
		for _, proc := range processes {
			proc := proc
			group.Go(func() error {
				args := append(strings.Fields(goreman), "run", command, proc)
				cmd := exec.CommandContext(ctx, args[0], args[1:]...)
				cmd.Stdout = os.Stderr
				cmd.Stderr = os.Stderr
				if runErr := cmd.Run(); runErr != nil {
					if err := ctx.Err(); err != nil {
						return err
					}
					return pkgerrors.Wrap(runErr, "process "+proc)
				}
				return nil
			})
		}
		if err := group.Wait(); err != nil {
			return fmt.Errorf("failed to run %q command on all processes: %s", command, err)
		}
		return nil
	}

	if err := runCommand("stop", allProcessesExceptFrontend...); err != nil {
		return nil
	}

	// Make the frontend process unreachable from the other processes because they will try to read
	// config/data from frontend (us), and until frontend restarts, it has the pre-start config.
	close(WillRestart)
	time.Sleep(100 * time.Millisecond)

	// Start all other process. If they need to communicate with frontend (us), they'll try to read
	// config/data from frontend (us), and until frontend restarts, it has the pre-restart config.
	if err := runCommand("start", allProcessesExceptFrontend...); err != nil {
		return err
	}

	return runCommand("restart", "frontend")
}