#!/usr/bin/env sh
set -eu
BASE_URL="${BASE_URL:?set BASE_URL}"

for path in "/api/v1/../../etc/passwd" "/api/v1/%00" "/api/v1/scans?x=' OR 1=1 --"; do
  code=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL$path")
  echo "$path => $code"
done
