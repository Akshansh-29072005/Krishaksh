#!/usr/bin/env sh
set -eu
BASE_URL="${BASE_URL:?set BASE_URL}"

echo "Testing invalid webhook signature rejection..."
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/payments/webhook" -H 'Content-Type: application/json' -H 'X-Razorpay-Signature: invalid' -d '{"id":"evt_stage_test","event":"payment.captured","payload":{"payment":{"entity":{"id":"pay_test","order_id":"order_test","amount":100,"status":"captured"}}}}')
[ "$code" = "401" ] || { echo "expected 401 got $code"; exit 1; }
echo "payment_webhook_test_passed"
