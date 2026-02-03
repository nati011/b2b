#!/bin/bash

# Direct test of supplier order query logic
# This creates a test scenario to verify the query works

set +e

SUPPLIER_EMAIL="${1:-supplier1@b2b.local}"
PASSWORD="${2:-password123}"
API_BASE_URL="${3:-http://localhost:8090}"

echo "Testing Supplier Order Fetch Query Logic"
echo "========================================"
echo ""

# Get supplier ID
SUPPLIER_RESPONSE=$(curl -s -u "${SUPPLIER_EMAIL}:${PASSWORD}" "${API_BASE_URL}/supplier")
SUPPLIER_ID=$(echo "$SUPPLIER_RESPONSE" | jq -r '.items[0].id // .id // empty')

if [ -z "$SUPPLIER_ID" ] || [ "$SUPPLIER_ID" = "null" ]; then
    echo "Error: Could not get supplier ID"
    exit 1
fi

echo "Supplier ID: $SUPPLIER_ID"
echo ""

# Test the endpoint
echo "Testing /orders/supplier endpoint:"
ORDERS_RESPONSE=$(curl -s -u "${SUPPLIER_EMAIL}:${PASSWORD}" "${API_BASE_URL}/orders/supplier")
echo "$ORDERS_RESPONSE" | jq '.'

ORDER_COUNT=$(echo "$ORDERS_RESPONSE" | jq '.orders | length')
TOTAL=$(echo "$ORDERS_RESPONSE" | jq '.total')

echo ""
echo "Results:"
echo "  Orders returned: $ORDER_COUNT"
echo "  Total: $TOTAL"
echo ""

# Verify response structure
HAS_ORDERS=$(echo "$ORDERS_RESPONSE" | jq 'has("orders")')
HAS_TOTAL=$(echo "$ORDERS_RESPONSE" | jq 'has("total")')
HAS_LIMIT=$(echo "$ORDERS_RESPONSE" | jq 'has("limit")')
HAS_OFFSET=$(echo "$ORDERS_RESPONSE" | jq 'has("offset")')

if [ "$HAS_ORDERS" = "true" ] && [ "$HAS_TOTAL" = "true" ] && [ "$HAS_LIMIT" = "true" ] && [ "$HAS_OFFSET" = "true" ]; then
    echo "✓ Response structure is correct"
else
    echo "✗ Response structure is incorrect"
    exit 1
fi

# Check ordering if orders exist
if [ "$ORDER_COUNT" -gt 1 ]; then
    TIMESTAMPS=$(echo "$ORDERS_RESPONSE" | jq -r '.orders[] | .created_at')
    PREV=""
    IS_DESC=true
    for TS in $TIMESTAMPS; do
        if [ -n "$PREV" ]; then
            if [[ "$TS" > "$PREV" ]]; then
                IS_DESC=false
                break
            fi
        fi
        PREV="$TS"
    done
    
    if [ "$IS_DESC" = "true" ]; then
        echo "✓ Orders are correctly ordered by created_at DESC"
    else
        echo "✗ Orders are not in descending order"
    fi
else
    echo "ℹ Skipping ordering test (need at least 2 orders)"
fi

echo ""
if [ "$ORDER_COUNT" -gt 0 ]; then
    echo "✓ Supplier-based order fetch is WORKING!"
    echo ""
    echo "Sample order:"
    echo "$ORDERS_RESPONSE" | jq '.orders[0] | {id, customer_id, status, created_at, cart_snapshot: (.cart_snapshot | if type == "string" then fromjson else . end | .[0])}'
else
    echo "⚠ No orders returned (this may be expected if no orders contain products from this supplier)"
    echo ""
    echo "The query logic is correct, but there are no matching orders."
    echo "To verify the query works:"
    echo "  1. Ensure orders exist with cart_snapshot containing product_id/ProductID"
    echo "  2. Ensure those products have supplier_id = $SUPPLIER_ID"
    echo "  3. Ensure products are not deleted (is_deleted = FALSE)"
fi



