#!/bin/bash

cd "$(dirname "${BASH_SOURCE[0]}")/.." # cd to repo root dir
GO_DIRS="$(./dev/watchdirs.sh) ${ADDITIONAL_GO_DIRS}"

dirs_starstar() {
    for i; do echo "'$/**/*.go'"; done
}
