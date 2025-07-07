#!/bin/bash

# Database connection string for Neon
DB_CONNECTION_STRING="${DB_CONNECTION_STRING:-postgresql://neondb_owner:npg_vI3ByRrKu9PX@ep-nameless-sun-a55nmt8x-pooler.us-east-2.aws.neon.tech/neondb?sslmode=require}"

# Function to execute SQL and create resource
create_resource() {
    local name="$1"
    local action="$2"
    local resource="$3"
    
    echo "Creating resource: $name"
    
    # Execute the PostgreSQL function using connection string
    psql "$DB_CONNECTION_STRING" -t -c \
        "SELECT public.create_resource('$name', '$action', '$resource');"
    
    if [ $? -eq 0 ]; then
        echo "✓ Successfully created resource: $name"
    else
        echo "✗ Failed to create resource: $name"
        return 1
    fi
}

# Main script
echo "Starting resource creation..."
echo "================================"

# Create resources based on your Go code
# Auth resources
create_resource "auth_login" "ALL" "/api/v1/auth/login"
create_resource "auth_logout" "ALL" "/api/v1/auth/logout"
create_resource "auth_refresh" "ALL" "/api/v1/auth/refresh"
create_resource "auth_reset" "ALL" "/api/v1/auth/reset/{param}"

# Email resources
create_resource "email_template" "ALL" "/api/v1/email_template"
create_resource "email" "ALL" "/api/v1/email"

# Health resources
create_resource "health" "ALL" "/health"
create_resource "health_metrics" "ALL" "/metrics/promethus"

# Identity resources
create_resource "identity_user" "ALL" "/api/v1/identity/user"

# Payment option resources
create_resource "payment_option" "ALL" "/api/v1/payment_option"
create_resource "activate_payment_option" "ALL" "/api/v1/payment_option/activate"
create_resource "payment_option_status" "ALL" "/api/v1/payment_option/{param}/status"
create_resource "payment_option_secret" "ALL" "/api/v1/payment_option/{param}/secret"

# Resource management
create_resource "resource" "ALL" "/api/v1/resource"

# Role resources
create_resource "role" "ALL" "/api/v1/role"
create_resource "role_resource" "ALL" "/api/v1/role/resource"
create_resource "role_list" "ALL" "/api/v1/role/{param}"

# Transaction resources
create_resource "transaction" "ALL" "/api/v1/transaction"

# User resources
create_resource "user" "ALL" "/api/v1/user"
create_resource "user_status" "ALL" "/api/v1/user/{param}/status"
create_resource "user_role" "ALL" "/api/v1/user/{param}/role/{param}"
create_resource "user_role_list" "ALL" "/api/v1/user/{param}/role"
create_resource "user_rest" "ALL" "/api/v1/user/init_reset"

# Catalogue resources
create_resource "catalogue" "ALL" "/api/v1/catalogue"

# Category resources
create_resource "category_list" "ALL" "/api/v1/category"
create_resource "category" "ALL" "/api/v1/category/{param}"

# Config resources
create_resource "config_get" "GET" "/api/v1/config"
create_resource "config_post" "POST" "/api/v1/config"

# Configurable product resources
create_resource "configurable_product" "ALL" "/api/v1/configurable_product"
create_resource "configurable_product_status" "ALL" "/api/v1/configurable_product/{param}/status"

# Distributor resources
create_resource "distributor" "ALL" "/api/v1/distributor"
create_resource "distributor_user" "ALL" "/api/v1/distributor/{param}/user"
create_resource "create_distributor_user" "POST" "/api/v1/distributor/user"
create_resource "distributor_status" "ALL" "/api/v1/distributor/{param}/status"
create_resource "distributor_onboarding_review" "ALL" "/api/v1/distributor/{param}/onboarding_review"

# Invoice resources
create_resource "invoice" "ALL" "/api/v1/invoice"
create_resource "invoice_document" "ALL" "/api/v1/invoice/html"

# Order resources
create_resource "order" "ALL" "/api/v1/order"
create_resource "order_init_payment" "ALL" "/api/v1/order/init_payment"
create_resource "orders_retailer" "ALL" "/api/v1/orders/retailer"
create_resource "orders_distributor" "ALL" "/api/v1/orders/distributor"

# Payment resources
create_resource "payment_webook_callback_var_1" "ALL" "/api/v1/payment/webhook/{param}/{param}"
create_resource "payment_webook_callback_var_2" "ALL" "/api/v1/payment/webhook/{param}"
create_resource "payment_verify" "ALL" "/api/v1/payment/verify"
create_resource "payment_confirm" "ALL" "/api/v1/payment/confirm/{param}"

# Product resources
create_resource "product" "ALL" "/api/v1/product"
create_resource "stock_ledger" "ALL" "/api/v1/stock_ledger"
create_resource "product_status" "ALL" "/api/v1/product/{param}/status"
create_resource "product_stock" "ALL" "/api/v1/product/{param}/stock"

# Retailer resources
create_resource "retailer" "ALL" "/api/v1/retailer"
create_resource "retailer_user" "ALL" "/api/v1/retailer/{param}/user"

echo "================================"
echo "Resource creation completed!"