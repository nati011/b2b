#!/bin/bash

echo "Stopping all running containers..."
sudo docker stop $(sudo docker ps -q) 2>/dev/null || true

echo "Removing all containers..."
sudo docker rm -f $(sudo docker ps -aq) 2>/dev/null || true

echo "Verifying all containers are removed..."
sudo docker ps -a

echo "Done! All containers have been purged."

