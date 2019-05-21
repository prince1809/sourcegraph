#!/bin/bash

cd $(dirname "${BASH_SOURCE[0]}")/../migrations
set -e

# The name is intentionally empty ('') so that it forces a merge conflict if two branches attempt to
# create a migration at the same sequence number (because they will both add a file with the same
# name, like `migration/1528277032_.up.sql`)

migrate create -ext sql -dir . -digits 10 -seq ''
