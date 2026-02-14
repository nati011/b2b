#!/usr/bin/env bash
# Start Docker Compose stack with Nginx as the single entry point.
# All services (backend, pgAdmin) are reached via Nginx on port 80.
#
# Usage: ./scripts/setup-compose-nginx.sh [up|down|restart]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
COMPOSE_FILE="$PROJECT_ROOT/docker-compose.yml"
NGINX_CONF="$PROJECT_ROOT/nginx/default.conf"

cd "$PROJECT_ROOT"

if ! command -v docker &>/dev/null; then
  echo "Docker is not installed or not in PATH."
  exit 1
fi

if [ ! -f "$COMPOSE_FILE" ]; then
  echo "Compose file not found: $COMPOSE_FILE"
  exit 1
fi

if [ ! -f "$NGINX_CONF" ]; then
  echo "Nginx config not found: $NGINX_CONF"
  exit 1
fi

case "${1:-up}" in
  up)
    echo "Starting stack (db -> migrations -> seed -> backend -> nginx, pgadmin)..."
    docker compose -f "$COMPOSE_FILE" up -d
    echo ""
    echo "Stack is up. All traffic goes through Nginx on port 80:"
    echo "  Backend API:  http://localhost/"
    echo "  pgAdmin:      http://localhost/pgadmin/"
    echo ""
    echo "pgAdmin login: admin@local.local / admin"
    ;;
  down)
    echo "Stopping stack..."
    docker compose -f "$COMPOSE_FILE" down
    ;;
  restart)
    docker compose -f "$COMPOSE_FILE" restart nginx
    echo "Nginx restarted."
    ;;
  *)
    echo "Usage: $0 [up|down|restart]"
    exit 1
    ;;
esac
