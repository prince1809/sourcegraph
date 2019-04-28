#!/usr/bin/env bash

# set to true if unset so set -u won't break us
: ${SOURCEGRAPH_COMBINE_CONFIG:=false}

set -euf -o pipefail

unset CDPATH
cd "$(dirname "${BASH_SOURCE[0]}")/.." # cd ro repo root dir

if [ -f .env ]; then
set -o allexport
source .env
set +o allexport
fi

export GO111MODULE=on