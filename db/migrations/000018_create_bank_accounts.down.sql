-- Drop indexes
DROP INDEX IF EXISTS public.idx_bank_accounts_supplier_primary;
DROP INDEX IF EXISTS public.idx_bank_accounts_is_primary;
DROP INDEX IF EXISTS public.idx_bank_accounts_is_active;
DROP INDEX IF EXISTS public.idx_bank_accounts_supplier_id;

-- Drop table
DROP TABLE IF EXISTS public.bank_accounts;

