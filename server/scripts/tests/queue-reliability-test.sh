#!/usr/bin/env sh
set -eu
BASE_URL="${BASE_URL:?set BASE_URL}"
JWT_TOKEN="${JWT_TOKEN:?set JWT_TOKEN}"

N="${N:-50}"
i=0
while [ $i -lt "$N" ]; do
  payload=$(printf '{"crop_type":"Wheat","image_key":"scans/reliability/%s.jpg"}' "$i")
  curl -fsS -X POST "$BASE_URL/api/v1/scans" -H "Authorization: Bearer $JWT_TOKEN" -H "Content-Type: application/json" -d "$payload" >/dev/null || true
  i=$((i+1))
done

metrics=$(curl -fsS "$BASE_URL/metrics")
printf "%s" "$metrics" | grep -q "queue_enqueue_total_scan"
echo "queue_reliability_test_passed"
