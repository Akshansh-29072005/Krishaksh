#!/usr/bin/env sh
set -eu

BASE_URL="${BASE_URL:-https://staging-api.krisho.example.com}"
JWT_TOKEN="${JWT_TOKEN:-}"
TEST_ORDER_ID="${TEST_ORDER_ID:-}"

export STAGING_BASE_URL="$BASE_URL"
export STAGING_JWT_TOKEN="$JWT_TOKEN"
export STAGING_TEST_ORDER_ID="$TEST_ORDER_ID"

go test ./tests/e2e -v

echo "Running k6 staging load suite..."
k6 run deploy/k6/staging-upload-concurrency.js -e BASE_URL="$BASE_URL" -e JWT_TOKEN="$JWT_TOKEN"
k6 run deploy/k6/staging-ai-queue-pressure.js -e BASE_URL="$BASE_URL" -e JWT_TOKEN="$JWT_TOKEN"
k6 run deploy/k6/staging-order-concurrency.js -e BASE_URL="$BASE_URL" -e JWT_TOKEN="$JWT_TOKEN"
k6 run deploy/k6/staging-webhook-concurrency.js -e BASE_URL="$BASE_URL"
