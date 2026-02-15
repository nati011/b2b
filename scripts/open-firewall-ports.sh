#!/usr/bin/env bash
# Open ports 80 and 443 for Docker nginx (B2B API).
# Run on the server: sudo ./scripts/open-firewall-ports.sh

set -euo pipefail

if [[ $EUID -ne 0 ]]; then
  echo "Run as root: sudo $0"
  exit 1
fi

if command -v ufw &>/dev/null; then
  echo "Opening UFW ports 80 and 443..."
  ufw allow 80/tcp
  ufw allow 443/tcp
  ufw status | grep -E '80|443' || true
  echo "Done. Reload firewall with: ufw reload"
else
  echo "UFW not found. If using firewalld:"
  echo "  firewall-cmd --permanent --add-service=http --add-service=https"
  echo "  firewall-cmd --reload"
  echo "If using iptables, allow: -A INPUT -p tcp --dport 80 -j ACCEPT (and 443)."
  echo "On cloud (AWS/GCP/Azure), open ports 80 and 443 in the security group / firewall rules."
fi
