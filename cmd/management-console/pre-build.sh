#!/bin/bash

cd $(dirname "${BASH_SOURCE[0]}")
set -ex

go generate ./assets
