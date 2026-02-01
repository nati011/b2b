
CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID PRIMARY KEY,
    "external_id" VARCHAR(255),
    "email" VARCHAR(255),
    "phone_number" VARCHAR(255),
    "name" VARCHAR(255),
    "status" VARCHAR(50),
    "user_type" VARCHAR(50),
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE
);

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_users_email ON "users"(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON "users"(phone_number) WHERE phone_number IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_external_id ON "users"(external_id) WHERE external_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_status ON "users"(status) WHERE status IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_user_type ON "users"(user_type) WHERE user_type IS NOT NULL;
