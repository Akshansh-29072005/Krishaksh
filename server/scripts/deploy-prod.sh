#!/usr/bin/env sh
set -eu

ENV_FILE="${1:-.env.prod}"
if [ ! -f "$ENV_FILE" ]; then
  echo "Missing env file: $ENV_FILE"
  exit 1
fi

set -a
. "$ENV_FILE"
set +a

./scripts/validate-env.sh

echo "Pull/build latest images..."
docker compose --env-file "$ENV_FILE" -f docker-compose.prod.yml build --pull

echo "Starting production stack..."
docker compose --env-file "$ENV_FILE" -f docker-compose.prod.yml up -d --remove-orphans

echo "Waiting for readiness..."
sleep 10
curl -fsS https://${DOMAIN:-localhost}/health >/dev/null || curl -fsS http://localhost/health >/dev/null

echo "Deployment completed"
