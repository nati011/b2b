#!/bin/bash

# Script to create authentication credentials for a user (requires admin authentication)
# Usage: ./scripts/create-credentials-with-auth.sh [email] [password] [user_id] [api_base_url] [admin_token]

set -e

EMAIL="${1:-customer1@b2b.local}"
PASSWORD="${2:-password123}"
USER_ID="${3:-550e8400-e29b-41d4-a716-446655440004}"
API_BASE_URL="${4:-http://localhost:8090}"
ADMIN_TOKEN="${5:-}"

echo "Creating credentials for user: $EMAIL"
echo "User ID: $USER_ID"
echo "API Base URL: $API_BASE_URL"
echo ""

if [ -z "$ADMIN_TOKEN" ]; then
  echo "⚠️  No admin token provided. Attempting unauthenticated request..."
  echo "   (This will likely fail if user is not self-service type)"
  echo ""
  
  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE_URL}/auth/basic/credentials" \
    -H "Content-Type: application/json" \
    -d "{
      \"username\": \"${EMAIL}\",
      \"password\": \"${PASSWORD}\",
      \"user_id\": \"${USER_ID}\"
    }")
else
  echo "Using admin token for authenticated request..."
  echo ""
  
  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE_URL}/auth/basic/credentials" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${ADMIN_TOKEN}" \
    -d "{
      \"username\": \"${EMAIL}\",
      \"password\": \"${PASSWORD}\",
      \"user_id\": \"${USER_ID}\"
    }")
fi

HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -eq 201 ]; then
  echo "✓ Credentials created successfully!"
  echo "Response: $BODY"
  echo ""
  echo "You can now login with:"
  echo "  Email: $EMAIL"
  echo "  Password: $PASSWORD"
  exit 0
elif [ "$HTTP_CODE" -eq 400 ]; then
  echo "✗ Bad Request (HTTP $HTTP_CODE)"
  echo "Response: $BODY"
  echo ""
  echo "Possible issues:"
  echo "  - Missing required fields"
  echo "  - Invalid user_id"
  echo "  - Username already exists"
  exit 1
elif [ "$HTTP_CODE" -eq 401 ]; then
  echo "✗ Unauthorized (HTTP $HTTP_CODE)"
  echo "Response: $BODY"
  echo ""
  echo "You need to provide a valid admin token:"
  echo "  ./scripts/create-credentials-with-auth.sh \"$EMAIL\" \"$PASSWORD\" \"$USER_ID\" \"$API_BASE_URL\" \"YOUR_ADMIN_TOKEN\""
  exit 1
elif [ "$HTTP_CODE" -eq 403 ]; then
  echo "✗ Forbidden (HTTP $HTTP_CODE)"
  echo "Response: $BODY"
  echo ""
  echo "The authenticated user doesn't have permission to create credentials."
  exit 1
else
  echo "✗ Unexpected response (HTTP $HTTP_CODE)"
  echo "Response: $BODY"
  exit 1
fi


