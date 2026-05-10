#!/usr/bin/env sh
set -eu
BASE_URL="${BASE_URL:?set BASE_URL}"

m=$(curl -fsS "$BASE_URL/metrics")
for metric in \
  http_request_duration_ms_total \
  errors_total \
  queue_enqueue_total_scan \
  ai_provider_fallback_total \
  ai_inference_failure_total
 do
  printf "%s" "$m" | grep -q "$metric" || { echo "missing metric $metric"; exit 1; }
 done

echo "observability metrics integrity passed"
