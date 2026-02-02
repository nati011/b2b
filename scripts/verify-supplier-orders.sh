#!/bin/bash

# Script to verify supplier-based order fetch works correctly
# This script checks database state and tests the query logic

set +e

SUPPLIER_EMAIL="${1:-supplier1@b2b.local}"
PASSWORD="${2:-password123}"
API_BASE_URL="${3:-http://localhost:8090}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Supplier Order Fetch Verification${NC}"
echo -e "${BLUE}========================================${NC}"
echo "Supplier Email: $SUPPLIER_EMAIL"
echo "API Base URL: $API_BASE_URL"
echo ""

# Step 1: Get supplier info
echo -e "${BLUE}Step 1: Get Supplier Information${NC}"
SUPPLIER_RESPONSE=$(curl -s -u "${SUPPLIER_EMAIL}:${PASSWORD}" "${API_BASE_URL}/supplier")
SUPPLIER_ID=$(echo "$SUPPLIER_RESPONSE" | jq -r '.items[0].id // .id // empty')

if [ -z "$SUPPLIER_ID" ] || [ "$SUPPLIER_ID" = "null" ]; then
    echo -e "${RED}✗ Could not get supplier ID${NC}"
    echo "Response: $SUPPLIER_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✓ Supplier ID: $SUPPLIER_ID${NC}"
echo "Supplier Info:"
echo "$SUPPLIER_RESPONSE" | jq '.items[0] // .'
echo ""

# Step 2: Check if products exist for this supplier
echo -e "${BLUE}Step 2: Check Products for Supplier${NC}"
PRODUCTS_RESPONSE=$(curl -s -u "${SUPPLIER_EMAIL}:${PASSWORD}" "${API_BASE_URL}/product?limit=100")
PRODUCT_COUNT=$(echo "$PRODUCTS_RESPONSE" | jq '.products | length // 0')

if [ "$PRODUCT_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓ Found $PRODUCT_COUNT products for this supplier${NC}"
    echo "Sample products:"
    echo "$PRODUCTS_RESPONSE" | jq '.products[0:3] | .[] | {id, name, supplier_id}'
else
    echo -e "${YELLOW}⚠ No products found for this supplier${NC}"
fi
echo ""

# Step 3: Check all orders (to see what exists)
echo -e "${BLUE}Step 3: Check All Orders (Customer Endpoint)${NC}"
ALL_ORDERS_RESPONSE=$(curl -s -u "${SUPPLIER_EMAIL}:${PASSWORD}" "${API_BASE_URL}/orders/customer?limit=100")
ALL_ORDER_COUNT=$(echo "$ALL_ORDERS_RESPONSE" | jq '.orders | length // 0')

if [ "$ALL_ORDER_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓ Found $ALL_ORDER_COUNT total orders${NC}"
    echo "Sample order cart_snapshot:"
    FIRST_ORDER=$(echo "$ALL_ORDERS_RESPONSE" | jq '.orders[0]')
    echo "$FIRST_ORDER" | jq '{id, customer_id, status, cart_snapshot}'
else
    echo -e "${YELLOW}⚠ No orders found in system${NC}"
fi
echo ""

# Step 4: Test supplier orders endpoint
echo -e "${BLUE}Step 4: Test Supplier Orders Endpoint${NC}"
SUPPLIER_ORDERS_RESPONSE=$(curl -s -u "${SUPPLIER_EMAIL}:${PASSWORD}" "${API_BASE_URL}/orders/supplier")
SUPPLIER_ORDER_COUNT=$(echo "$SUPPLIER_ORDERS_RESPONSE" | jq '.orders | length // 0')
TOTAL=$(echo "$SUPPLIER_ORDERS_RESPONSE" | jq '.total // 0')

echo -e "${GREEN}✓ Supplier orders endpoint returned: $SUPPLIER_ORDER_COUNT orders (total: $TOTAL)${NC}"
echo "Response:"
echo "$SUPPLIER_ORDERS_RESPONSE" | jq '.'
echo ""

# Step 5: Analyze why orders might not be returned
if [ "$SUPPLIER_ORDER_COUNT" -eq 0 ] && [ "$ALL_ORDER_COUNT" -gt 0 ]; then
    echo -e "${YELLOW}Step 5: Analyzing Why Orders Are Not Returned${NC}"
    echo ""
    echo "Possible reasons:"
    echo "1. Products in orders don't have supplier_id = $SUPPLIER_ID"
    echo "2. cart_snapshot format doesn't match expected structure"
    echo "3. Products are deleted (is_deleted = TRUE)"
    echo ""
    
    # Check first order's cart_snapshot
    if [ "$ALL_ORDER_COUNT" -gt 0 ]; then
        FIRST_CART=$(echo "$ALL_ORDERS_RESPONSE" | jq -r '.orders[0].cart_snapshot // "[]"')
        echo "First order cart_snapshot structure:"
        echo "$FIRST_CART" | jq '.'
        echo ""
        
        # Extract product_ids from cart_snapshot
        PRODUCT_IDS=$(echo "$FIRST_CART" | jq -r '.[] | .product_id // empty' | head -5)
        if [ -n "$PRODUCT_IDS" ]; then
            echo "Product IDs in first order: $PRODUCT_IDS"
            echo ""
            echo "Checking if these products belong to supplier $SUPPLIER_ID:"
            for PID in $PRODUCT_IDS; do
                if [ -n "$PID" ] && [ "$PID" != "null" ]; then
                    PRODUCT_INFO=$(curl -s -u "${SUPPLIER_EMAIL}:${PASSWORD}" "${API_BASE_URL}/product?id=${PID}")
                    PROD_SUPPLIER_ID=$(echo "$PRODUCT_INFO" | jq -r '.supplier_id // empty')
                    if [ "$PROD_SUPPLIER_ID" = "$SUPPLIER_ID" ]; then
                        echo -e "  ${GREEN}✓ Product $PID belongs to supplier $SUPPLIER_ID${NC}"
                    else
                        echo -e "  ${RED}✗ Product $PID belongs to supplier $PROD_SUPPLIER_ID (not $SUPPLIER_ID)${NC}"
                    fi
                fi
            done
        fi
    fi
fi

# Step 6: Test with different suppliers
echo ""
echo -e "${BLUE}Step 6: Test with Different Suppliers${NC}"
for SUPPLIER in "supplier1@b2b.local" "supplier2@b2b.local" "supplier3@b2b.local"; do
    if [ "$SUPPLIER" != "$SUPPLIER_EMAIL" ]; then
        OTHER_ORDERS=$(curl -s -u "${SUPPLIER}:${PASSWORD}" "${API_BASE_URL}/orders/supplier")
        OTHER_COUNT=$(echo "$OTHER_ORDERS" | jq '.orders | length // 0')
        echo "$SUPPLIER: $OTHER_COUNT orders"
    fi
done

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Summary${NC}"
echo -e "${BLUE}========================================${NC}"
echo "Supplier ID: $SUPPLIER_ID"
echo "Products for supplier: $PRODUCT_COUNT"
echo "Total orders in system: $ALL_ORDER_COUNT"
echo "Orders for this supplier: $SUPPLIER_ORDER_COUNT"
echo ""

if [ "$SUPPLIER_ORDER_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓ Supplier-based order fetch is working!${NC}"
    exit 0
else
    if [ "$ALL_ORDER_COUNT" -eq 0 ]; then
        echo -e "${YELLOW}⚠ No orders exist in the system${NC}"
    else
        echo -e "${YELLOW}⚠ No orders found for this supplier${NC}"
        echo "This could mean:"
        echo "  - Orders don't contain products from this supplier"
        echo "  - Products in orders have different supplier_id"
        echo "  - cart_snapshot format issue"
    fi
    exit 1
fi

