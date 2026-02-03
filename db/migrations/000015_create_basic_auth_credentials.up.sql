-- Create basic_auth_credentials table
CREATE TABLE IF NOT EXISTS basic_auth_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_basic_auth_credentials_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create index on username for faster lookups
CREATE INDEX IF NOT EXISTS idx_basic_auth_credentials_username ON basic_auth_credentials(LOWER(username));

-- Create index on user_id
CREATE INDEX IF NOT EXISTS idx_basic_auth_credentials_user_id ON basic_auth_credentials(user_id);

-- Create unique constraint on username (case-insensitive would be ideal but PostgreSQL doesn't support it directly)
-- We'll rely on application-level case-insensitive matching








