-- Seed data for user_roles table
-- Assigns roles to users based on their user_type
-- This must be run after users are created and roles are bootstrapped by the application
-- Development seed data

-- Assign Admin role to users with user_type='officer'
-- Use case-insensitive matching for role name
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.user_type = 'officer'
  AND LOWER(r.name) = LOWER('Admin')
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = u.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Explicitly assign Admin role to admin-admin@admin.local (fallback)
-- This ensures the admin user always gets the Admin role even if user_type check fails
DO $$
DECLARE
    admin_user_id UUID := '550e8400-e29b-41d4-a716-446655440001';
    admin_role_id UUID;
BEGIN
    -- Find Admin role by name (case-insensitive)
    SELECT id INTO admin_role_id FROM roles WHERE LOWER(name) = LOWER('Admin') LIMIT 1;
    IF admin_role_id IS NOT NULL THEN
        INSERT INTO user_roles (user_id, role_id)
        VALUES (admin_user_id, admin_role_id)
        ON CONFLICT (user_id, role_id) DO NOTHING;
    END IF;
END $$;

-- Assign customer role to users with user_type='customer'
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.user_type = 'customer'
  AND r.name = 'Customer'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = u.id AND ur.role_id = r.id
  )
ON CONFLICT (user_id, role_id) DO NOTHING;

-- Assign supplier role to users with user_type='supplier'
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

