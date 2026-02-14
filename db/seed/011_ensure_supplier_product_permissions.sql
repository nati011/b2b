-- Ensure Supplier Role Has Product Permissions
-- This script ensures the supplier role has all necessary product permissions
-- including products:create, products:view, products:update, products:delete
-- 
-- Note: The role_permissions table stores permissions directly with resource and action,
-- not via permission IDs. This script directly inserts into role_permissions.

-- Ensure the Supplier role has products:create permission
INSERT INTO role_permissions (role_id, resource, action)
SELECT 
    r.id,
    'products',
    'create'
FROM roles r
WHERE r.name = 'Supplier'
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id 
      AND rp.resource = 'products' 
      AND rp.action = 'create'
  )
ON CONFLICT (role_id, resource, action) DO NOTHING;

-- Ensure the Supplier role has products:view permission
INSERT INTO role_permissions (role_id, resource, action)
SELECT 
    r.id,
    'products',
    'view'
FROM roles r
WHERE r.name = 'Supplier'
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id 
      AND rp.resource = 'products' 
      AND rp.action = 'view'
  )
ON CONFLICT (role_id, resource, action) DO NOTHING;

-- Ensure the Supplier role has products:update permission
INSERT INTO role_permissions (role_id, resource, action)
SELECT 
    r.id,
    'products',
    'update'
FROM roles r
WHERE r.name = 'Supplier'
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id 
      AND rp.resource = 'products' 
      AND rp.action = 'update'
  )
ON CONFLICT (role_id, resource, action) DO NOTHING;

-- Ensure the Supplier role has products:delete permission
INSERT INTO role_permissions (role_id, resource, action)
SELECT 
    r.id,
    'products',
    'delete'
FROM roles r
WHERE r.name = 'Supplier'
  AND NOT EXISTS (
    SELECT 1 FROM role_permissions rp
    WHERE rp.role_id = r.id 
      AND rp.resource = 'products' 
      AND rp.action = 'delete'
  )
ON CONFLICT (role_id, resource, action) DO NOTHING;

-- Verification query (uncomment to verify permissions are assigned)
-- SELECT 
--     r.name as role_name,
--     rp.resource,
--     rp.action
-- FROM roles r
-- JOIN role_permissions rp ON r.id = rp.role_id
-- WHERE r.name = 'Supplier'
--   AND rp.resource = 'products'
-- ORDER BY rp.action;

