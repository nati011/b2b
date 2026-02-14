#!/usr/bin/env bash
# Nginx setup script for Ubuntu
# Usage: sudo ./setup-nginx-ubuntu.sh [--reverse-proxy PORT] [--ssl DOMAIN]

set -euo pipefail

REVERSE_PROXY_PORT=""
SSL_DOMAIN=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --reverse-proxy)
      REVERSE_PROXY_PORT="$2"
      shift 2
      ;;
    --ssl)
      SSL_DOMAIN="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--reverse-proxy PORT] [--ssl DOMAIN]"
      exit 1
      ;;
  esac
done

if [[ $EUID -ne 0 ]]; then
  echo "Run this script as root (e.g. sudo $0)"
  exit 1
fi

echo "[1/6] Updating package index..."
apt-get update -qq

echo "[2/6] Installing Nginx..."
apt-get install -y -qq nginx

echo "[3/6] Enabling and starting Nginx..."
systemctl enable nginx
systemctl start nginx
systemctl status nginx --no-pager || true

echo "[4/6] Configuring Nginx..."

if [[ -n "$REVERSE_PROXY_PORT" ]]; then
  cat > /etc/nginx/sites-available/default << 'NGINX_UPSTREAM'
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:REVERSE_PROXY_PORT_PLACEHOLDER;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
NGINX_UPSTREAM
  sed -i "s/REVERSE_PROXY_PORT_PLACEHOLDER/$REVERSE_PROXY_PORT/" /etc/nginx/sites-available/default
  echo "  Reverse proxy to port $REVERSE_PROXY_PORT configured."
else
  cat > /etc/nginx/sites-available/default << 'NGINX_DEFAULT'
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    root /var/www/html;
    index index.html index.htm;
    server_name _;
    location / {
        try_files $uri $uri/ =404;
    }
}
NGINX_DEFAULT
  echo "  Default static site configured (root /var/www/html)."
fi

echo "[5/6] Testing Nginx config and reloading..."
nginx -t && systemctl reload nginx

echo "[6/6] Configuring firewall (UFW)..."
if command -v ufw &>/dev/null; then
  if ufw status | grep -q "Status: active"; then
    ufw allow 'Nginx HTTP' 2>/dev/null || ufw allow 80/tcp
    ufw --force reload 2>/dev/null || true
    echo "  UFW: Nginx HTTP (80) allowed."
  else
    echo "  UFW is installed but inactive. Enable with: ufw enable"
  fi
else
  echo "  UFW not installed. Install with: apt-get install ufw"
fi

if [[ -n "$SSL_DOMAIN" ]]; then
  echo ""
  echo "Installing Certbot for SSL..."
  apt-get install -y -qq certbot python3-certbot-nginx 2>/dev/null || apt-get install -y -qq certbot
  echo "Requesting certificate for $SSL_DOMAIN..."
  certbot --nginx -d "$SSL_DOMAIN" --non-interactive --agree-tos --register-unsafely-without-email || true
fi

echo ""
echo "Nginx setup complete."
echo "  Service: systemctl status nginx"
echo "  Logs:    journalctl -u nginx -f"
echo "  Config:  /etc/nginx/sites-available/default"
