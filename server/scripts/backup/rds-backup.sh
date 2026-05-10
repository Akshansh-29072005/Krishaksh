#!/usr/bin/env sh
set -eu

: "${DB_HOST:?missing DB_HOST}"
: "${DB_PORT:?missing DB_PORT}"
: "${DB_USER:?missing DB_USER}"
: "${DB_NAME:?missing DB_NAME}"
: "${DB_PASSWORD:?missing DB_PASSWORD}"
: "${BACKUP_S3_URI:?missing BACKUP_S3_URI like s3://bucket/path}"

TS=$(date +%Y%m%d_%H%M%S)
OUT="/tmp/krisho_${DB_NAME}_${TS}.sql.gz"

export PGPASSWORD="$DB_PASSWORD"
pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" | gzip > "$OUT"
aws s3 cp "$OUT" "$BACKUP_S3_URI/"
rm -f "$OUT"

echo "Backup uploaded to $BACKUP_S3_URI"
