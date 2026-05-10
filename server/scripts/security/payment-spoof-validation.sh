#!/usr/bin/env sh
set -eu
BASE_URL="${BASE_URL:?set BASE_URL}"

payload='{"id":"evt_spoof","event":"payment.captured","payload":{"payment":{"entity":{"id":"pay_spoof","order_id":"order_spoof","amount":100,"status":"captured"}}}}'
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/payments/webhook" \
  -H 'Content-Type: application/json' \
  -H 'X-Razorpay-Signature: invalid' \
  -d "$payload")

echo "spoof_webhook_status=$code (expect 401)"
