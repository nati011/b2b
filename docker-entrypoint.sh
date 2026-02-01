#!/bin/sh
set -e

echo "=========================================="
echo "Starting B2B Marketplace Application"
echo "=========================================="
echo "Entrypoint script is running..."
echo ""

# Load database config from the config file
CONFIG_FILE="${CONFIG_FILE:-/root/config/config.docker.yaml}"

# Extract database connection details from environment or defaults
# These can be overridden via docker-compose environment variables
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-b2b}"
DB_SSLMODE="${DB_SSLMODE:-disable}"

# Construct database URL
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "Database: ${DB_NAME}@${DB_HOST}:${DB_PORT}"
echo ""

# Wait for database server to be ready (check postgres database first)
echo "Waiting for database server to be ready..."
until PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "postgres" -c '\q' 2>/dev/null; do
  echo "Database server is unavailable - sleeping"
  sleep 1
done
echo "✓ Database server is ready"

# Check if database exists, create if it doesn't
echo "Checking if database '${DB_NAME}' exists..."
if ! PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -c '\q' 2>/dev/null; then
  echo "Database '${DB_NAME}' does not exist, creating it..."
  PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "postgres" -c "CREATE DATABASE ${DB_NAME};" 2>/dev/null || {
    echo "  (Database may already exist or creation failed)"
  }
fi
echo "✓ Database '${DB_NAME}' is ready"
echo ""

# Run migrations using golang-migrate
echo "Running database migrations with golang-migrate..."
MIGRATIONS_DIR="/root/db/migrations"

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
  echo "  ERROR: Migrations directory not found at $MIGRATIONS_DIR"
  exit 1
fi

# Verify golang-migrate tool is installed
if ! command -v migrate >/dev/null 2>&1; then
  echo "  ERROR: golang-migrate tool not found!"
  echo "  Please ensure the migrate binary is installed in the Docker image"
  echo "  Expected location: /usr/local/bin/migrate"
  exit 1
fi

# Display migrate tool version
MIGRATE_VERSION=$(migrate -version 2>&1 || echo "unknown")
echo "  Using golang-migrate: $MIGRATE_VERSION"
echo "  Migrations directory: $MIGRATIONS_DIR"
echo "  Database: $DB_NAME@$DB_HOST:$DB_PORT"
echo ""

# Run migrations
echo "  Executing: migrate -path $MIGRATIONS_DIR -database [DB_URL] up"
if migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" up; then
  echo "✓ Migrations completed successfully"
else
  MIGRATE_EXIT_CODE=$?
  if [ $MIGRATE_EXIT_CODE -eq 0 ]; then
    echo "✓ Migrations completed (already up to date)"
  else
    echo "  WARNING: Migration command exited with code $MIGRATE_EXIT_CODE"
    echo "  This may indicate migrations are already applied or there was an issue"
    echo "  Continuing with application startup..."
  fi
fi
echo ""

# Run seed data
echo "Seeding database..."
SEED_DIR="/root/db/seed"
if [ -d "$SEED_DIR" ]; then
  SEED_FILES=(
    "001_seed_users.sql"
    "002_seed_suppliers.sql"
    "003_seed_customers.sql"
    "004_seed_categories.sql"
    "005_seed_products.sql"
    "006_seed_product_categories.sql"
    "007_seed_orders.sql"
  )
  
  for file in "${SEED_FILES[@]}"; do
    filepath="$SEED_DIR/$file"
    if [ -f "$filepath" ]; then
      echo "  Seeding $file..."
      PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -f "$filepath" >/dev/null 2>&1 || {
        echo "    (Already seeded or skipped)"
      }
    fi
  done
  echo "✓ Seed data loaded"
else
  echo "Warning: Seed directory not found at $SEED_DIR"
fi
echo ""

# Start the application
echo "Starting application..."
echo "=========================================="
exec "$@"
