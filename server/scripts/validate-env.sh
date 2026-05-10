#!/usr/bin/env sh
set -eu

required_vars="SERVER_ENV SERVER_PORT DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME DB_SSLMODE REDIS_ADDR JWT_SECRET AWS_REGION S3_BUCKET GOOGLE_CLIENT_ID RAZORPAY_KEY_ID RAZORPAY_KEY_SECRET RAZORPAY_WEBHOOK_SECRET GEMINI_API_KEY OPENAI_API_KEY"
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
