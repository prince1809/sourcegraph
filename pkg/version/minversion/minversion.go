package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	version "github.com/mcuadros/go-version"
)

func main() {
	minimumVersion := "1.12"
	rawVersion := runtime.Version()
	versionNumber := strings.TrimPrefix(rawVersion, "go")
	minimumVersionMet := version.Compare(minimumVersion, versionNumber, "<=")
	if !minimumVersionMet {
		fmt.Printf("Go version %s or newer must be used; found: %s\n", minimumVersion, versionNumber)
		os.Exit(1) // minimum version not met means non-zero exit code
	}
	fmt.Printf("Go Version %s\n", versionNumber)
}
