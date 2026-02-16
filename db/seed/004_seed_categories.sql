-- Seed data for category table
-- Industrial and machinery focused

INSERT INTO public.category (name, created_date, last_modified, is_deleted)
VALUES
  ('Industrial Equipment', NOW(), NOW(), FALSE),
  ('Machinery', NOW(), NOW(), FALSE),
  ('Safety Equipment', NOW(), NOW(), FALSE),
  ('Construction Materials', NOW(), NOW(), FALSE),
  ('Metalworking', NOW(), NOW(), FALSE),
  ('Electrical Equipment', NOW(), NOW(), FALSE),
  ('Hydraulics & Pneumatics', NOW(), NOW(), FALSE),
  ('Tools & Abrasives', NOW(), NOW(), FALSE)
ON CONFLICT DO NOTHING;
