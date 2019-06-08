// +build linux darwin

package sysreq

import (
	"context"
	"fmt"
	"github.com/prince1809/sourcegraph/pkg/conf"
	"golang.org/x/sys/unix"
)

func rlimitCheck(ctx context.Context) (problem, fix string, err error) {
	const minLimit = 10000

	var limit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &limit); err != nil {
		return "", "", err
	}

	if limit.Cur < minLimit {
		fix := fmt.Sprintf(`Please increase the open file limit by running "ulimit -n %d."`, minLimit)

		if conf.IsDeployTypeDockerContainer(conf.DeployType()) {
			fix = fmt.Sprintf("Add --ulimit nofile=%d:%d to the docker run command", minLimit, minLimit)
		}

		return "Insufficient file descriptor limit", fix, nil
	}
	return
}
