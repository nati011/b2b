-- Seed data for supplier users with credentials and roles
-- This file creates supplier users, their authentication credentials, and assigns supplier roles
-- Development seed data
--
-- Note: This must be run after:
--   1. Users table exists (migration 000001)
--   2. Roles table exists and roles are bootstrapped (migration 000002 + config/roles.yaml)
--   3. Basic auth credentials table exists (migration 000015)
--   4. User roles table exists (migration 000007)
--
-- Default password for all supplier users: password123
-- Password hash (bcrypt): $2a$10$dk1d.WykwJNoDGWUti1tFemAFgv5/TozdTGLlzY9jFsk1xAHc.sn.

-- ============================================================================
-- STEP 1: Create Supplier Users
-- ============================================================================
INSERT INTO users (id, external_id, email, phone_number, name, status, user_type, created_at, updated_at)
VALUES
  -- Supplier User 1
  (
    '550e8400-e29b-41d4-a716-446655440002',
    'EXT-SUPPLIER-001',
    'supplier1@b2b.local',
    '+251911000002',
    'Supplier One',
    'active',
    'supplier',
    NOW(),
    NOW()
  ),
  -- Supplier User 2
  (
    '550e8400-e29b-41d4-a716-446655440003',
    'EXT-SUPPLIER-002',
    'supplier2@b2b.local',
    '+251911000003',
    'Supplier Two',
    'active',
    'supplier',
    NOW(),
    NOW()
  ),
  -- Supplier User 3 (Additional supplier for testing)
  (
    '550e8400-e29b-41d4-a716-446655440010',
    'EXT-SUPPLIER-003',
    'supplier3@b2b.local',
    '+251911000010',
    'Supplier Three',
    'active',
    'supplier',
    NOW(),
    NOW()
  )
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- STEP 2: Create Authentication Credentials for Supplier Users
-- ============================================================================
-- Password for all: password123
-- Bcrypt hash: $2a$10$dk1d.WykwJNoDGWUti1tFemAFgv5/TozdTGLlzY9jFsk1xAHc.sn.
INSERT INTO basic_auth_credentials (id, username, password, user_id, active, created_at, updated_at)
VALUES
  -- Credentials for supplier1@b2b.local
  (
    gen_random_uuid(),
    'supplier1@b2b.local',
    '$2a$10$dk1d.WykwJNoDGWUti1tFemAFgv5/TozdTGLlzY9jFsk1xAHc.sn.', -- password123
    '550e8400-e29b-41d4-a716-446655440002',
    true,
    NOW(),
    NOW()
  ),
  -- Credentials for supplier2@b2b.local
  (
    gen_random_uuid(),
    'supplier2@b2b.local',
    '$2a$10$dk1d.WykwJNoDGWUti1tFemAFgv5/TozdTGLlzY9jFsk1xAHc.sn.', -- password123
    '550e8400-e29b-41d4-a716-446655440003',
    true,
    NOW(),
    NOW()
  ),
  -- Credentials for supplier3@b2b.local
  (
    gen_random_uuid(),
    'supplier3@b2b.local',
    '$2a$10$dk1d.WykwJNoDGWUti1tFemAFgv5/TozdTGLlzY9jFsk1xAHc.sn.', -- password123
    '550e8400-e29b-41d4-a716-446655440010',
    true,
    NOW(),
    NOW()
  )
ON CONFLICT DO NOTHING;

-- ============================================================================
-- STEP 3: Assign Supplier Role to Supplier Users
-- ============================================================================
-- This assigns the 'Supplier' role to all users with user_type='supplier'
-- The role must exist in the roles table (bootstrapped from config/roles.yaml)
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.user_type = 'supplier'
  AND r.name = 'Supplier'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = u.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- ============================================================================
-- Verification Query (commented out - uncomment to verify)
-- ============================================================================
-- SELECT 
--   u.id,
--   u.email,
--   u.name,
--   u.user_type,
--   u.status,
--   r.name as role_name,
--   CASE WHEN bac.id IS NOT NULL THEN 'Yes' ELSE 'No' END as has_credentials
-- FROM users u
-- LEFT JOIN user_roles ur ON u.id = ur.user_id
-- LEFT JOIN roles r ON ur.role_id = r.id
-- LEFT JOIN basic_auth_credentials bac ON u.id = bac.user_id AND bac.active = true
-- WHERE u.user_type = 'supplier'
-- ORDER BY u.email;






