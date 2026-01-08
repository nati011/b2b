#!/bin/bash
# Force stop the client container
echo "Stopping client container..."
docker compose stop client --timeout 30
sleep 2
# If still running, force kill
if docker ps | grep -q "b2b-client"; then
    echo "Force killing client container..."
    docker compose kill client
    docker compose rm -f client
fi
echo "Client container stopped."
