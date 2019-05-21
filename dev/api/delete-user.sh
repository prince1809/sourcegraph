#!/bin/bash

set -e
unset CDPATH
cd "$(dirname "${BASH_SOURCE[0]}")/../.."

eval $(grep 'export OVERRIDE_AUTH_SECRET=' dev/launch.sh)
