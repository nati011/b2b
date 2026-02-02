-- Create bank_accounts table for supplier bank account information
CREATE TABLE IF NOT EXISTS public.bank_accounts (
  id SERIAL PRIMARY KEY,
  supplier_id INT NOT NULL,
  bank_name VARCHAR(255) NOT NULL,
  account_number VARCHAR(100) NOT NULL,
  account_holder_name VARCHAR(255) NOT NULL,
  branch_name VARCHAR(255),
  account_type VARCHAR(50), -- e.g., 'checking', 'savings', 'current'
  is_primary BOOLEAN DEFAULT FALSE,
  is_active BOOLEAN DEFAULT TRUE,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.bank_accounts IS 'stores supplier bank account information for payments';

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_bank_accounts_supplier_id ON public.bank_accounts(supplier_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_active ON public.bank_accounts(is_active) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_primary ON public.bank_accounts(is_primary) WHERE is_deleted = FALSE AND is_primary = TRUE;

-- Ensure only one primary account per supplier
CREATE UNIQUE INDEX IF NOT EXISTS idx_bank_accounts_supplier_primary 
ON public.bank_accounts(supplier_id) 
WHERE is_deleted = FALSE AND is_primary = TRUE;

