-- Seed data for suppliers table
-- Development seed data
-- Note: support_email must match user email for supplier filtering to work
-- Users: supplier1@b2b.local, supplier2@b2b.local, supplier3@b2b.local

INSERT INTO public.suppliers (business_name, status, support_email, support_phone, is_active, created_date, last_modified, is_deleted)
VALUES
  ('TechSupply Co.', 'active', 'supplier1@b2b.local', '+251911100001', TRUE, NOW(), NOW(), FALSE),
  ('Industrial Solutions Ltd', 'active', 'supplier2@b2b.local', '+251911100002', TRUE, NOW(), NOW(), FALSE),
  ('Global Equipment Inc', 'active', 'supplier3@b2b.local', '+251911100003', TRUE, NOW(), NOW(), FALSE),
  ('Quality Tools Ltd', 'active', 'sales@qualitytools.local', '+251911100004', TRUE, NOW(), NOW(), FALSE),
  ('Professional Supplies', 'active', 'support@prosupplies.local', '+251911100005', TRUE, NOW(), NOW(), FALSE),
  ('Business Essentials Co', 'active', 'info@bizessentials.local', '+251911100006', TRUE, NOW(), NOW(), FALSE),
  ('Office Depot Ethiopia', 'active', 'orders@officedepot.local', '+251911100007', TRUE, NOW(), NOW(), FALSE),
  ('Construction Materials Hub', 'active', 'sales@constructionhub.local', '+251911100008', TRUE, NOW(), NOW(), FALSE)
ON CONFLICT DO NOTHING;

