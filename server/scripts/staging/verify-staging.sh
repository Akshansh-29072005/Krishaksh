#!/usr/bin/env sh
set -eu

ENV_FILE="${1:-.env.staging}"
set -a
. "$ENV_FILE"
set +a

BASE_URL="${BASE_URL:-https://${DOMAIN}}"

echo "Checking health endpoints..."
curl -fsS "${BASE_URL}/health" >/dev/null
curl -fsS "${BASE_URL}/ready" >/dev/null

echo "Checking metrics integrity..."
METRICS=$(curl -fsS "${BASE_URL}/metrics")
printf "%s" "$METRICS" | grep -q "http_request_duration_ms_total"
printf "%s" "$METRICS" | grep -q "ai_inference"

echo "Checking container health..."
docker compose --env-file "$ENV_FILE" -f docker-compose.staging.yml ps

echo "Staging verification passed"
