-- Fix admin user and credentials
-- This ensures the admin user exists with the correct email and credentials

-- Step 1: Create or update the admin user
INSERT INTO users (id, external_id, email, phone_number, name, status, user_type, created_at, updated_at)
VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'EXT-ADMIN-admin', 'admin-admin@admin.local', '+251911000001', 'Super Administrator', 'active', 'officer', NOW(), NOW())
ON CONFLICT (id) DO UPDATE
SET 
  email = EXCLUDED.email,
  name = EXCLUDED.name,
  status = EXCLUDED.status,
  user_type = EXCLUDED.user_type,
  updated_at = NOW();

-- Step 2: Create or update credentials for admin-admin@admin.local
-- Password: changeme
-- Bcrypt hash: $2a$10$1OcPvbmRKzgnIwwWAuMq3.MSj5H.5MA9DYt8/QaqzdY2kQkOBeolS
-- Delete existing credential if it exists (to avoid conflicts)
DELETE FROM basic_auth_credentials WHERE username = 'admin-admin@admin.local';

-- Insert the credential
INSERT INTO basic_auth_credentials (id, username, password, user_id, active, created_at, updated_at)
VALUES
  (
    gen_random_uuid(),
    'admin-admin@admin.local',
    '$2a$10$1OcPvbmRKzgnIwwWAuMq3.MSj5H.5MA9DYt8/QaqzdY2kQkOBeolS',
    '550e8400-e29b-41d4-a716-446655440001',
    true,
    NOW(),
    NOW()
  );

-- Step 3: Ensure admin role is assigned (if roles table exists)
-- The role ID in roles.yaml is 'admin', but it's stored as a UUID in the database
-- The bootstrap function generates a deterministic UUID from the role name if needed
-- We'll find the role by name "Admin" instead
DO $$
DECLARE
    admin_role_id UUID;
BEGIN
    -- Try to find role by name "Admin" (case-insensitive)
    SELECT id INTO admin_role_id FROM roles WHERE LOWER(name) = LOWER('Admin') LIMIT 1;
    IF admin_role_id IS NOT NULL THEN
        INSERT INTO user_roles (user_id, role_id)
        VALUES ('550e8400-e29b-41d4-a716-446655440001', admin_role_id)
        ON CONFLICT (user_id, role_id) DO NOTHING;
    END IF;
END $$;

