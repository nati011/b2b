-- Allowed actions for each resource.
CREATE TABLE resource_actions (
    id UUID PRIMARY KEY,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    action VARCHAR(128) NOT NULL,
    description TEXT NOT NULL,
    deprecated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Unique constraint on resource_id and action (actions are normalized to lowercase in application code)
ALTER TABLE resource_actions
    ADD CONSTRAINT resource_actions_resource_action_uc
    UNIQUE (resource_id, action);

CREATE INDEX resource_actions_resource_id_idx ON resource_actions (resource_id);

