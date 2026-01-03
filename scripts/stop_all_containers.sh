#!/bin/bash

# ============================================
# Stop All Docker Containers
# ============================================
# This script stops all running Docker containers
# ============================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Stopping all Docker containers...${NC}"

# Method 1: Stop all containers using docker stop
if docker stop $(docker ps -q) 2>/dev/null; then
    echo -e "${GREEN}✓ All containers stopped successfully${NC}"
    exit 0
fi

# Method 2: If that fails, try docker-compose down
echo -e "${YELLOW}Trying docker-compose down...${NC}"
if docker-compose down 2>/dev/null; then
    echo -e "${GREEN}✓ All containers stopped successfully${NC}"
    exit 0
fi

# Method 3: Force kill if needed (requires sudo)
echo -e "${RED}Standard stop failed. Trying force kill...${NC}"
echo "You may need to run with sudo:"
echo "sudo docker kill \$(sudo docker ps -q)"

# List running containers
echo ""
echo "Currently running containers:"
docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Status}}"

