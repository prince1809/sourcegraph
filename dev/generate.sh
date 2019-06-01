#!/bin/bash

export PGPORT=5432
export PGHOST=localhost
export PGUSER=sourcegraph
export PGPASSWORD=sourcegraph
export PGDATABASE=sourcegraph
export PGSSLMODE=disable

cd "$(dirname "${BASH_SOURCE[0]}")/.." # cd to repo root dir

go list ./... | xargs go generate -x
