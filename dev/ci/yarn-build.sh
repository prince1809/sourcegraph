#!/bin/bash

set -e

echo 'ENTERPRISE='$ENTERPRISE
echo 'NODE_ENV='$NODE_ENV
echo "# Note: NODE_ENV only used for build command"

echo "--- yarn in root"
NODE_ENV= yarn --frozen-lockfile --network-timeout 60000

cd $1
echo "--- browserslist"
NODE_ENV= yarn -s run browserslist

echo "--- build"
yarn -s run build --color

if jq -e '.scripts.bundlesize' package.json > /dev/null; then
    echo "--- bundlesize"
    NODE_ENV= GITHUB_TOKEN= yarn -s run bundlesize
fi

