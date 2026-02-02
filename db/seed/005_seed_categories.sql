-- Seed data for category table
-- Development seed data

INSERT INTO public.category (name, created_date, last_modified, is_deleted)
VALUES
  ('Electronics', NOW(), NOW(), FALSE),
  ('Office Supplies', NOW(), NOW(), FALSE),
  ('Tools & Equipment', NOW(), NOW(), FALSE),
  ('Furniture', NOW(), NOW(), FALSE),
  ('Safety & Security', NOW(), NOW(), FALSE),
  ('Industrial', NOW(), NOW(), FALSE),
  ('Construction Materials', NOW(), NOW(), FALSE),
  ('Cleaning Supplies', NOW(), NOW(), FALSE),
  ('Packaging Materials', NOW(), NOW(), FALSE),
  ('Food & Beverages', NOW(), NOW(), FALSE)
ON CONFLICT DO NOTHING;










