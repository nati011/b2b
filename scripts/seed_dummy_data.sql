-- ============================================
-- Seed Dummy Data for Testing
-- ============================================
-- This script creates:
-- - Categories
-- - Test User and Distributor
-- - Active Subscription
-- - Simple Products (with images, categories, stock)
-- - Configurable Products (with images and linked products)
-- ============================================

BEGIN;

-- ============================================
-- 1. Create Categories
-- ============================================
INSERT INTO public.category (name, is_deleted) VALUES
    ('Electronics', false),
    ('Clothing', false),
    ('Home & Kitchen', false),
    ('Sports & Outdoors', false),
    ('Beauty & Personal Care', false),
    ('Books', false),
    ('Toys & Games', false),
    ('Automotive', false)
ON CONFLICT DO NOTHING;

-- ============================================
-- Main DO Block for all operations
-- ============================================
DO $$
DECLARE
    -- Category IDs
    cat_electronics_id INT;
    cat_clothing_id INT;
    cat_home_id INT;
    cat_sports_id INT;
    cat_beauty_id INT;
    cat_books_id INT;
    
    -- User and Distributor IDs
    user_id_var INT;
    distributor_id_var INT;
    subscription_plan_id_var INT;
    subscription_id_var INT;
    
    -- Product IDs
    product1_id INT;
    product2_id INT;
    product3_id INT;
    product4_id INT;
    product5_id INT;
    product6_id INT;
    
    -- Configurable Product IDs
    cp1_id INT;
    cp2_id INT;
    
    -- Variant Product IDs
    cp1_var1_id INT;
    cp1_var2_id INT;
    cp1_var3_id INT;
    cp2_var1_id INT;
    cp2_var2_id INT;
BEGIN
    -- Get category IDs
    SELECT id INTO cat_electronics_id FROM public.category WHERE name = 'Electronics' AND is_deleted = false LIMIT 1;
    SELECT id INTO cat_clothing_id FROM public.category WHERE name = 'Clothing' AND is_deleted = false LIMIT 1;
    SELECT id INTO cat_home_id FROM public.category WHERE name = 'Home & Kitchen' AND is_deleted = false LIMIT 1;
    SELECT id INTO cat_sports_id FROM public.category WHERE name = 'Sports & Outdoors' AND is_deleted = false LIMIT 1;
    SELECT id INTO cat_beauty_id FROM public.category WHERE name = 'Beauty & Personal Care' AND is_deleted = false LIMIT 1;
    SELECT id INTO cat_books_id FROM public.category WHERE name = 'Books' AND is_deleted = false LIMIT 1;

    -- ============================================
    -- 2. Create Test User
    -- ============================================
    SELECT id INTO user_id_var FROM public.users WHERE username = 'testdistributor' AND is_deleted = false LIMIT 1;
    
    IF user_id_var IS NULL THEN
        INSERT INTO public.users (firstName, lastName, email, phone_number, username, birth_date, is_active, is_deleted)
        VALUES ('Test', 'Distributor', 'test.distributor@efoyeta.com', '+251911234567', 'testdistributor', '1990-01-01', true, false)
        RETURNING id INTO user_id_var;
    END IF;

    -- ============================================
    -- 3. Create Distributor
    -- ============================================
    SELECT id INTO distributor_id_var FROM public.distributors WHERE is_deleted = false ORDER BY id DESC LIMIT 1;
    
    IF distributor_id_var IS NULL THEN
        INSERT INTO public.distributors (is_active, is_deleted)
        VALUES (true, false)
        RETURNING id INTO distributor_id_var;
    END IF;

    -- ============================================
    -- 4. Link User to Distributor
    -- ============================================
    INSERT INTO public.distributor_users (user_id, distributor_id, is_deleted)
    VALUES (user_id_var, distributor_id_var, false)
    ON CONFLICT DO NOTHING;

    -- ============================================
    -- 5. Create Distributor Business Info
    -- ============================================
    INSERT INTO public.distributor_business_info (name, tin, distributor_id, is_deleted)
    VALUES ('Efoyeta Test Distributor', 'TIN123456', distributor_id_var, false)
    ON CONFLICT DO NOTHING;

    -- ============================================
    -- 6. Create Active Subscription
    -- ============================================
    SELECT id INTO subscription_plan_id_var FROM public.subscription_plans WHERE name = 'Year Plan' AND is_deleted = false LIMIT 1;
    
    IF subscription_plan_id_var IS NULL THEN
        INSERT INTO public.subscription_plans (name, price, term_in_month, description, is_deleted)
        VALUES ('Year Plan', 11500.00, 12, 'Annual subscription plan with best value', false)
        RETURNING id INTO subscription_plan_id_var;
    END IF;

    INSERT INTO public.distributor_subscriptions (subscription_plan_id, distributor_id, payment_partner_Id, status, is_deleted)
    VALUES (subscription_plan_id_var, distributor_id_var, 1, 'active', false)
    ON CONFLICT DO NOTHING;

    -- ============================================
    -- 7. Create Simple Products
    -- ============================================
    
    -- Product 1: Wireless Headphones
    SELECT * FROM public.create_product(
        'Wireless Bluetooth Headphones',
        'Premium wireless headphones with noise cancellation and 30-hour battery life',
        'PROD-001',
        distributor_id_var,
        2500.00
    ) INTO product1_id;

    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', product1_id);
    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1484704849700-f032a568e944?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', product1_id);
    PERFORM public.add_category_to_product(product1_id, cat_electronics_id);
    PERFORM public.update_product_stock(product1_id, 150);
    UPDATE public.products SET is_active = true WHERE id = product1_id;

    -- Product 2: Cotton T-Shirt
    SELECT * FROM public.create_product(
        'Premium Cotton T-Shirt',
        '100% organic cotton t-shirt, comfortable and breathable',
        'PROD-002',
        distributor_id_var,
        450.00
    ) INTO product2_id;

    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', product2_id);
    PERFORM public.add_category_to_product(product2_id, cat_clothing_id);
    PERFORM public.update_product_stock(product2_id, 300);
    UPDATE public.products SET is_active = true WHERE id = product2_id;

    -- Product 3: Coffee Maker
    SELECT * FROM public.create_product(
        'Automatic Coffee Maker',
        'Programmable coffee maker with thermal carafe, makes 12 cups',
        'PROD-003',
        distributor_id_var,
        3500.00
    ) INTO product3_id;

    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1517487881594-2787fef5ebf7?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', product3_id);
    PERFORM public.add_category_to_product(product3_id, cat_home_id);
    PERFORM public.update_product_stock(product3_id, 50);
    UPDATE public.products SET is_active = true WHERE id = product3_id;

    -- Product 4: Running Shoes
    SELECT * FROM public.create_product(
        'Professional Running Shoes',
        'Lightweight running shoes with cushioned sole for maximum comfort',
        'PROD-004',
        distributor_id_var,
        3200.00
    ) INTO product4_id;

    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', product4_id);
    PERFORM public.add_category_to_product(product4_id, cat_sports_id);
    PERFORM public.update_product_stock(product4_id, 120);
    UPDATE public.products SET is_active = true WHERE id = product4_id;

    -- Product 5: Face Moisturizer
    SELECT * FROM public.create_product(
        'Hydrating Face Moisturizer',
        'Daily face moisturizer with SPF 30, suitable for all skin types',
        'PROD-005',
        distributor_id_var,
        850.00
    ) INTO product5_id;

    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1556228578-0d85b1a4d571?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', product5_id);
    PERFORM public.add_category_to_product(product5_id, cat_beauty_id);
    PERFORM public.update_product_stock(product5_id, 200);
    UPDATE public.products SET is_active = true WHERE id = product5_id;

    -- Product 6: Programming Book
    SELECT * FROM public.create_product(
        'Complete Guide to Web Development',
        'Comprehensive guide covering HTML, CSS, JavaScript, and modern frameworks',
        'PROD-006',
        distributor_id_var,
        1200.00
    ) INTO product6_id;

    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', product6_id);
    PERFORM public.add_category_to_product(product6_id, cat_books_id);
    PERFORM public.update_product_stock(product6_id, 80);
    UPDATE public.products SET is_active = true WHERE id = product6_id;

    -- ============================================
    -- 8. Create Configurable Products
    -- ============================================
    
    -- Configurable Product 1: Smartphone (with different storage options)
    SELECT * FROM public.create_configurable_product(
        'Smartphone Pro',
        'Latest smartphone with multiple storage and color options',
        'CPROD-001'
    ) INTO cp1_id;

    PERFORM public.add_image_to_configurable_product('https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', cp1_id);
    PERFORM public.add_image_to_configurable_product('https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', cp1_id);

    -- Create variant products for the configurable product
    -- Variant 1: 64GB
    SELECT * FROM public.create_product(
        'Smartphone Pro 64GB',
        'Smartphone Pro with 64GB storage',
        'CPROD-001-V1',
        distributor_id_var,
        8500.00
    ) INTO cp1_var1_id;
    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', cp1_var1_id);
    PERFORM public.add_category_to_product(cp1_var1_id, cat_electronics_id);
    PERFORM public.update_product_stock(cp1_var1_id, 25);
    UPDATE public.products SET is_active = true WHERE id = cp1_var1_id;
    PERFORM public.add_product_to_configurable_product(cp1_id, cp1_var1_id);

    -- Variant 2: 128GB
    SELECT * FROM public.create_product(
        'Smartphone Pro 128GB',
        'Smartphone Pro with 128GB storage',
        'CPROD-001-V2',
        distributor_id_var,
        9500.00
    ) INTO cp1_var2_id;
    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', cp1_var2_id);
    PERFORM public.add_category_to_product(cp1_var2_id, cat_electronics_id);
    PERFORM public.update_product_stock(cp1_var2_id, 30);
    UPDATE public.products SET is_active = true WHERE id = cp1_var2_id;
    PERFORM public.add_product_to_configurable_product(cp1_id, cp1_var2_id);

    -- Variant 3: 256GB
    SELECT * FROM public.create_product(
        'Smartphone Pro 256GB',
        'Smartphone Pro with 256GB storage',
        'CPROD-001-V3',
        distributor_id_var,
        11000.00
    ) INTO cp1_var3_id;
    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', cp1_var3_id);
    PERFORM public.add_category_to_product(cp1_var3_id, cat_electronics_id);
    PERFORM public.update_product_stock(cp1_var3_id, 20);
    UPDATE public.products SET is_active = true WHERE id = cp1_var3_id;
    PERFORM public.add_product_to_configurable_product(cp1_id, cp1_var3_id);

    -- Activate configurable product
    UPDATE public.configurable_products SET is_available = true WHERE id = cp1_id;

    -- Configurable Product 2: Laptop (with different RAM options)
    SELECT * FROM public.create_configurable_product(
        'Gaming Laptop',
        'High-performance gaming laptop with multiple RAM and storage configurations',
        'CPROD-002'
    ) INTO cp2_id;

    PERFORM public.add_image_to_configurable_product('https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', cp2_id);

    -- Variant 1: 8GB RAM
    SELECT * FROM public.create_product(
        'Gaming Laptop 8GB RAM',
        'Gaming laptop with 8GB RAM and 256GB SSD',
        'CPROD-002-V1',
        distributor_id_var,
        45000.00
    ) INTO cp2_var1_id;
    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', cp2_var1_id);
    PERFORM public.add_category_to_product(cp2_var1_id, cat_electronics_id);
    PERFORM public.update_product_stock(cp2_var1_id, 15);
    UPDATE public.products SET is_active = true WHERE id = cp2_var1_id;
    PERFORM public.add_product_to_configurable_product(cp2_id, cp2_var1_id);

    -- Variant 2: 16GB RAM
    SELECT * FROM public.create_product(
        'Gaming Laptop 16GB RAM',
        'Gaming laptop with 16GB RAM and 512GB SSD',
        'CPROD-002-V2',
        distributor_id_var,
        55000.00
    ) INTO cp2_var2_id;
    PERFORM public.add_image_to_product('https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=800', 'LGF5]+Yk^6#M@-5c,1J5@[or[Q6.', cp2_var2_id);
    PERFORM public.add_category_to_product(cp2_var2_id, cat_electronics_id);
    PERFORM public.update_product_stock(cp2_var2_id, 10);
    UPDATE public.products SET is_active = true WHERE id = cp2_var2_id;
    PERFORM public.add_product_to_configurable_product(cp2_id, cp2_var2_id);

    UPDATE public.configurable_products SET is_available = true WHERE id = cp2_id;

    RAISE NOTICE 'Seed data created successfully!';
    RAISE NOTICE 'Created % categories', (SELECT COUNT(*) FROM public.category WHERE is_deleted = false);
    RAISE NOTICE 'Created % active products', (SELECT COUNT(*) FROM public.products WHERE is_deleted = false AND is_active = true);
    RAISE NOTICE 'Created % configurable products', (SELECT COUNT(*) FROM public.configurable_products WHERE is_deleted = false AND is_available = true);
END $$;

COMMIT;

-- ============================================
-- Verification Queries (uncomment to run)
-- ============================================
-- SELECT COUNT(*) as total_categories FROM public.category WHERE is_deleted = false;
-- SELECT COUNT(*) as total_products FROM public.products WHERE is_deleted = false AND is_active = true;
-- SELECT COUNT(*) as total_configurable_products FROM public.configurable_products WHERE is_deleted = false AND is_available = true;
-- SELECT d.id, d.is_active, ds.status FROM public.distributors d 
-- JOIN public.distributor_subscriptions ds ON d.id = ds.distributor_id 
-- WHERE d.is_deleted = false AND ds.status = 'active';
