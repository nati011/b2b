CREATE TABLE IF NOT EXISTS public."templates"
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE,
    html BYTEA,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

COMMENT ON TABLE public."template" IS 'stores templates.';
