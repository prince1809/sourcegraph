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
go run ./pkg/version/minversion

export GOMOD_ROOT="${GOMOD_ROOT:-$PWD}"

# Verify postgresql config.
hash psql 2>/dev/null || {
  # "brew install postgresql@9.6" does not put psql on the $PATH by default;
  # try to fix this automatically if we can.
  hash brew 2>/dev/null && {
    if [[ -x "$(brew --prefix)/opt/postgresql@9.6/bin/psql" ]]; then
      export PATH="$(brew --prefix)/opt/postgresql@9.6/bin:$PATH"
    fi
  }
}

if ! psql -wc '\x' >/dev/null; then
  echo "FAIL: postgreSQL config invalid or missing or postgreSQL is still starting up"
  echo "You probably need, at least, PGUSER and PGPASSWD set in the environment"
  exit 1
fi