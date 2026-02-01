CREATE TABLE IF NOT EXISTS registration_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_registration_tokens_users FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
);

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_registration_tokens_user_id ON registration_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_registration_tokens_used ON registration_tokens(used) WHERE used = FALSE;
CREATE INDEX IF NOT EXISTS idx_registration_tokens_expires_at ON registration_tokens(expires_at) WHERE used = FALSE;

