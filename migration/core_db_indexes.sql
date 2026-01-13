-- Performance indexes for catalogue API optimization
-- Version: 0.2

-- Index for products by distributor and active status (used in catalogue queries)
CREATE INDEX IF NOT EXISTS idx_products_distributor_active 
ON public."products"(distributor_id, is_active) 
WHERE is_active = true;

-- Index for configurable products by availability status
CREATE INDEX IF NOT EXISTS idx_configurable_products_available 
ON public."configurable_products"(is_available) 
WHERE is_available = true;

-- Index for product-category relationships (used in filtering)
CREATE INDEX IF NOT EXISTS idx_p_category_product_id 
ON public."p_category"(product_id);

CREATE INDEX IF NOT EXISTS idx_p_category_category_id 
ON public."p_category"(category_id);

-- Composite index for product-category lookups
CREATE INDEX IF NOT EXISTS idx_p_category_composite 
ON public."p_category"(product_id, category_id);

-- Index for product images (used in catalogue responses)
CREATE INDEX IF NOT EXISTS idx_p_images_product_id 
ON public."p_images"(product_id);

-- Index for configurable product images
CREATE INDEX IF NOT EXISTS idx_cp_images_configurable_product_id 
ON public."cp_images"(configurable_product_id);

-- Index for configurable product members (used in catalogue queries)
CREATE INDEX IF NOT EXISTS idx_cp_members_cp_id 
ON public."cp_members"(cp_id);

CREATE INDEX IF NOT EXISTS idx_cp_members_product_id 
ON public."cp_members"(product_id);

-- Index for product stock (used in catalogue for available stock)
CREATE INDEX IF NOT EXISTS idx_p_stock_product_id 
ON public."p_stock"(product_id);

-- Index for distributor subscriptions (used in catalogue function)
CREATE INDEX IF NOT EXISTS idx_distributor_subscriptions_distributor_id 
ON public."distributor_subscriptions"(distributor_id);

-- Index for product attributes (used in catalogue responses)
CREATE INDEX IF NOT EXISTS idx_p_attribute_values_product_id 
ON public."p_attribute_values"(product_id);

CREATE INDEX IF NOT EXISTS idx_p_attribute_values_attribute_id 
ON public."p_attribute_values"(attribute_id);

-- Index for configurable product attributes
CREATE INDEX IF NOT EXISTS idx_cp_attributes_configurable_product_id 
ON public."cp_attributes"(configurable_product_id);

-- Index for products by external_id (used in lookups)
CREATE INDEX IF NOT EXISTS idx_products_external_id 
ON public."products"(external_id);

-- Index for configurable products by external_id
CREATE INDEX IF NOT EXISTS idx_configurable_products_external_id 
ON public."configurable_products"(external_id);

COMMENT ON INDEX idx_products_distributor_active IS 'Optimizes catalogue queries filtering by distributor and active status';
COMMENT ON INDEX idx_configurable_products_available IS 'Optimizes catalogue queries for configurable products';
COMMENT ON INDEX idx_p_category_composite IS 'Optimizes product-category relationship lookups';

