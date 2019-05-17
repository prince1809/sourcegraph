#!/bin/bash

cd $(dirname "${BASH_SOURCE[0]}")/../..

set -ex

./cmd/frontend/pre-build.sh
