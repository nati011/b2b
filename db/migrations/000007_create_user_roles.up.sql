CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_roles_users FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_roles FOREIGN KEY (role_id)
        REFERENCES roles (id) ON DELETE CASCADE
);

-- Indexes for reverse lookups (role_id is already in PK but useful for queries filtering by role)
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);

