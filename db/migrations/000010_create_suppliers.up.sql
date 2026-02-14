CREATE TABLE IF NOT EXISTS public.suppliers (
  id SERIAL PRIMARY KEY,
  business_name VARCHAR(255),
  status VARCHAR(50),
  support_email VARCHAR(255),
  support_phone VARCHAR(50),
  is_active BOOLEAN DEFAULT FALSE,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE
);

COMMENT ON TABLE public.suppliers IS 'stores supplier specific information (not user)';

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_suppliers_status ON public.suppliers(status) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_suppliers_is_active ON public.suppliers(is_active) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_suppliers_support_email ON public.suppliers(support_email) WHERE is_deleted = FALSE AND support_email IS NOT NULL;