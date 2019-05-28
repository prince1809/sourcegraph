#!/bin/bash

yarn run serve 2>&1 | grep -v 'Could not load existing sourcemap of'
