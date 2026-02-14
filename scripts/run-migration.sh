#!/bin/bash

# Script to run database migrations
# Usage: ./scripts/run-migration.sh [migration_number]
# If no migration number is provided, runs all pending migrations

set -e

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Default values from config.yaml
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-b2b}"
DB_SSLMODE="${DB_SSLMODE:-disable}"

# Construct database URL
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

# Migrations directory
MIGRATIONS_DIR="${PROJECT_ROOT}/db/migrations"

echo "Running database migrations..."
echo "Database: ${DB_NAME}@${DB_HOST}:${DB_PORT}"
echo "Migrations directory: ${MIGRATIONS_DIR}"
echo ""

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Error: Migrations directory not found at $MIGRATIONS_DIR"
    exit 1
fi

# Run migrations
if [ -n "$1" ]; then
    # Run specific migration
    echo "Running migration: $1"
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" up "$1"
else
    # Run all pending migrations
    echo "Running all pending migrations..."
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" up
fi

echo ""
echo "Migration completed successfully!"

