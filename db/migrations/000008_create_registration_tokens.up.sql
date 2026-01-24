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

