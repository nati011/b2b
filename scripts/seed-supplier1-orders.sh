#!/bin/bash

# Script to seed test orders for supplier1@b2b.local
# Usage: ./scripts/seed-supplier1-orders.sh

set -e

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-b2b}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

SEED_FILE="db/seed/011_seed_supplier1_orders.sql"

echo "Seeding test orders for supplier1@b2b.local"
echo "============================================"
echo "Database: $DB_NAME on $DB_HOST:$DB_PORT"
echo ""

if [ ! -f "$SEED_FILE" ]; then
    echo "✗ Seed file not found: $SEED_FILE"
    exit 1
fi

echo "Running seed file: $SEED_FILE"
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$SEED_FILE"

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ Seed file executed successfully!"
    echo ""
    echo "You can now test the supplier orders endpoint:"
    echo "  curl -u supplier1@b2b.local:password123 http://localhost:8090/orders/supplier"
    echo ""
    echo "Or run the test script:"
    echo "  ./scripts/test-supplier-orders.sh"
else
    echo ""
    echo "✗ Failed to execute seed file"
    exit 1
fi



