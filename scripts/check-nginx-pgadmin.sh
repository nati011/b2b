#!/usr/bin/env bash
# Quick check that nginx is on port 80 and /pgadmin is routed to pgadmin container.
set -e
echo "=== Port 80 listener ==="
sudo ss -tlnp | grep ':80 ' || true
echo ""
echo "=== Docker containers (nginx, pgadmin) ==="
docker ps --format 'table {{.Names}}\t{{.Ports}}' | grep -E 'nginx|pgadmin|NAMES' || true
echo ""
echo "=== Curl http://localhost/pgadmin/ (expect 200 or 302 from pgAdmin, not 404) ==="
curl -sI http://127.0.0.1/pgadmin/ 2>/dev/null | head -5
