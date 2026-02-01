-- Seed data for orders and order items tables
-- Development seed data

-- Insert orders with cart_snapshot
INSERT INTO public.orders (customer_id, status, payment_status, delivery_status, confirmation_status, total, customer_snapshot, shipping_address_snapshot, billing_address_snapshot, cart_snapshot, created_date, last_modified, is_deleted)
VALUES
  -- Order 1: Customer 1
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
   )::jsonb,
   NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', FALSE),
  
  -- Order 2: Customer 2
  (2, 'confirmed', 'paid', 'processing', 'confirmed', 370.00,
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
     WHERE p.external_id IN ('PROD-OFF-001', 'PROD-OFF-002')
   )::jsonb,
   NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', FALSE),
  
  -- Order 3: Customer 3
  (3, 'shipped', 'paid', 'in_transit', 'confirmed', 9700.00,
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
     WHERE p.external_id IN ('PROD-TOOL-001', 'PROD-TOOL-002')
   )::jsonb,
   NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', FALSE),
  
  -- Order 4: Customer 4
  (4, 'delivered', 'paid', 'delivered', 'confirmed', 4500.00,
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
     WHERE p.external_id = 'PROD-FURN-001'
   )::jsonb,
   NOW() - INTERVAL '7 days', NOW() - INTERVAL '1 day', FALSE),
  
  -- Order 5: Customer 5
  (5, 'pending', 'unpaid', 'not_shipped', 'pending', 4550.00,
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
     WHERE p.external_id IN ('PROD-SAFE-001', 'PROD-SAFE-002')
   )::jsonb,
   NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', FALSE)
ON CONFLICT DO NOTHING;
