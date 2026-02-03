-- Add supplier_id column to orders table
ALTER TABLE public.orders
ADD COLUMN IF NOT EXISTS supplier_id INT;

-- Add foreign key constraint
ALTER TABLE public.orders
ADD CONSTRAINT fk_orders_supplier_id
FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id) ON DELETE SET NULL;

-- Add index for supplier_id
CREATE INDEX IF NOT EXISTS idx_orders_supplier_id ON public.orders(supplier_id) WHERE is_deleted = FALSE;

COMMENT ON COLUMN public.orders.supplier_id IS 'supplier_id of the first product in the order';


