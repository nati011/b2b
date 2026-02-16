ALTER TABLE public.orders
  ADD COLUMN IF NOT EXISTS delivery_address VARCHAR(500) NOT NULL DEFAULT '';

COMMENT ON COLUMN public.orders.delivery_address IS 'customer delivery address for the order (mandatory)';
