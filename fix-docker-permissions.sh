#!/bin/bash
# Script to fix Docker container permission issues

echo "Attempting to stop containers with sudo..."

# Try to stop with sudo
sudo docker compose stop admin client 2>&1

# If that doesn't work, try kill
if [ $? -ne 0 ]; then
    echo "Stop failed, trying kill..."
    sudo docker compose kill admin client 2>&1
fi

# Force remove if still running
echo "Cleaning up..."
sudo docker compose rm -f admin client 2>&1

echo "Done. You can now restart with: docker compose up -d"

