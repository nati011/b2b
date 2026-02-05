-- Populate supplier_id for existing orders
-- Assigns 4 orders to supplier_id = 1 and the rest to supplier_id = 2

-- Update first 4 orders to supplier_id = 1
UPDATE public.orders
SET supplier_id = 1,
    last_modified = NOW()
WHERE id IN (
    SELECT id 
    FROM public.orders 
    WHERE is_deleted = FALSE 
      AND supplier_id IS NULL
    ORDER BY created_date ASC
    LIMIT 4
);

-- Update remaining orders to supplier_id = 2
UPDATE public.orders
SET supplier_id = 2,
    last_modified = NOW()
WHERE is_deleted = FALSE 
  AND supplier_id IS NULL;

-- Verify the update
DO $$
DECLARE
    supplier1_count INT;
    supplier2_count INT;
    null_count INT;
BEGIN
    SELECT COUNT(*) INTO supplier1_count 
    FROM public.orders 
    WHERE is_deleted = FALSE AND supplier_id = 1;
    
    SELECT COUNT(*) INTO supplier2_count 
    FROM public.orders 
    WHERE is_deleted = FALSE AND supplier_id = 2;
    
    SELECT COUNT(*) INTO null_count 
    FROM public.orders 
    WHERE is_deleted = FALSE AND supplier_id IS NULL;
    
    RAISE NOTICE 'Orders with supplier_id = 1: %', supplier1_count;
    RAISE NOTICE 'Orders with supplier_id = 2: %', supplier2_count;
    RAISE NOTICE 'Orders with NULL supplier_id: %', null_count;
END $$;




