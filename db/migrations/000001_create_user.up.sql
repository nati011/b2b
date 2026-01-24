
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
