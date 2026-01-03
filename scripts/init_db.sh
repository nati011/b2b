#!/bin/bash

# ============================================
# Database Initialization Script
# ============================================
# This script runs automatically when the database container starts
# It clears all data and seeds fresh data
# ============================================

set -e

# Wait for PostgreSQL to be ready
until pg_isready -U "${POSTGRES_USER:-postgres}" -d "${POSTGRES_DB:-b2b}"; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 1
done

echo "PostgreSQL is ready. Clearing and seeding database..."

# Run clear script
psql -U "${POSTGRES_USER:-postgres}" -d "${POSTGRES_DB:-b2b}" -f /docker-entrypoint-initdb.d/clear_database.sql

# Run seed script
psql -U "${POSTGRES_USER:-postgres}" -d "${POSTGRES_DB:-b2b}" -f /docker-entrypoint-initdb.d/seed_dummy_data.sql

echo "Database initialization complete!"

