#!/bin/bash

# Script to seed all development data
# Usage: ./db/seed/seed_all.sh

set -e

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-b2b}"
DB_USER="${DB_USER:-postgres}"

echo "Seeding database: $DB_NAME on $DB_HOST:$DB_PORT"
echo "================================================"

SEED_DIR="$(dirname "$0")"

# Array of seed files in order
# Note: Roles are configured via config/roles.yaml and bootstrapped by the application
# Note: 009_seed_user_roles.sql must be run after the application has bootstrapped roles
SEED_FILES=(
  "001_seed_users.sql"
  "008_seed_credentials.sql"  # Credentials must be seeded after users
  "009_seed_user_roles.sql"   # User roles must be seeded after users and roles are bootstrapped
  "002_seed_suppliers.sql"
  "003_seed_customers.sql"
  "004_seed_categories.sql"
  "005_seed_products.sql"
  "006_seed_product_categories.sql"
  "007_seed_orders.sql"
)

# Execute each seed file
for file in "${SEED_FILES[@]}"; do
  filepath="$SEED_DIR/$file"
  if [ -f "$filepath" ]; then
    echo "Seeding $file..."
    PGPASSWORD="${DB_PASSWORD:-postgres}" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$filepath"
    echo "✓ Completed $file"
  else
    echo "✗ File not found: $filepath"
    exit 1
  fi
done

echo "================================================"
echo "✓ All seed data loaded successfully!"

