#!/bin/bash

# ============================================
# Clear Database and Seed Script
# ============================================
# This script clears all data from the database
# and then seeds it with fresh dummy data
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
echo -e "${GREEN}Clear Database and Seed Script${NC}"
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

# SQL files
CLEAR_FILE="$SCRIPT_DIR/clear_database.sql"
SEED_FILE="$SCRIPT_DIR/seed_dummy_data.sql"

# Check if SQL files exist
if [ ! -f "$CLEAR_FILE" ]; then
    echo -e "${RED}Error: Clear SQL file not found at $CLEAR_FILE${NC}"
    exit 1
fi

if [ ! -f "$SEED_FILE" ]; then
    echo -e "${RED}Error: Seed SQL file not found at $SEED_FILE${NC}"
    exit 1
fi

echo -e "${YELLOW}Step 1: Clearing all database data...${NC}"
if docker-compose -f "$PROJECT_ROOT/docker-compose.yaml" exec -T db psql -U postgres -d b2b < "$CLEAR_FILE"; then
    echo -e "${GREEN}✓ Database cleared successfully!${NC}"
else
    echo -e "${RED}✗ Error clearing database.${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}Step 2: Seeding database with dummy data...${NC}"
if docker-compose -f "$PROJECT_ROOT/docker-compose.yaml" exec -T db psql -U postgres -d b2b < "$SEED_FILE"; then
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
    echo -e "${GREEN}Done! Database has been cleared and reseeded.${NC}"
else
    echo ""
    echo -e "${RED}✗ Error seeding database. Please check the error messages above.${NC}"
    exit 1
fi

