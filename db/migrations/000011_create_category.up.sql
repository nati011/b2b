CREATE TABLE IF NOT EXISTS public.category (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE
);

COMMENT ON TABLE public.category IS 'category tags for products';
