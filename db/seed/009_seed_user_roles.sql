-- Seed data for user_roles table
-- Assigns customer role to customer users and supplier role to supplier users
-- This must be run after users are created and roles are bootstrapped by the application
-- Development seed data

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

