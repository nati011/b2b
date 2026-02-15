#!/bin/bash

# Stop backend, frontend, and nginx (keeps db and pgadmin running).
# Use "docker compose down" to stop everything.

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"
COMPOSE_FILE="$PROJECT_ROOT/docker-compose.yml"

if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
else
    COMPOSE_CMD="docker-compose"
fi

echo "Stopping backend, frontend, and nginx..."

$COMPOSE_CMD -f "$COMPOSE_FILE" stop backend frontend nginx

echo "Containers stopped. DB and pgadmin still running."
$COMPOSE_CMD -f "$COMPOSE_FILE" ps -a

