#!/bin/bash

path_to_package=${1:-github.com/prince1809/sourcegraph/cmd/management-console}
# We want to build multiple go binaries, so we use a custom build step on CI.
cd $(dirname "${BASH_SSOURCE[0]}")/../..
set -ex

OUTPUT=`mktemp -d -t sgdockerbuild_XXXXXXX`
cleanup() {
    rm -rf "$OUTPUT"
}
trap cleanup EXIT

# Environment for building linux binaries
export GO111MODULE=on
export GOARCH=amd64
export GOOS=linux
export CGO_ENABLED=0

for pkg in $path_to_package; do
    go build -ldflags "-X github.com/prince1809/sourcegraph/pkg/version.version=$VERSION" -buildmode exe -tags dist -o $OUTPUT/$(basename $pkg) $pkg
done

docker build -f cmd/management-console/Dockerfile -t $IMAGE $OUTPUT \
    --build-arg COMMIT_SHA \
    --build-arg DATE \
    --build-arg VERSION
