#!/bin/bash

# Script to properly stop client and admin containers
# This script uses docker-compose stop which respects restart policies

set -e

echo "Stopping client and admin containers..."

# Stop containers gracefully using docker-compose
docker-compose stop client admin

echo "✅ Containers stopped successfully"

# Optional: Show container status
echo ""
echo "Container status:"
docker-compose ps client admin

