-- Seed data for product_categories junction table
-- Maps products to categories (1=Industrial Equipment, 2=Machinery, 3=Safety, 4=Construction, 5=Metalworking, 6=Electrical, 7=Hydraulics, 8=Tools)

INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 1, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-IND-001', 'PROD-IND-002', 'PROD-IND-003')
ON CONFLICT DO NOTHING;

INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 2, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-MACH-001', 'PROD-MACH-002', 'PROD-MACH-003')
ON CONFLICT DO NOTHING;

INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 3, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-SAFE-001', 'PROD-SAFE-002', 'PROD-SAFE-003')
ON CONFLICT DO NOTHING;

INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 4, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-CONST-001', 'PROD-CONST-002', 'PROD-CONST-003')
ON CONFLICT DO NOTHING;

INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 5, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-METAL-001', 'PROD-METAL-002', 'PROD-METAL-003')
ON CONFLICT DO NOTHING;

INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 6, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-ELEC-001', 'PROD-ELEC-002', 'PROD-ELEC-003')
ON CONFLICT DO NOTHING;

INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 7, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-HYD-001', 'PROD-HYD-002', 'PROD-HYD-003')
ON CONFLICT DO NOTHING;

INSERT INTO public.product_categories (product_id, category_id, created_date, last_modified, is_deleted)
SELECT p.id, 8, NOW(), NOW(), FALSE
FROM public.products p
WHERE p.external_id IN ('PROD-TOOL-001', 'PROD-TOOL-002', 'PROD-TOOL-003')
ON CONFLICT DO NOTHING;
