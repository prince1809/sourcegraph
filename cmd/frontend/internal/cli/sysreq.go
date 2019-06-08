package cli

import (
	"context"
	"fmt"
	"github.com/kr/text"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prince1809/sourcegraph/pkg/sysreq"
	"gopkg.in/inconshreveable/log15.v2"
	"io"
	"strings"
)

const skipSysReqsEnvVar = "SRC_SKIP_REQS"

var skipSysReqEnv = env.Get(skipSysReqsEnvVar, "false", "skip system requirement checks")

// skippedSysReqs returns a list of sysreq names to skip (e.g.,
// "Docker").
func skippedSysReqs() []string {
	return strings.Fields(skipSysReqEnv)
}

// checkSysReqs uses package sysreq to check for the presence of
// system requirements. If any are missing, it prints a message to
// w and returns a non-nil error.
func checkSysReqs(ctx context.Context, w io.Writer) error {
	wrap := func(s string) string {
		const indent = "\t\t"
		return strings.TrimPrefix(text.Indent(text.Wrap(s, 72), "\t\t"), indent)
	}

	var failed []string
	for _, st := range sysreq.Check(ctx, skippedSysReqs()) {
		if st.Failed() {
			failed = append(failed, st.Name)

			fmt.Fprint(w, " !!!!! ")
			fmt.Fprintf(w, " %s is required\n", st.Name)
			if st.Problem != "" {
				fmt.Fprint(w, "\tProblem: ")
				fmt.Fprintln(w, wrap(st.Problem))
			}
			if st.Err != nil {
				fmt.Fprint(w, "\tError: ")
				fmt.Fprintln(w, wrap(st.Err.Error()))
			}
			if st.Fix != "" {
				fmt.Fprint(w, "\tPossible fix: ")
				fmt.Fprintln(w, wrap(st.Fix))
			}
			fmt.Fprintln(w, "\t"+wrap(fmt.Sprintf("Skip this check by setting the env var %s=%q (separate multiple entries with spaces). Note: Sourcegraph may not function properly without %s.", skipSysReqsEnvVar, st.Name, st.Name)))
		}
	}

	if failed != nil {
		log15.Error("System requirement checks failed", "failed", failed)
	}
	return nil
}
