CREATE TABLE IF NOT EXISTS public.orders (
  id SERIAL PRIMARY KEY,
  customer_id INT NOT NULL,
  status VARCHAR(255),
  payment_status VARCHAR(255),
  delivery_status VARCHAR(255),
  confirmation_status VARCHAR(255),
  total DECIMAL(10,2),
  customer_snapshot JSONB,
  cart_snapshot JSONB,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.orders IS 'stores orders';

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON public.orders(customer_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_orders_status ON public.orders(status) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_orders_customer_status ON public.orders(customer_id, status) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_orders_created_date ON public.orders(created_date DESC) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_orders_customer_created ON public.orders(customer_id, created_date DESC) WHERE is_deleted = FALSE;
