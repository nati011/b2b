CREATE TABLE IF NOT EXISTS "idempotency_keys" (
    "idempotency_key" VARCHAR(200) PRIMARY KEY,
    "method"          VARCHAR(10)  NOT NULL,
    "path"            TEXT         NOT NULL,
    "request_hash"    VARCHAR(128) NOT NULL,
    "processing"      BOOLEAN      NOT NULL DEFAULT TRUE,
    "response_status" INTEGER      NULL,
    "response_headers" JSONB       NULL,
    "response_body"   BYTEA        NULL,
    "created_at"      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    "updated_at"      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_idempotency_keys_updated_at
    ON idempotency_keys (updated_at DESC);

