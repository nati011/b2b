CREATE TABLE IF NOT EXISTS public.products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  description VARCHAR(255),
  external_id VARCHAR(255),
  attributes JSONB,
  unit VARCHAR(50),
  is_active BOOLEAN DEFAULT FALSE,
  supplier_id INT,
  price DECIMAL(12,2),
  total_quantity INT,
  reserved_quantity INT DEFAULT 0,
  available_quantity INT GENERATED ALWAYS AS (total_quantity - reserved_quantity) STORED,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (supplier_id) REFERENCES public.suppliers(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.product_categories (
  id SERIAL PRIMARY KEY,
  product_id INT NOT NULL,
  category_id INT NOT NULL,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE,
  FOREIGN KEY (category_id) REFERENCES public.category(id) ON DELETE CASCADE
);

COMMENT ON TABLE public.products IS 'stores products';
COMMENT ON TABLE public.product_categories IS 'maps products to multiple categories';