// Package shared contains the frontend command implementation shared
package shared

import (
	"fmt"
	"os"

	"github.com/prince1809/sourcegraph/cmd/frontend/internal/cli"
	"github.com/prince1809/sourcegraph/pkg/env"
)

// Main is the main function that runs the frontend process.
//
// It is exposed as function in a package so that is can be called by other
// main package implementation such as Sourcegraph Enterprise, which import
// proprietary/private code.
func Main() {
	env.Lock()
	err := cli.Main()
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}
