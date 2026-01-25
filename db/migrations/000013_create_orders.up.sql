CREATE TABLE IF NOT EXISTS public.orders (
  id SERIAL PRIMARY KEY,
  customer_id INT NOT NULL,
  status VARCHAR(255),
  payment_status VARCHAR(255),
  delivery_status VARCHAR(255),
  confirmation_status VARCHAR(255),
  total JSONB,
  customer_snapshot JSONB,
  shipping_address_snapshot JSONB,
  billing_address_snapshot JSONB,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.orders IS 'stores orders';

CREATE TABLE IF NOT EXISTS public.o_items (
  order_id INT,
  product_id INT,
  quantity INT,
  price JSONB,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE,
  FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.o_items IS 'stores order items';
