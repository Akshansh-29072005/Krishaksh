#!/usr/bin/env sh
set -eu
ENV_FILE="${1:-.env.staging}"
docker compose --env-file "$ENV_FILE" -f docker-compose.staging.yml down

echo "Rollback complete: stack stopped."
echo "Redeploy previous image tag by setting image tags in compose and rerunning deploy-staging.sh"
