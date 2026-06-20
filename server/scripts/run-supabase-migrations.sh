#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

if [ -f .env ]; then
  echo "Loading environment from .env"
  set -a
  # shellcheck disable=SC1091
  . .env
  set +a
fi

if [ -z "$DATABASE_URL" ]; then
  echo "ERROR: DATABASE_URL environment variable is not set."
  echo "Set it in server/.env or export it before running this script."
  exit 1
fi

echo "Using DATABASE_URL=$DATABASE_URL"

echo "Running Supabase migrations using Go migrator..."
go run ./cmd/migrate

echo "Migration run complete. Applied migrations:"
PGSSLMODE=require psql "$DATABASE_URL" -Atc "SELECT filename FROM schema_migrations ORDER BY filename;"
