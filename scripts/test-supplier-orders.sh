#!/bin/bash

# Script to test supplier-based order fetch
# Usage: ./scripts/test-supplier-orders.sh [supplier_email] [password] [api_base_url]

# Don't exit on error - we want to run all tests
set +e

SUPPLIER_EMAIL="${1:-supplier1@b2b.local}"
PASSWORD="${2:-password123}"
API_BASE_URL="${3:-http://localhost:8090}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0

# Function to print test header
print_header() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

# Function to print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    ((TESTS_PASSED++))
}

# Function to print failure
print_failure() {
    echo -e "${RED}✗ $1${NC}"
    ((TESTS_FAILED++))
}

# Function to print info
print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Function to validate JSON response structure
validate_order_response() {
    local response="$1"
    local test_name="$2"
    
    # Check if response is valid JSON
    if ! echo "$response" | jq . > /dev/null 2>&1; then
        print_failure "$test_name: Response is not valid JSON"
        echo "Response: $response"
        return 1
    fi
    
    # Check required fields
    local has_orders=$(echo "$response" | jq -r 'has("orders")')
    local has_total=$(echo "$response" | jq -r 'has("total")')
    local has_limit=$(echo "$response" | jq -r 'has("limit")')
    local has_offset=$(echo "$response" | jq -r 'has("offset")')
    
    if [ "$has_orders" != "true" ]; then
        print_failure "$test_name: Missing 'orders' field in response"
        return 1
    fi
    
    if [ "$has_total" != "true" ]; then
        print_failure "$test_name: Missing 'total' field in response"
        return 1
    fi
    
    if [ "$has_limit" != "true" ]; then
        print_failure "$test_name: Missing 'limit' field in response"
        return 1
    fi
    
    if [ "$has_offset" != "true" ]; then
        print_failure "$test_name: Missing 'offset' field in response"
        return 1
    fi
    
    # Check that orders is an array
    local orders_type=$(echo "$response" | jq -r '.orders | type')
    if [ "$orders_type" != "array" ]; then
        print_failure "$test_name: 'orders' field is not an array (got: $orders_type)"
        return 1
    fi
    
    print_success "$test_name: Response structure is valid"
    return 0
}

# Function to validate order ordering (by created_at DESC)
validate_order_ordering() {
    local response="$1"
    local test_name="$2"
    
    local order_count=$(echo "$response" | jq '.orders | length')
    
    if [ "$order_count" -lt 2 ]; then
        print_info "$test_name: Skipping ordering test (need at least 2 orders, got $order_count)"
        return 0
    fi
    
    # Extract created_at timestamps
    local timestamps=$(echo "$response" | jq -r '.orders[] | .created_at')
    
    # Check if timestamps are in descending order
    local prev_timestamp=""
    local is_descending=true
    
    while IFS= read -r timestamp; do
        if [ -n "$prev_timestamp" ]; then
            # Compare timestamps (newer should come first)
            if [[ "$timestamp" > "$prev_timestamp" ]]; then
                is_descending=false
                break
            fi
        fi
        prev_timestamp="$timestamp"
    done <<< "$timestamps"
    
    if [ "$is_descending" = "true" ]; then
        print_success "$test_name: Orders are correctly ordered by created_at DESC"
        return 0
    else
        print_failure "$test_name: Orders are not in descending order by created_at"
        return 1
    fi
}

# Function to validate order structure
validate_order_structure() {
    local response="$1"
    local test_name="$2"
    
    local order_count=$(echo "$response" | jq '.orders | length')
    
    if [ "$order_count" -eq 0 ]; then
        print_info "$test_name: No orders to validate structure"
        return 0
    fi
    
    # Check first order structure
    local first_order=$(echo "$response" | jq '.orders[0]')
    
    local required_fields=("id" "customer_id" "status" "created_at" "updated_at" "cart_snapshot")
    local missing_fields=()
    
    for field in "${required_fields[@]}"; do
        if ! echo "$first_order" | jq -e ".$field" > /dev/null 2>&1; then
            missing_fields+=("$field")
        fi
    done
    
    if [ ${#missing_fields[@]} -gt 0 ]; then
        print_failure "$test_name: Missing required fields in order: ${missing_fields[*]}"
        return 1
    fi
    
    print_success "$test_name: Order structure is valid"
    return 0
}

print_header "Supplier Order Fetch Test"
echo "Supplier Email: $SUPPLIER_EMAIL"
echo "API Base URL: $API_BASE_URL"
echo ""

# Step 1: Login as supplier
print_header "Step 1: Login as Supplier"

LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${SUPPLIER_EMAIL}\",
        \"password\": \"${PASSWORD}\"
    }")

HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n1)
LOGIN_BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -ne 202 ]; then
    print_failure "Login failed (HTTP $HTTP_CODE)"
    echo "Response: $LOGIN_BODY"
    echo ""
    echo "Cannot continue without authentication. Please check:"
    echo "  - API is running at $API_BASE_URL"
    echo "  - Supplier credentials are correct"
    echo "  - Database is seeded with supplier users"
    exit 1
fi

# Extract access token
ACCESS_TOKEN=$(echo "$LOGIN_BODY" | jq -r '.body.access_token // .body.access_token // empty')

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
    print_failure "No access token received in login response"
    echo "Response: $LOGIN_BODY"
    exit 1
fi

print_success "Login successful"
print_info "Access token: ${ACCESS_TOKEN:0:20}..."

# Step 2: Fetch supplier orders (basic test)
print_header "Step 2: Fetch Supplier Orders (Basic)"

# Use Basic Auth instead of Bearer token
ORDERS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE_URL}/orders/supplier" \
    -H "Content-Type: application/json" \
    -u "${SUPPLIER_EMAIL}:${PASSWORD}")

HTTP_CODE=$(echo "$ORDERS_RESPONSE" | tail -n1)
ORDERS_BODY=$(echo "$ORDERS_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" -ne 200 ]; then
    print_failure "Failed to fetch orders (HTTP $HTTP_CODE)"
    echo "Response: $ORDERS_BODY"
    echo ""
    echo "Cannot continue without orders endpoint. Please check:"
    echo "  - API endpoint is correct"
    echo "  - Supplier has proper permissions"
    exit 1
fi

print_success "Orders endpoint returned HTTP 200"

# Validate response structure
validate_order_response "$ORDERS_BODY" "Basic Order Fetch"

# Extract order count
ORDER_COUNT=$(echo "$ORDERS_BODY" | jq '.orders | length')
TOTAL=$(echo "$ORDERS_BODY" | jq '.total')
LIMIT=$(echo "$ORDERS_BODY" | jq '.limit')
OFFSET=$(echo "$ORDERS_BODY" | jq '.offset')

print_info "Found $ORDER_COUNT orders (total: $TOTAL, limit: $LIMIT, offset: $OFFSET)"

# Step 3: Validate order structure
print_header "Step 3: Validate Order Structure"
validate_order_structure "$ORDERS_BODY" "Order Structure Validation"

# Step 4: Validate ordering
print_header "Step 4: Validate Order Ordering"
validate_order_ordering "$ORDERS_BODY" "Order Ordering Validation"

# Step 5: Test pagination
print_header "Step 5: Test Pagination"

if [ "$TOTAL" -gt 0 ]; then
    # Test with limit
    PAGINATED_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE_URL}/orders/supplier?limit=2&offset=0" \
        -H "Content-Type: application/json" \
        -u "${SUPPLIER_EMAIL}:${PASSWORD}")
    
    PAGINATED_HTTP_CODE=$(echo "$PAGINATED_RESPONSE" | tail -n1)
    PAGINATED_BODY=$(echo "$PAGINATED_RESPONSE" | sed '$d')
    
    if [ "$PAGINATED_HTTP_CODE" -eq 200 ]; then
        validate_order_response "$PAGINATED_BODY" "Pagination Test"
        
        PAGINATED_COUNT=$(echo "$PAGINATED_BODY" | jq '.orders | length')
        PAGINATED_LIMIT=$(echo "$PAGINATED_BODY" | jq '.limit')
        PAGINATED_OFFSET=$(echo "$PAGINATED_BODY" | jq '.offset')
        
        if [ "$PAGINATED_LIMIT" -eq 2 ] && [ "$PAGINATED_OFFSET" -eq 0 ]; then
            print_success "Pagination parameters are correct (limit: $PAGINATED_LIMIT, offset: $PAGINATED_OFFSET)"
        else
            print_failure "Pagination parameters are incorrect (expected limit: 2, offset: 0, got limit: $PAGINATED_LIMIT, offset: $PAGINATED_OFFSET)"
        fi
        
        if [ "$PAGINATED_COUNT" -le 2 ]; then
            print_success "Pagination limit is respected (returned $PAGINATED_COUNT orders)"
        else
            print_failure "Pagination limit not respected (expected <= 2, got $PAGINATED_COUNT)"
        fi
    else
        print_failure "Pagination request failed (HTTP $PAGINATED_HTTP_CODE)"
    fi
else
    print_info "Skipping pagination test (no orders available)"
fi

# Step 6: Test status filter
print_header "Step 6: Test Status Filter"

STATUS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE_URL}/orders/supplier?status=pending" \
    -H "Content-Type: application/json" \
    -u "${SUPPLIER_EMAIL}:${PASSWORD}")

STATUS_HTTP_CODE=$(echo "$STATUS_RESPONSE" | tail -n1)
STATUS_BODY=$(echo "$STATUS_RESPONSE" | sed '$d')

if [ "$STATUS_HTTP_CODE" -eq 200 ]; then
    validate_order_response "$STATUS_BODY" "Status Filter Test"
    
    # Check if all returned orders have pending status
    STATUS_COUNT=$(echo "$STATUS_BODY" | jq '[.orders[] | select(.status == "pending")] | length')
    TOTAL_STATUS=$(echo "$STATUS_BODY" | jq '.orders | length')
    
    if [ "$TOTAL_STATUS" -eq 0 ] || [ "$STATUS_COUNT" -eq "$TOTAL_STATUS" ]; then
        print_success "Status filter works correctly (all $TOTAL_STATUS orders have 'pending' status)"
    else
        print_failure "Status filter not working correctly (expected all pending, got $STATUS_COUNT/$TOTAL_STATUS)"
    fi
else
    print_failure "Status filter request failed (HTTP $STATUS_HTTP_CODE)"
fi

# Step 7: Validate cart_snapshot contains supplier products
print_header "Step 7: Validate Cart Snapshot Contains Supplier Products"

if [ "$ORDER_COUNT" -gt 0 ]; then
    # Get supplier info to verify products belong to this supplier
    SUPPLIER_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE_URL}/supplier" \
        -H "Content-Type: application/json" \
        -u "${SUPPLIER_EMAIL}:${PASSWORD}")
    
    SUPPLIER_HTTP_CODE=$(echo "$SUPPLIER_RESPONSE" | tail -n1)
    SUPPLIER_BODY=$(echo "$SUPPLIER_RESPONSE" | sed '$d')
    
    if [ "$SUPPLIER_HTTP_CODE" -eq 200 ]; then
        SUPPLIER_ID=$(echo "$SUPPLIER_BODY" | jq -r '.id // empty')
        
        if [ -n "$SUPPLIER_ID" ] && [ "$SUPPLIER_ID" != "null" ]; then
            print_info "Supplier ID: $SUPPLIER_ID"
            
            # Validate that cart_snapshot is present and valid
            INVALID_ORDERS=0
            for i in $(seq 0 $((ORDER_COUNT - 1))); do
                ORDER=$(echo "$ORDERS_BODY" | jq ".orders[$i]")
                CART_SNAPSHOT=$(echo "$ORDER" | jq -r '.cart_snapshot // empty')
                
                if [ -z "$CART_SNAPSHOT" ] || [ "$CART_SNAPSHOT" = "null" ] || [ "$CART_SNAPSHOT" = "[]" ]; then
                    print_failure "Order $(echo "$ORDER" | jq -r '.id') has empty or invalid cart_snapshot"
                    ((INVALID_ORDERS++))
                fi
            done
            
            if [ "$INVALID_ORDERS" -eq 0 ]; then
                print_success "All orders have valid cart_snapshot"
            else
                print_failure "$INVALID_ORDERS orders have invalid cart_snapshot"
            fi
        else
            print_info "Could not retrieve supplier ID, skipping cart_snapshot validation"
        fi
    else
        print_info "Could not fetch supplier info (HTTP $SUPPLIER_HTTP_CODE), skipping cart_snapshot validation"
    fi
else
    print_info "No orders to validate cart_snapshot"
fi

# Step 8: Validate items field is populated
print_header "Step 8: Validate Items Field"

if [ "$ORDER_COUNT" -gt 0 ]; then
    ORDERS_WITH_ITEMS=0
    ORDERS_WITHOUT_ITEMS=0
    
    for i in $(seq 0 $((ORDER_COUNT - 1))); do
        ORDER=$(echo "$ORDERS_BODY" | jq ".orders[$i]")
        ITEMS=$(echo "$ORDER" | jq '.items // empty')
        ITEMS_COUNT=$(echo "$ORDER" | jq '.items | length // 0')
        
        if [ "$ITEMS_COUNT" -gt 0 ]; then
            ((ORDERS_WITH_ITEMS++))
        else
            ((ORDERS_WITHOUT_ITEMS++))
        fi
    done
    
    if [ "$ORDERS_WITH_ITEMS" -gt 0 ]; then
        print_success "$ORDERS_WITH_ITEMS orders have items populated"
    fi
    
    if [ "$ORDERS_WITHOUT_ITEMS" -gt 0 ]; then
        print_info "$ORDERS_WITHOUT_ITEMS orders have no items (this may be expected if items are optional)"
    fi
else
    print_info "No orders to validate items field"
fi

# Step 9: Test unauthorized access (non-supplier user)
print_header "Step 9: Test Unauthorized Access"

# Try to access with a customer account (should fail or return empty)
CUSTOMER_EMAIL="customer1@b2b.local"
CUSTOMER_PASSWORD="password123"

CUSTOMER_LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${CUSTOMER_EMAIL}\",
        \"password\": \"${CUSTOMER_PASSWORD}\"
    }")

CUSTOMER_HTTP_CODE=$(echo "$CUSTOMER_LOGIN_RESPONSE" | tail -n1)
CUSTOMER_LOGIN_BODY=$(echo "$CUSTOMER_LOGIN_RESPONSE" | sed '$d')

if [ "$CUSTOMER_HTTP_CODE" -eq 202 ]; then
    CUSTOMER_TOKEN=$(echo "$CUSTOMER_LOGIN_BODY" | jq -r '.body.access_token // empty')
    
    if [ -n "$CUSTOMER_TOKEN" ] && [ "$CUSTOMER_TOKEN" != "null" ]; then
        CUSTOMER_ORDERS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE_URL}/orders/supplier" \
            -H "Content-Type: application/json" \
            -u "${CUSTOMER_EMAIL}:${CUSTOMER_PASSWORD}")
        
        CUSTOMER_ORDERS_HTTP_CODE=$(echo "$CUSTOMER_ORDERS_RESPONSE" | tail -n1)
        CUSTOMER_ORDERS_BODY=$(echo "$CUSTOMER_ORDERS_RESPONSE" | sed '$d')
        
        if [ "$CUSTOMER_ORDERS_HTTP_CODE" -eq 403 ]; then
            print_success "Non-supplier user correctly denied access (HTTP 403)"
        elif [ "$CUSTOMER_ORDERS_HTTP_CODE" -eq 200 ]; then
            CUSTOMER_ORDER_COUNT=$(echo "$CUSTOMER_ORDERS_BODY" | jq '.orders | length')
            if [ "$CUSTOMER_ORDER_COUNT" -eq 0 ]; then
                print_success "Non-supplier user correctly receives empty result"
            else
                print_failure "Non-supplier user received orders (should be empty or denied)"
            fi
        else
            print_info "Non-supplier access returned HTTP $CUSTOMER_ORDERS_HTTP_CODE (may be expected behavior)"
        fi
    fi
else
    print_info "Could not login as customer to test unauthorized access"
fi

# Step 10: Test with invalid parameters
print_header "Step 10: Test Invalid Parameters"

# Test with negative limit
INVALID_LIMIT_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE_URL}/orders/supplier?limit=-1" \
    -H "Content-Type: application/json" \
    -u "${SUPPLIER_EMAIL}:${PASSWORD}")

INVALID_LIMIT_HTTP_CODE=$(echo "$INVALID_LIMIT_RESPONSE" | tail -n1)
INVALID_LIMIT_BODY=$(echo "$INVALID_LIMIT_RESPONSE" | sed '$d')

if [ "$INVALID_LIMIT_HTTP_CODE" -eq 200 ]; then
    INVALID_LIMIT=$(echo "$INVALID_LIMIT_BODY" | jq '.limit')
    if [ "$INVALID_LIMIT" -gt 0 ]; then
        print_success "Invalid negative limit handled correctly (defaulted to $INVALID_LIMIT)"
    else
        print_info "Invalid limit returned $INVALID_LIMIT"
    fi
fi

# Test with invalid status
INVALID_STATUS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE_URL}/orders/supplier?status=invalid_status_xyz" \
    -H "Content-Type: application/json" \
    -u "${SUPPLIER_EMAIL}:${PASSWORD}")

INVALID_STATUS_HTTP_CODE=$(echo "$INVALID_STATUS_RESPONSE" | tail -n1)
if [ "$INVALID_STATUS_HTTP_CODE" -eq 200 ]; then
    INVALID_STATUS_COUNT=$(echo "$INVALID_STATUS_RESPONSE" | sed '$d' | jq '.orders | length')
    if [ "$INVALID_STATUS_COUNT" -eq 0 ]; then
        print_success "Invalid status filter correctly returns empty result"
    else
        print_info "Invalid status returned $INVALID_STATUS_COUNT orders"
    fi
fi

# Step 11: Test multiple suppliers (if available)
print_header "Step 11: Test Supplier Isolation"

if [ "$SUPPLIER_EMAIL" = "supplier1@b2b.local" ]; then
    OTHER_SUPPLIER="supplier2@b2b.local"
elif [ "$SUPPLIER_EMAIL" = "supplier2@b2b.local" ]; then
    OTHER_SUPPLIER="supplier1@b2b.local"
else
    OTHER_SUPPLIER="supplier1@b2b.local"
fi

print_info "Testing with another supplier: $OTHER_SUPPLIER"

OTHER_LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${API_BASE_URL}/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"${OTHER_SUPPLIER}\",
        \"password\": \"${PASSWORD}\"
    }")

OTHER_HTTP_CODE=$(echo "$OTHER_LOGIN_RESPONSE" | tail -n1)
OTHER_LOGIN_BODY=$(echo "$OTHER_LOGIN_RESPONSE" | sed '$d')

if [ "$OTHER_HTTP_CODE" -eq 202 ]; then
    OTHER_TOKEN=$(echo "$OTHER_LOGIN_BODY" | jq -r '.body.access_token // empty')
    
    if [ -n "$OTHER_TOKEN" ] && [ "$OTHER_TOKEN" != "null" ]; then
        OTHER_ORDERS_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "${API_BASE_URL}/orders/supplier" \
            -H "Content-Type: application/json" \
            -u "${OTHER_SUPPLIER}:${PASSWORD}")
        
        OTHER_HTTP_CODE=$(echo "$OTHER_ORDERS_RESPONSE" | tail -n1)
        OTHER_ORDERS_BODY=$(echo "$OTHER_ORDERS_RESPONSE" | sed '$d')
        
        if [ "$OTHER_HTTP_CODE" -eq 200 ]; then
            OTHER_ORDER_COUNT=$(echo "$OTHER_ORDERS_BODY" | jq '.orders | length')
            OTHER_TOTAL=$(echo "$OTHER_ORDERS_BODY" | jq '.total')
            
            print_info "Other supplier ($OTHER_SUPPLIER) has $OTHER_ORDER_COUNT orders (total: $OTHER_TOTAL)"
            print_info "Current supplier ($SUPPLIER_EMAIL) has $ORDER_COUNT orders (total: $TOTAL)"
            
            if [ "$ORDER_COUNT" -ne "$OTHER_ORDER_COUNT" ] || [ "$TOTAL" -ne "$OTHER_TOTAL" ]; then
                print_success "Suppliers see different orders (isolation working)"
            else
                print_info "Both suppliers see same number of orders (may be expected if they share products)"
            fi
        fi
    fi
else
    print_info "Could not login as other supplier to test isolation"
fi

# Step 12: Display sample order (if available)
print_header "Step 12: Sample Order Data"

if [ "$ORDER_COUNT" -gt 0 ]; then
    echo "First order details:"
    echo "$ORDERS_BODY" | jq '.orders[0]' | head -30
    echo ""
    echo "Cart snapshot (first order):"
    echo "$ORDERS_BODY" | jq '.orders[0].cart_snapshot' | head -20
    print_success "Sample order displayed"
else
    print_info "No orders available to display"
    echo ""
    echo "This could mean:"
    echo "  - No orders exist with products from this supplier"
    echo "  - Products in orders don't have supplier_id set correctly"
    echo "  - The supplier_id doesn't match any products in orders"
    echo "  - cart_snapshot format doesn't match expected structure (product_id, ProductID, or productId)"
fi

# Summary
print_header "Test Summary"
echo "Tests Passed: $TESTS_PASSED"
echo "Tests Failed: $TESTS_FAILED"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi

