-- Drop indexes
DROP INDEX IF EXISTS public.idx_price_history_changed_by_user_id;
DROP INDEX IF EXISTS public.idx_price_history_created_date;
DROP INDEX IF EXISTS public.idx_price_history_product_id;

-- Drop table
DROP TABLE IF EXISTS public.price_history;

