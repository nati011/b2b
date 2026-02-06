-- Seed data for customers table
-- Development seed data

INSERT INTO public.customers (full_name, status, city, region, woreda, phone_number, email, is_active, created_date, last_modified, is_deleted)
VALUES
  ('Abebe Bekele', 'active', 'Addis Ababa', 'Addis Ababa', 'Bole', '+251911200001', 'abebe@example.local', TRUE, NOW(), NOW(), FALSE),
  ('Meron Tadesse', 'active', 'Addis Ababa', 'Addis Ababa', 'Kirkos', '+251911200002', 'meron@example.local', TRUE, NOW(), NOW(), FALSE),
  ('Yonas Alemayehu', 'active', 'Addis Ababa', 'Addis Ababa', 'Arada', '+251911200003', 'yonas@example.local', TRUE, NOW(), NOW(), FALSE),
  ('Sara Getachew', 'active', 'Addis Ababa', 'Addis Ababa', 'Lideta', '+251911200004', 'sara@example.local', TRUE, NOW(), NOW(), FALSE),
  ('Daniel Haile', 'active', 'Addis Ababa', 'Addis Ababa', 'Nifas Silk', '+251911200005', 'daniel@example.local', TRUE, NOW(), NOW(), FALSE),
  ('Tigist Worku', 'active', 'Addis Ababa', 'Addis Ababa', 'Yeka', '+251911200006', 'tigist@example.local', TRUE, NOW(), NOW(), FALSE),
  ('Kebede Tesfaye', 'active', 'Addis Ababa', 'Addis Ababa', 'Gulele', '+251911200007', 'kebede@example.local', TRUE, NOW(), NOW(), FALSE),
  ('Alemitu Fekadu', 'active', 'Addis Ababa', 'Addis Ababa', 'Kolfe', '+251911200008', 'alemitu@example.local', TRUE, NOW(), NOW(), FALSE),
  ('Bereket Assefa', 'active', 'Addis Ababa', 'Addis Ababa', 'Addis Ketema', '+251911200009', 'bereket@example.local', TRUE, NOW(), NOW(), FALSE),
  ('Hanna Solomon', 'active', 'Addis Ababa', 'Addis Ababa', 'Akaki', '+251911200010', 'hanna@example.local', TRUE, NOW(), NOW(), FALSE)
ON CONFLICT DO NOTHING;














