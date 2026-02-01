-- Seed data for basic_auth_credentials table
-- Development seed data
-- 
-- This creates credentials for the seeded users
-- Password for all: password123
--
-- Note: The bcrypt hash for "password123" is pre-generated
-- If you need a different password, generate a new hash using:
--   python3 -c "import bcrypt; print(bcrypt.hashpw(b'your_password', bcrypt.gensalt()).decode())"

-- Password hash for "password123" (bcrypt)
-- This hash was generated with: bcrypt.hashpw(b'password123', bcrypt.gensalt())
INSERT INTO basic_auth_credentials (id, username, password, user_id, active, created_at, updated_at)
VALUES
  -- Credentials for customer1@b2b.local (User ID: 550e8400-e29b-41d4-a716-446655440004)
  (
    gen_random_uuid(),
    'customer1@b2b.local',
    '$2a$10$dk1d.WykwJNoDGWUti1tFemAFgv5/TozdTGLlzY9jFsk1xAHc.sn.', -- password123
    '550e8400-e29b-41d4-a716-446655440004',
    true,
    NOW(),
    NOW()
  ),
  -- Credentials for customer2@b2b.local (User ID: 550e8400-e29b-41d4-a716-446655440005)
  (
    gen_random_uuid(),
    'customer2@b2b.local',
    '$2a$10$dk1d.WykwJNoDGWUti1tFemAFgv5/TozdTGLlzY9jFsk1xAHc.sn.', -- password123
    '550e8400-e29b-41d4-a716-446655440005',
    true,
    NOW(),
    NOW()
  ),
  -- Credentials for supplier1@b2b.local (User ID: 550e8400-e29b-41d4-a716-446655440002)
  (
    gen_random_uuid(),
    'supplier1@b2b.local',
    '$2a$10$dk1d.WykwJNoDGWUti1tFemAFgv5/TozdTGLlzY9jFsk1xAHc.sn.', -- password123
    '550e8400-e29b-41d4-a716-446655440002',
    true,
    NOW(),
    NOW()
  ),
  -- Credentials for supplier2@b2b.local (User ID: 550e8400-e29b-41d4-a716-446655440003)
  (
    gen_random_uuid(),
    'supplier2@b2b.local',
    '$2a$10$dk1d.WykwJNoDGWUti1tFemAFgv5/TozdTGLlzY9jFsk1xAHc.sn.', -- password123
    '550e8400-e29b-41d4-a716-446655440003',
    true,
    NOW(),
    NOW()
  )
ON CONFLICT DO NOTHING;

