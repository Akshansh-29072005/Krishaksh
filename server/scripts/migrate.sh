#!/usr/bin/env sh
set -eu
RUN_DB_MIGRATIONS=true ENABLE_EMBEDDED_WORKER=false ./krisho-api
