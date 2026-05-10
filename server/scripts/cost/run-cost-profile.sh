#!/usr/bin/env sh
set -eu
INPUT="${1:-deploy/staging/cost-input.sample.json}"
go run ./cmd/costprofiler -in "$INPUT"
