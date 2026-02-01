#!/bin/bash

# Script to create authentication credentials for a user
# Usage: ./scripts/create-credentials.sh [email] [password] [user_id] [api_base_url]

set -e

EMAIL="${1:-customer1@b2b.local}"
PASSWORD="${2:-password123}"
USER_ID="${3:-550e8400-e29b-41d4-a716-446655440004}"
API_BASE_URL="${4:-http://localhost:8090}"

echo "Creating credentials for user: $EMAIL"
echo "User ID: $USER_ID"
echo "API Base URL: $API_BASE_URL"
echo ""

# Try to create credentials with user_id
echo "Attempting to create credentials..."
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE_URL}/auth/basic/credentials" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"${EMAIL}\",
    \"password\": \"${PASSWORD}\",
    \"user_id\": \"${USER_ID}\"
  }")

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -eq 201 ]; then
  echo "✓ Credentials created successfully!"
  echo "Response: $BODY"
  echo ""
  echo "You can now login with:"
  echo "  Email: $EMAIL"
  echo "  Password: $PASSWORD"
elif [ "$HTTP_CODE" -eq 400 ] || [ "$HTTP_CODE" -eq 403 ]; then
  echo "✗ Failed to create credentials (HTTP $HTTP_CODE)"
  echo "Response: $BODY"
  echo ""
  echo "Note: This endpoint may require authentication or a registration token."
  echo "If you have an admin token, you can use it like this:"
  echo ""
  echo "  curl -X POST \"${API_BASE_URL}/auth/basic/credentials\" \\"
  echo "    -H \"Content-Type: application/json\" \\"
  echo "    -H \"Authorization: Bearer YOUR_ADMIN_TOKEN\" \\"
  echo "    -d '{\"username\": \"${EMAIL}\", \"password\": \"${PASSWORD}\", \"user_id\": \"${USER_ID}\"}'"
else
  echo "✗ Unexpected response (HTTP $HTTP_CODE)"
  echo "Response: $BODY"
fi


