-- Create price_history table for auditing product price changes
CREATE TABLE IF NOT EXISTS public.price_history (
  id SERIAL PRIMARY KEY,
  product_id INT NOT NULL,
  old_price DECIMAL(10,2),
  new_price DECIMAL(10,2) NOT NULL,
  changed_by_user_id INT,
  changed_by_user_email VARCHAR(255),
  reason TEXT,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.price_history IS 'audit trail for product price changes';

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_price_history_product_id ON public.price_history(product_id);
CREATE INDEX IF NOT EXISTS idx_price_history_created_date ON public.price_history(created_date DESC);
CREATE INDEX IF NOT EXISTS idx_price_history_changed_by_user_id ON public.price_history(changed_by_user_id) WHERE changed_by_user_id IS NOT NULL;

