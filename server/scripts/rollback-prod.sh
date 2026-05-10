#!/usr/bin/env sh
set -eu

ENV_FILE="${1:-.env.prod}"
if [ ! -f "$ENV_FILE" ]; then
  echo "Missing env file: $ENV_FILE"
  exit 1
fi

echo "Rollback strategy: redeploy previous image tag by setting IMAGE_TAG in env and re-running deploy."
echo "Stopping current stack..."
docker compose --env-file "$ENV_FILE" -f docker-compose.prod.yml down

echo "Re-run deploy with previous image tag env values."
