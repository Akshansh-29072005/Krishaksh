#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}ERROR: DATABASE_URL environment variable is not set${NC}"
    echo "Set it with: export DATABASE_URL='postgresql://user:password@host:port/database'"
    exit 1
fi

MIGRATIONS_DIR="./server/db/migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}ERROR: Migrations directory not found at $MIGRATIONS_DIR${NC}"
    exit 1
fi

echo -e "${YELLOW}Starting database migrations...${NC}"
echo -e "${YELLOW}Database: $(echo $DATABASE_URL | sed 's/.*@//' | sed 's/:.*@//')${NC}"
echo ""

# Array of migrations to run
MIGRATIONS=(
    "000001_init_schema"
    "000002_commerce_payments"
    "000003_observability_analytics_admin"
    "000004_scan_ai_inference"
    "000005_crop_and_app_config"
    "000006_add_support_tickets_description"
    "000007_sync_support_tickets_schema"
)

connect_with_retry() {
    local retries=5
    local delay=5
    local count=0

    while true; do
        if PGSSLMODE=require psql "$DATABASE_URL" --set=ON_ERROR_STOP=on -Atc "$1" >/dev/null 2>&1; then
            return 0
        fi

        count=$((count + 1))
        if [ $count -ge $retries ]; then
            return 1
        fi

        echo -e "${YELLOW}Connection attempt $count failed. Retrying in ${delay}s...${NC}"
        sleep $delay
    done
}

# Run each migration
for migration in "${MIGRATIONS[@]}"; do
    UP_FILE="$MIGRATIONS_DIR/${migration}.up.sql"
    
    if [ ! -f "$UP_FILE" ]; then
        echo -e "${RED}✗ Migration file not found: $UP_FILE${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}Running: $migration${NC}"

    if connect_with_retry "SELECT 1"; then
        if PGSSLMODE=require psql "$DATABASE_URL" --set=ON_ERROR_STOP=on -f "$UP_FILE"; then
            echo -e "${GREEN}✓ Successfully applied: $migration${NC}"
        else
            echo -e "${RED}✗ Failed to apply: $migration${NC}"
            echo "Error details:"
            PGSSLMODE=require psql "$DATABASE_URL" --set=ON_ERROR_STOP=on -f "$UP_FILE"
            exit 1
        fi
    else
        echo -e "${RED}✗ Could not establish a database connection after retries.${NC}"
        exit 1
    fi
    
    echo ""
done

echo -e "${GREEN}All migrations completed successfully!${NC}"
