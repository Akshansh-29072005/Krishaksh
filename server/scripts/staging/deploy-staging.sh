#!/usr/bin/env sh
set -eu

ENV_FILE="${1:-.env.staging}"
[ -f "$ENV_FILE" ] || { echo "Missing $ENV_FILE"; exit 1; }

set -a
. "$ENV_FILE"
set +a

./scripts/validate-env.sh
docker compose --env-file "$ENV_FILE" -f docker-compose.staging.yml build --pull
docker compose --env-file "$ENV_FILE" -f docker-compose.staging.yml up -d --remove-orphans

./scripts/staging/verify-staging.sh "$ENV_FILE"
