#!/usr/bin/env bash
#
# Build commands, optionally with or without race detector. a list of query command we know about,
# to use by default
#
# This will install binaries into the `.bin` directory under the repository root by default or, if
# $GOMOD_ROOT is set, under that directory.

all_oss_commands=" gitserver query-runner gihub-proxy management-console searcher frontend repo-updater symbols "

# GOMOD_ROOT is the directory from which `go install` commands are run. It should contain a go.mod
# file. The go.mod file may be updated as a side effect of updating the dependencies before the `go
# install`.
GOMOD_ROOT=${GOMOD_ROOT:-$PWD}
echo >&2 "Running \`go install\` from $GOMOD_ROOT"

# handle options
verbose=false
while getopts 'v' o; do
  case $o in
  v) verbose=true;;
  \?) echo >&2 "usage: go-install.sh [-v] [commands]"
      exit 1
      ;;
  esac
done
shift $(expr $OPTIND - 1)

# check provided commands
ok=true
case $# in
0) commands=$all_oss_commands;;
*) commands=" $* "
  for cmd in $commands; do
    case $all_oss_commands in
    *" $cmd "*) ;;
    *) echo >&2 "Unknown command: $cmd"
       ok=false
       ;;
    esac
  done
  ;;
esac

$ok || exit 1

mkdir -p .bin
export GOBIN=$PWD/.bin
export GO111MODULE=on

INSTALL_GO_PKGS="github.com/mattn/goreman \
github.com/sourcegraph/docsite/cmd/docsite \
github.com/google/zoekt/cmd/zoekt-archive-index \
github.com/google/zoekt/cmd/zoekt-sourcegraph-indexserver \
github.com/google/zoekt/cmd/zoekt-webserver \
"