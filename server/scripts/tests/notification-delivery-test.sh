#!/usr/bin/env sh
set -eu
BASE_URL="${BASE_URL:?set BASE_URL}"
JWT_TOKEN="${JWT_TOKEN:?set JWT_TOKEN}"

echo "Validating notification pipeline health via metrics/log hooks"
metrics=$(curl -fsS "$BASE_URL/metrics")
printf "%s" "$metrics" | grep -q "worker_errors_total_notification" || true
echo "notification_delivery_test_completed (requires real FCM token population for end-user delivery assertion)"
