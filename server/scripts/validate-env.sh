#!/usr/bin/env sh
set -eu

# Check for new consolidated variable
if [ -z "${DATABASE_URL}" ]; then
    echo "Error: DATABASE_URL is not set"
    exit 1
fi

# Individual vars are now optional/fallback
required_vars="SERVER_ENV SERVER_PORT REDIS_ADDR JWT_SECRET AWS_REGION S3_BUCKET GOOGLE_CLIENT_ID RAZORPAY_KEY_ID RAZORPAY_KEY_SECRET RAZORPAY_WEBHOOK_SECRET GEMINI_API_KEY OPENAI_API_KEY"
for var in $required_vars; do
  eval "val=\${$var:-}"
  if [ -z "$val" ]; then
    echo "Missing required env: $var"
    exit 1
  fi
done

if [ "${SERVER_ENV}" = "production" ] && [ "${DB_SSLMODE}" = "disable" ]; then
  echo "DB_SSLMODE=disable is not allowed in production"
  exit 1
fi

echo "Environment validation passed"
