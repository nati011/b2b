-- Seed data for product_categories junction table
-- Development seed data
-- Maps products to categories

-- Electronics products (category_id = 1)
INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 1, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-ELEC-001', 'PROD-ELEC-002', 'PROD-ELEC-003')
ON CONFLICT DO NOTHING;

-- Office Supplies products (category_id = 2)
INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 2, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-OFF-001', 'PROD-OFF-002', 'PROD-OFF-003')
ON CONFLICT DO NOTHING;

-- Tools & Equipment products (category_id = 3)
INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 3, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-TOOL-001', 'PROD-TOOL-002', 'PROD-TOOL-003')
ON CONFLICT DO NOTHING;

-- Furniture products (category_id = 4)
INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 4, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-FURN-001', 'PROD-FURN-002', 'PROD-FURN-003')
ON CONFLICT DO NOTHING;

-- Safety & Security products (category_id = 5)
INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 5, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-SAFE-001', 'PROD-SAFE-002', 'PROD-SAFE-003')
ON CONFLICT DO NOTHING;

-- Industrial products (category_id = 6)
INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 6, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-IND-001', 'PROD-IND-002', 'PROD-IND-003')
ON CONFLICT DO NOTHING;

-- Construction Materials products (category_id = 7)
INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 7, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-CONST-001', 'PROD-CONST-002', 'PROD-CONST-003')
ON CONFLICT DO NOTHING;

-- Cleaning Supplies products (category_id = 8)
INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 8, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-CLEAN-001', 'PROD-CLEAN-002', 'PROD-CLEAN-003')
ON CONFLICT DO NOTHING;

