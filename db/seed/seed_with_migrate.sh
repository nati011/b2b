#!/bin/bash

# Script to seed database using go-migrate database connection format
# Usage: ./db/seed/seed_with_migrate.sh
# 
# This script uses go-migrate's database URL format and connection approach
# to execute seed SQL files in the correct order.

set -e

# Default database connection values (can be overridden via environment variables)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-b2b}"
DB_SSLMODE="${DB_SSLMODE:-disable}"

# Construct database URL in go-migrate format
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SEED_DIR="$SCRIPT_DIR"

echo "=========================================="
echo "Seeding database using go-migrate format"
echo "=========================================="
echo "Database: ${DB_NAME}@${DB_HOST}:${DB_PORT}"
echo "Database URL: postgres://${DB_USER}:***@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
echo "Seed directory: ${SEED_DIR}"
echo ""

# Verify go-migrate tool is installed (optional check)
if command -v migrate >/dev/null 2>&1; then
  MIGRATE_VERSION=$(migrate -version 2>&1 || echo "unknown")
  echo "Using golang-migrate: $MIGRATE_VERSION"
  echo ""
fi

# Wait for database server to be ready
echo "Waiting for database server to be ready..."
counter=0
max_attempts=60
while [ $counter -lt $max_attempts ]; do
  if PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "postgres" -c '\q' 2>/dev/null; then
    echo "✓ Database server is ready"
    break
  fi
  counter=$((counter + 1))
  if [ $counter -eq $max_attempts ]; then
    echo "✗ Database connection timeout after ${max_attempts} attempts"
    exit 1
  fi
  echo "Database is unavailable - sleeping (attempt ${counter}/${max_attempts})"
  sleep 1
done

# Check if database exists
echo "Checking if database '${DB_NAME}' exists..."
if ! PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -c '\q' 2>/dev/null; then
  echo "Database '${DB_NAME}' does not exist, creating it..."
  PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "postgres" -c "CREATE DATABASE ${DB_NAME};" 2>/dev/null || {
    echo "  (Database may already exist or creation failed)"
  }
fi
echo "✓ Database '${DB_NAME}' is ready"
echo ""

# Array of seed files in execution order
# Note: Roles are configured via config/roles.yaml and bootstrapped by the application
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

# Execute each seed file using psql with go-migrate's connection format
echo "Starting seed data execution..."
echo "----------------------------------------"

for file in "${SEED_FILES[@]}"; do
  filepath="${SEED_DIR}/${file}"
  if [ -f "$filepath" ]; then
    echo "Seeding ${file}..."
    if PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -f "$filepath" >/dev/null 2>&1; then
      echo "✓ Completed ${file}"
    else
      # Check if it's just a conflict (already seeded)
      if PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -f "$filepath" 2>&1 | grep -q "already exists\|duplicate\|conflict"; then
        echo "  (Already seeded or skipped: ${file})"
      else
        echo "✗ Error seeding ${file}"
        exit 1
      fi
    fi
  else
    echo "⚠ Warning: Seed file not found: ${filepath}"
    echo "  Skipping..."
  fi
done

echo "----------------------------------------"
echo "=========================================="
echo "✓ All seed data loaded successfully!"
echo "=========================================="




