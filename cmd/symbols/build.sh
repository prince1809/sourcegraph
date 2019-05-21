#!/bin/bash

set -e
IMAGE="${IMAGE:=dev-symbols}"
CTAGS_IMAGE="${CTAGS_IMAGE:=ctags}"
BUILD_TYPE="${BUILD_TYPE:=dev}"

repositoryRoot="$PWD"
SYMBOLS_EXECUTABLE_OUTPUT_PATH="${SYMBOLS_EXECUTABLE_OUTPUT_PATH:=$repositoryRoot/.bin/symbols}"
case "$OSTYPE" in
    darwin*)
        libsqlite3PcrePath="$repositoryRoot/libsqlite3-pcre.dylib"
        ;;
    linux*)
        libsqlite3PcrePath="$repositoryRoot/libsqlite3-pcre.so"
        ;;
    *)
        echo "Unknown platform $OSTYPE"
        exit 1
        ;;
esac


