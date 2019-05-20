#!/usr/bin/env bash

# We want to build multiple go binaries, so we use a custom build step on CI.
cd "$(dirname "${BASH_SOURCE[0]}")/../.."
set -eux

OUTPUT=$(mktemp -d -t sgserver_XXXXXXX)
cleanup() {
  rm -rf "$OUTPUT"
}
trap cleanup EXIT

# Environment for building linux binaries
export GO111MODULE=on
export GOARCH=amd64
export GOOS=linux
export CGO_ENABLED=0

# Additional images passed in here when this script is called externally by our
# enterprise build scripts.
additional_images=${@:-github.com/prince1809/sourcegraph/cmd/frontend github.com/prince1809/sourcegraph/cmd/management-console}

# Overridable server packages path for when this script is called externally by
# our enterprise build scripts.
server_pkg=${SERVER_PKG:-github.com/prince1809/sourcegraph/cmd/server}

cp -a ./cmd/server/rootfs/. "$OUTPUT"
bindir="$OUTPUT/usr/local/bin"
mkdir -p "$bindir"

echo "--- go build"
for pkg in $server_pkg \
    github.com/prince1809/sourcegraph/cmd/github-proxy \
    github.com/prince1809/sourcegraph/cmd/gitserver \
    github.com/prince1809/sourcegraph/cmd/query-runner \
    github.com/prince1809/sourcegraph/cmd/repo-updater \
    github.com/prince1809/sourcegraph/cmd/searcher \
    github.com/prince1809/sourcegraph/cmd/zoekt-archive-index \
    github.com/prince1809/sourcegraph/cmd/zoekt-sourcegraph-indexserver \
    github.com/prince1809/sourcegraph/cmd/zoekt-webserver $additional_images; do

    go build \
        -a \
        -ldflags "-X github.com/prince1809/sourcegraph/pkg/version.version=$VERSION" \
        -buildmode exe \
        -installsuffix netgo \
        -tags "dis netgo" \
        -o "$bindir/$(basename "$pkg")" "$pkg"
done

echo "--- build sqlite for symbols"
env CTAGS_D_OUTPUT_PATH="$OUTPUT/.ctags.d" SYMBOLS_EXECUTABLE_OUTPUT_PATH="$bindir/symbols" BUILD_TYPE=dis ./cmd/symbols/build.sh buildSymbolsDockerImageDependencies

echo "--- docker build"
docker build -f cmd/server/Dockerfile -t "$IMAGE" "$OUTPUT" \
    --build-arg COMMIT_SHA \
    --build-arg DATE \
    --build-arg VERSION
