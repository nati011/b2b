CREATE TABLE IF NOT EXISTS public.customers (
  id SERIAL PRIMARY KEY,
  full_name VARCHAR(255),
  status VARCHAR(50),
  city VARCHAR(100),
  region VARCHAR(100),
  woreda VARCHAR(100),
  phone_number VARCHAR(50),
  email VARCHAR(255),
  is_active BOOLEAN DEFAULT FALSE,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE
);

COMMENT ON TABLE public.customers IS 'stores customer specific information (not user)';

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_customers_email ON public.customers(email) WHERE is_deleted = FALSE AND email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_customers_phone_number ON public.customers(phone_number) WHERE is_deleted = FALSE AND phone_number IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_customers_status ON public.customers(status) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_customers_is_active ON public.customers(is_active) WHERE is_deleted = FALSE;
