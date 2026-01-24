CREATE TABLE IF NOT EXISTS "permissions" (
    "id" UUID PRIMARY KEY,
    "resource" VARCHAR(255) NOT NULL,
    "action" VARCHAR(255) NOT NULL,
    "description" TEXT DEFAULT '',
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS permissions_resource_action_uidx
    ON permissions (resource, action);

