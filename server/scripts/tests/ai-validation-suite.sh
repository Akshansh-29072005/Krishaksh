#!/usr/bin/env sh
set -eu

BASE_URL="${BASE_URL:?set BASE_URL}"
JWT_TOKEN="${JWT_TOKEN:?set JWT_TOKEN}"
CROP_TYPE="${CROP_TYPE:-Wheat}"

# Input file format: one S3 image_key per line already uploaded in staging bucket.
IMAGE_KEYS_FILE="${IMAGE_KEYS_FILE:-./deploy/staging/sample-image-keys.txt}"
[ -f "$IMAGE_KEYS_FILE" ] || { echo "Missing $IMAGE_KEYS_FILE"; exit 1; }

OUT_DIR="${OUT_DIR:-./deploy/staging/ai-validation-output}"
mkdir -p "$OUT_DIR"
RAW="$OUT_DIR/raw_results.jsonl"
: > "$RAW"

while IFS= read -r key; do
  [ -z "$key" ] && continue
  payload=$(printf '{"crop_type":"%s","image_key":"%s"}' "$CROP_TYPE" "$key")
  create=$(curl -fsS -X POST "$BASE_URL/api/v1/scans" -H "Authorization: Bearer $JWT_TOKEN" -H "Content-Type: application/json" -d "$payload")
  scan_id=$(printf "%s" "$create" | jq -r '.data.id')

  i=0
  status=""
  confidence=""
  provider=""
  while [ $i -lt 20 ]; do
    s=$(curl -fsS -H "Authorization: Bearer $JWT_TOKEN" "$BASE_URL/api/v1/scans/$scan_id")
    status=$(printf "%s" "$s" | jq -r '.data.prediction_status')
    confidence=$(printf "%s" "$s" | jq -r '.data.confidence_score // empty')
    provider=$(printf "%s" "$s" | jq -r '.data.ai_provider // empty')
    if [ "$status" = "COMPLETED" ] || [ "$status" = "UNCERTAIN" ] || [ "$status" = "FAILED" ]; then
      break
    fi
    i=$((i+1))
    sleep 2
  done

  printf '{"scan_id":"%s","status":"%s","confidence":"%s","provider":"%s"}\n' "$scan_id" "$status" "$confidence" "$provider" >> "$RAW"
done < "$IMAGE_KEYS_FILE"

uncertain_count=$(jq -s '[.[] | select(.status=="UNCERTAIN")] | length' "$RAW")
failed_count=$(jq -s '[.[] | select(.status=="FAILED")] | length' "$RAW")
low_conf_count=$(jq -s '[.[] | select((.confidence|tonumber?) < 0.65)] | length' "$RAW")
provider_fallback_proxy=$(jq -s '[.[] | select(.provider=="openai")] | length' "$RAW")

echo "AI validation summary"
echo "uncertain_predictions=$uncertain_count"
echo "failed_predictions=$failed_count"
echo "low_confidence_lt_0_65=$low_conf_count"
echo "fallback_provider_openai_count=$provider_fallback_proxy"
