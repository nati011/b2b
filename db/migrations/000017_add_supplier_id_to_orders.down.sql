-- Remove index
DROP INDEX IF EXISTS public.idx_orders_supplier_id;

-- Remove foreign key constraint
ALTER TABLE public.orders
DROP CONSTRAINT IF EXISTS fk_orders_supplier_id;

-- Remove supplier_id column
ALTER TABLE public.orders
DROP COLUMN IF EXISTS supplier_id;

