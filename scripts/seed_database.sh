#!/bin/bash

# ============================================
# Database Seeding Script
# ============================================
# This script seeds the database with dummy data for testing
# Usage: ./scripts/seed_database.sh
# ============================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Database Seeding Script${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running. Please start Docker and try again.${NC}"
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: docker-compose is not installed.${NC}"
    exit 1
fi

# Check if the database container is running
if ! docker-compose -f "$PROJECT_ROOT/docker-compose.yaml" ps db | grep -q "Up"; then
    echo -e "${YELLOW}Warning: Database container is not running. Starting it...${NC}"
    docker-compose -f "$PROJECT_ROOT/docker-compose.yaml" up -d db
    echo "Waiting for database to be ready..."
    sleep 5
fi

# SQL file path
SQL_FILE="$SCRIPT_DIR/seed_dummy_data.sql"

# Check if SQL file exists
if [ ! -f "$SQL_FILE" ]; then
    echo -e "${RED}Error: SQL file not found at $SQL_FILE${NC}"
    exit 1
fi

echo -e "${YELLOW}Seeding database with dummy data...${NC}"
echo ""

# Run the SQL file
if docker-compose -f "$PROJECT_ROOT/docker-compose.yaml" exec -T db psql -U postgres -d b2b < "$SQL_FILE"; then
    echo ""
    echo -e "${GREEN}✓ Database seeded successfully!${NC}"
    echo ""
    echo -e "${GREEN}Verification:${NC}"
    
    # Run verification queries
    docker-compose -f "$PROJECT_ROOT/docker-compose.yaml" exec -T db psql -U postgres -d b2b -c "
        SELECT 
            (SELECT COUNT(*) FROM public.category WHERE is_deleted = false) as categories,
            (SELECT COUNT(*) FROM public.products WHERE is_deleted = false AND is_active = true) as products,
            (SELECT COUNT(*) FROM public.configurable_products WHERE is_deleted = false AND is_available = true) as configurable_products;
    "
    
    echo ""
    echo -e "${GREEN}Done! You can now see products in the application.${NC}"
else
    echo ""
    echo -e "${RED}✗ Error seeding database. Please check the error messages above.${NC}"
    exit 1
fi




