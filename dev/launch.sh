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

export LIGHTSTEP_INCLUDE_SENSITIVE=true
export PGSSLMODE=disable

# Default to "info" level debugging, and "condensed" log format (nice for human readers)
export SRC_LOG_LEVEL=${SRC_LOG_LEVEL:-info}
export SRC_LOG_FORMAT=${SRC_LOG_FORMAT:-condensed}
export GITHUB_BASE_URL=http://127.0.0.1:3180
export SRC_REPOS_DIR=$HOME/.sourcegraph/repos
export INSECURE_DEV=1
export GOLANGSERVER_SRC_GIT_SERVERS=host.docker.internal:3178


# webpack-dev-server is a proxy running on port 3080 that (1) serves assets, waiting to respond
# until they are built and (2) otherwise proxies to nginx running on port 3081 (which proxies to
# sourcegraph running on port 3082). That is why Sourcegraph listens on 3081 despite the externalURL
# having port 3080.
export SRC_HTTP_ADDR=":3082"
export WEBPACK_DEV_SERVER=1

# WebApp
export NODE_ENV=development
export NODE_OPTIONS="--max_old_space_size=4096"

# Make sure chokidar-cli is installed in the background
printf >&2 "Currently installing Yarn and Go dependencies...\n\n"
yarn_pid=''
[ -n "${OFFLINE-}"] || {
  yarn --no-progress &
  yarn_pid="$!"
}

if ! ./dev/go-install.sh; then
  # let yarn finish, otherwise we get Yarn diagnostics AFTER the
  # actual reason we're failing.
  wait
  echo >&2 "WARNING: go-install.sh failed, some builds may have failed."
  exit 1
fi