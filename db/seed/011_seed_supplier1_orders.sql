-- Seed data for orders containing products from supplier1@b2b.local (supplier_id = 1)
-- Development seed data for testing supplier-based order fetch
-- Products from supplier 1: PROD-ELEC-001, PROD-ELEC-002, PROD-ELEC-003

-- Insert orders with cart_snapshot containing products from supplier 1
INSERT INTO public.orders (customer_id, status, payment_status, delivery_status, confirmation_status, total, customer_snapshot, shipping_address_snapshot, billing_address_snapshot, cart_snapshot, created_date, last_modified, is_deleted)
VALUES
  -- Order 1: Customer 1 - Pending order with Laptop and Mouse (Supplier 1 products)
  (1, 'pending', 'unpaid', 'not_shipped', 'pending', 45850.00, 
   '{"id": 1, "full_name": "Abebe Bekele", "email": "abebe@example.local", "phone": "+251911200001"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Bole", "address": "Bole Sub-city, Street 123"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Bole", "address": "Bole Sub-city, Street 123"}'::jsonb,
   (
     SELECT jsonb_agg(
       jsonb_build_object(
         'product_id', p.id,
         'quantity', 1,
         'price', p.price
       )
     )
     FROM public.products p
     WHERE p.external_id IN ('PROD-ELEC-001', 'PROD-ELEC-002')
       AND p.supplier_id = 1
       AND p.is_deleted = FALSE
   )::jsonb,
   NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', FALSE),
  
  -- Order 2: Customer 2 - Confirmed order with Keyboard (Supplier 1 product)
  (2, 'confirmed', 'paid', 'processing', 'confirmed', 3500.00,
   '{"id": 2, "full_name": "Meron Tadesse", "email": "meron@example.local", "phone": "+251911200002"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Kirkos", "address": "Kirkos Sub-city, Avenue 456"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Kirkos", "address": "Kirkos Sub-city, Avenue 456"}'::jsonb,
   (
     SELECT jsonb_agg(
       jsonb_build_object(
         'product_id', p.id,
         'quantity', 2,
         'price', p.price
       )
     )
     FROM public.products p
     WHERE p.external_id = 'PROD-ELEC-003'
       AND p.supplier_id = 1
       AND p.is_deleted = FALSE
   )::jsonb,
   NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', FALSE),
  
  -- Order 3: Customer 3 - Shipped order with all three products (Supplier 1 products)
  (3, 'shipped', 'paid', 'in_transit', 'confirmed', 49200.00,
   '{"id": 3, "full_name": "Yonas Alemayehu", "email": "yonas@example.local", "phone": "+251911200003"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Arada", "address": "Arada Sub-city, Road 789"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Arada", "address": "Arada Sub-city, Road 789"}'::jsonb,
   (
     SELECT jsonb_agg(
       jsonb_build_object(
         'product_id', p.id,
         'quantity', 1,
         'price', p.price
       )
     )
     FROM public.products p
     WHERE p.external_id IN ('PROD-ELEC-001', 'PROD-ELEC-002', 'PROD-ELEC-003')
       AND p.supplier_id = 1
       AND p.is_deleted = FALSE
   )::jsonb,
   NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', FALSE),
  
  -- Order 4: Customer 4 - Delivered order with Laptop (Supplier 1 product)
  (4, 'delivered', 'paid', 'delivered', 'confirmed', 45000.00,
   '{"id": 4, "full_name": "Sara Getachew", "email": "sara@example.local", "phone": "+251911200004"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Lideta", "address": "Lideta Sub-city, Lane 321"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Lideta", "address": "Lideta Sub-city, Lane 321"}'::jsonb,
   (
     SELECT jsonb_agg(
       jsonb_build_object(
         'product_id', p.id,
         'quantity', 1,
         'price', p.price
       )
     )
     FROM public.products p
     WHERE p.external_id = 'PROD-ELEC-001'
       AND p.supplier_id = 1
       AND p.is_deleted = FALSE
   )::jsonb,
   NOW() - INTERVAL '7 days', NOW() - INTERVAL '1 day', FALSE),
  
  -- Order 5: Customer 5 - Pending order with Mouse (Supplier 1 product)
  (5, 'pending', 'unpaid', 'not_shipped', 'pending', 850.00,
   '{"id": 5, "full_name": "Daniel Haile", "email": "daniel@example.local", "phone": "+251911200005"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Nifas Silk", "address": "Nifas Silk Sub-city, Street 654"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Nifas Silk", "address": "Nifas Silk Sub-city, Street 654"}'::jsonb,
   (
     SELECT jsonb_agg(
       jsonb_build_object(
         'product_id', p.id,
         'quantity', 1,
         'price', p.price
       )
     )
     FROM public.products p
     WHERE p.external_id = 'PROD-ELEC-002'
       AND p.supplier_id = 1
       AND p.is_deleted = FALSE
   )::jsonb,
   NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', FALSE),
  
  -- Order 6: Customer 1 - Recent order with Keyboard and Mouse (Supplier 1 products)
  (1, 'confirmed', 'paid', 'processing', 'confirmed', 4350.00,
   '{"id": 1, "full_name": "Abebe Bekele", "email": "abebe@example.local", "phone": "+251911200001"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Bole", "address": "Bole Sub-city, Street 123"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Bole", "address": "Bole Sub-city, Street 123"}'::jsonb,
   (
     SELECT jsonb_agg(
       jsonb_build_object(
         'product_id', p.id,
         'quantity', 1,
         'price', p.price
       )
     )
     FROM public.products p
     WHERE p.external_id IN ('PROD-ELEC-002', 'PROD-ELEC-003')
       AND p.supplier_id = 1
       AND p.is_deleted = FALSE
   )::jsonb,
   NOW() - INTERVAL '6 hours', NOW() - INTERVAL '6 hours', FALSE),
  
  -- Order 7: Customer 2 - Cancelled order with Laptop (Supplier 1 product)
  (2, 'cancelled', 'refunded', 'not_shipped', 'cancelled', 45000.00,
   '{"id": 2, "full_name": "Meron Tadesse", "email": "meron@example.local", "phone": "+251911200002"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Kirkos", "address": "Kirkos Sub-city, Avenue 456"}'::jsonb,
   '{"city": "Addis Ababa", "region": "Addis Ababa", "woreda": "Kirkos", "address": "Kirkos Sub-city, Avenue 456"}'::jsonb,
   (
     SELECT jsonb_agg(
       jsonb_build_object(
         'product_id', p.id,
         'quantity', 1,
         'price', p.price
       )
     )
     FROM public.products p
     WHERE p.external_id = 'PROD-ELEC-001'
       AND p.supplier_id = 1
       AND p.is_deleted = FALSE
   )::jsonb,
   NOW() - INTERVAL '10 days', NOW() - INTERVAL '8 days', FALSE)
ON CONFLICT DO NOTHING;

-- Verify the orders were created
-- SELECT 
--   o.id,
--   o.customer_id,
--   o.status,
--   o.total,
--   o.cart_snapshot,
--   o.created_date
-- FROM orders o
-- WHERE o.is_deleted = FALSE
--   AND EXISTS (
--     SELECT 1 
--     FROM jsonb_array_elements(o.cart_snapshot) AS item
--     JOIN products p ON (
--       ((item->>'product_id') IS NOT NULL AND (item->>'product_id')::int = p.id)
--       OR ((item->>'ProductID') IS NOT NULL AND (item->>'ProductID')::int = p.id)
--       OR ((item->>'productId') IS NOT NULL AND (item->>'productId')::int = p.id)
--     )
--     WHERE p.supplier_id = 1 AND p.is_deleted = FALSE
--   )
-- ORDER BY o.created_date DESC;


