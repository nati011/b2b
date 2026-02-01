-- Seed data for users table
-- Development seed data

INSERT INTO users (id, external_id, email, phone_number, name, status, user_type, created_at, updated_at)
VALUES
  ('550e8400-e29b-41d4-a716-446655440002', 'EXT-USER-002', 'supplier1@b2b.local', '+251911000002', 'Supplier One', 'active', 'supplier', NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440003', 'EXT-USER-003', 'supplier2@b2b.local', '+251911000003', 'Supplier Two', 'active', 'supplier', NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440004', 'EXT-USER-004', 'customer1@b2b.local', '+251911000004', 'Customer One', 'active', 'customer', NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440005', 'EXT-USER-005', 'customer2@b2b.local', '+251911000005', 'Customer Two', 'active', 'customer', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;




