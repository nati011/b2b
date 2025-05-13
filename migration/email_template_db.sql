-- Version 0.1 (corrected)

CREATE TABLE IF NOT EXISTS public.templates
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    html TEXT,  -- Changed from BYTEA to TEXT for HTML content
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

COMMENT ON TABLE public.templates IS 'Stores email templates.';

-- email_templates ---------------------------

    -- Write
CREATE OR REPLACE PROCEDURE public.insert_template(
   template_name VARCHAR(255),
   template_html TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.templates (name, html)
    VALUES (template_name, template_html);
    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE public.update_template(
   template_name VARCHAR(255),
   new_html TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.templates
    SET html = new_html,
        last_modified = CURRENT_TIMESTAMP
    WHERE name = template_name;
    COMMIT;
END;
$$;

    -- Read
CREATE OR REPLACE FUNCTION public.get_template(
   template_name VARCHAR(255)
)
RETURNS SETOF public.templates
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM public.templates t
    WHERE t.name = template_name
      AND t.is_deleted = FALSE;
END;
$$;

