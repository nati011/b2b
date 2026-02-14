#!/bin/bash

# ============================================
# Fix Docker Permission Issues
# ============================================
# This script helps resolve Docker permission denied errors
# ============================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Checking Docker service status...${NC}"

# Check if Docker service is running
if ! systemctl is-active --quiet docker; then
    echo -e "${RED}Docker service is not running!${NC}"
    echo -e "${YELLOW}Attempting to restart Docker service...${NC}"
    echo ""
    echo "Please run the following command with sudo:"
    echo -e "${GREEN}sudo systemctl restart docker${NC}"
    echo ""
    echo "Or if you have sudo access, you can run:"
    echo -e "${GREEN}sudo systemctl restart docker && newgrp docker${NC}"
    exit 1
fi

echo -e "${GREEN}Docker service is running.${NC}"

# Check if user is in docker group
if groups | grep -q docker; then
    echo -e "${GREEN}User is in docker group.${NC}"
else
    echo -e "${YELLOW}User is not in docker group.${NC}"
    echo "Adding user to docker group..."
    echo "Please run: sudo usermod -aG docker $USER"
    echo "Then log out and log back in, or run: newgrp docker"
    exit 1
fi

# Check Docker socket permissions
if [ -S /var/run/docker.sock ]; then
    PERMS=$(stat -c "%a" /var/run/docker.sock)
    OWNER=$(stat -c "%U:%G" /var/run/docker.sock)
    echo -e "${GREEN}Docker socket permissions: $PERMS ($OWNER)${NC}"
    
    if [ "$PERMS" != "660" ]; then
        echo -e "${YELLOW}Warning: Docker socket permissions are not 660${NC}"
    fi
else
    echo -e "${RED}Docker socket not found!${NC}"
    exit 1
fi

# Test Docker access
if docker ps > /dev/null 2>&1; then
    echo -e "${GREEN}Docker access test: SUCCESS${NC}"
else
    echo -e "${RED}Docker access test: FAILED${NC}"
    echo ""
    echo "Try running: newgrp docker"
    echo "Or log out and log back in to refresh group membership"
    exit 1
fi

echo ""
echo -e "${GREEN}All checks passed! Docker should be working correctly.${NC}"

