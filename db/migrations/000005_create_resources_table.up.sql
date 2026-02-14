-- Resources are permissionable surfaces identified by a code + service.
CREATE TABLE resources (
    id UUID PRIMARY KEY,
    code VARCHAR(128) NOT NULL,
    service VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    deprecated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX resources_code_uidx ON resources (LOWER(code));

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_resources_service ON resources(service) WHERE deprecated_at IS NULL;