#!/bin/sh
set -e
CERT_DIR=/etc/nginx/certs
CERT_FILE=$CERT_DIR/server.crt
KEY_FILE=$CERT_DIR/server.key

if [ ! -f "$CERT_FILE" ] || [ ! -f "$KEY_FILE" ]; then
  mkdir -p "$CERT_DIR"
  openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$KEY_FILE" -out "$CERT_FILE" \
    -subj "/CN=localhost/O=B2B/C=ET"
  echo "Generated self-signed certificate at $CERT_FILE"
fi

exec nginx -g "daemon off;"
