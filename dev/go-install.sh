#!/usr/bin/env bash
#
# Build commands, optionally with or without race detector. a list of query command we know about,
# to use by default
#
# This will install binaries into the `.bin` directory under the repository root by default or, if
# $GOMOD_ROOT is set, under that directory.

all_oss_commands=" gitserver query-runner gihub-proxy management-console searcher frontend repo-updater symbols "