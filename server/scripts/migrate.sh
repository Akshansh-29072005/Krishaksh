#!/usr/bin/env sh
set -eu

if [ ! -f ./krisho-api ]; then
  echo "krisho-api binary not found, building..."
  go build -o krisho-api ./cmd/api
fi

RUN_DB_MIGRATIONS=true ENABLE_EMBEDDED_WORKER=false ./krisho-api
