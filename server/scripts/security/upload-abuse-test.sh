#!/usr/bin/env sh
set -eu
BASE_URL="${BASE_URL:?set BASE_URL}"
JWT_TOKEN="${JWT_TOKEN:?set JWT_TOKEN}"

# Oversized body test (11MB)
BIG=$(mktemp)
dd if=/dev/zero of="$BIG" bs=1M count=11 >/dev/null 2>&1

status=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/scans" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  --data-binary "$(printf '{\"crop_type\":\"Wheat\",\"image_key\":\"%s\"}' "$(cat "$BIG" | tr -d '\n' | head -c 1024)")")

echo "upload_abuse_status=$status (expect 400/413)"
rm -f "$BIG"
