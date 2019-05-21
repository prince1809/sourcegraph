#!/bin/bash

psql -C "drop schema public cascade; create schema public;"
redis-cli -c flushall
