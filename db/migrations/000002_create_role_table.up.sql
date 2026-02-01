CREATE TABLE IF NOT EXISTS "roles" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT DEFAULT '',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS role_permissions (
    "role_id" UUID NOT NULL,
    "resource" VARCHAR(255) NOT NULL,
    "action" VARCHAR(255) NOT NULL,
    PRIMARY KEY (role_id, resource, action),
    CONSTRAINT fk_role_permissions_roles FOREIGN KEY (role_id)
        REFERENCES roles (id) ON DELETE CASCADE
);

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_roles_name ON "roles"("name");
CREATE INDEX IF NOT EXISTS idx_role_permissions_resource_action ON role_permissions("resource", "action");