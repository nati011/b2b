-- ============================================
-- Clear All Database Data
-- ============================================
-- This script truncates all tables to remove all data
-- while preserving the schema structure
-- ============================================

BEGIN;

-- Disable foreign key checks temporarily by truncating in the right order
-- First, truncate tables that have foreign keys pointing to them last

-- Truncate tables with foreign key dependencies first (child tables)
TRUNCATE TABLE IF EXISTS public.invoice_line_items CASCADE;
TRUNCATE TABLE IF EXISTS public.o_items CASCADE;
TRUNCATE TABLE IF EXISTS public.invoices CASCADE;
TRUNCATE TABLE IF EXISTS public.orders CASCADE;
TRUNCATE TABLE IF EXISTS public.s_ledger CASCADE;
TRUNCATE TABLE IF EXISTS public.p_stock CASCADE;
TRUNCATE TABLE IF EXISTS public.p_category CASCADE;
TRUNCATE TABLE IF EXISTS public.p_attribute_values CASCADE;
TRUNCATE TABLE IF EXISTS public.cp_attributes CASCADE;
TRUNCATE TABLE IF EXISTS public.cp_products CASCADE;
TRUNCATE TABLE IF EXISTS public.cp_images CASCADE;
TRUNCATE TABLE IF EXISTS public.configurable_products CASCADE;
TRUNCATE TABLE IF EXISTS public.p_images CASCADE;
TRUNCATE TABLE IF EXISTS public.products CASCADE;
TRUNCATE TABLE IF EXISTS public.p_attributes CASCADE;
TRUNCATE TABLE IF EXISTS public.category CASCADE;
TRUNCATE TABLE IF EXISTS public.distributor_subscriptions CASCADE;
TRUNCATE TABLE IF EXISTS public.distributor_business_info CASCADE;
TRUNCATE TABLE IF EXISTS public.distributor_reviews CASCADE;
TRUNCATE TABLE IF EXISTS public.distributor_users CASCADE;
TRUNCATE TABLE IF EXISTS public.distributors CASCADE;
TRUNCATE TABLE IF EXISTS public.retailer_users CASCADE;
TRUNCATE TABLE IF EXISTS public.retailers CASCADE;
TRUNCATE TABLE IF EXISTS public.transactions CASCADE;
TRUNCATE TABLE IF EXISTS public.payment_partners CASCADE;
TRUNCATE TABLE IF EXISTS public.user_providers CASCADE;
TRUNCATE TABLE IF EXISTS public.user_roles CASCADE;
TRUNCATE TABLE IF EXISTS public.role_resources CASCADE;
TRUNCATE TABLE IF EXISTS public.users CASCADE;
TRUNCATE TABLE IF EXISTS public.roles CASCADE;
TRUNCATE TABLE IF EXISTS public.resources CASCADE;
TRUNCATE TABLE IF EXISTS public.admins CASCADE;

-- Reset sequences (optional, but ensures IDs start from 1)
-- Note: TRUNCATE with CASCADE should reset sequences, but we'll do it explicitly
ALTER SEQUENCE IF EXISTS public.users_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.roles_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.resources_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.distributors_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.retailers_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.products_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.category_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.p_attributes_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.p_attribute_values_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.configurable_products_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.orders_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.invoices_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.invoice_line_items_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.transactions_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.s_ledger_id_seq RESTART WITH 1;
ALTER SEQUENCE IF EXISTS public.admins_id_seq RESTART WITH 1;

COMMIT;

