-- Fix supplier support_email to match user emails for proper supplier filtering
-- This ensures that supplier users can be properly identified and data is filtered correctly
-- This migration only runs if the suppliers table exists

DO $$
BEGIN
    -- Check if suppliers table exists before attempting updates
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'suppliers') THEN
        -- Update suppliers to match user emails
        -- Supplier 1: TechSupply Co. -> supplier1@b2b.local
        UPDATE suppliers 
        SET support_email = 'supplier1@b2b.local', last_modified = NOW()
        WHERE business_name = 'TechSupply Co.' AND support_email != 'supplier1@b2b.local';

        -- Supplier 2: Industrial Solutions Ltd -> supplier2@b2b.local
        UPDATE suppliers 
        SET support_email = 'supplier2@b2b.local', last_modified = NOW()
        WHERE business_name = 'Industrial Solutions Ltd' AND support_email != 'supplier2@b2b.local';

        -- Supplier 3: Global Equipment Inc -> supplier3@b2b.local
        UPDATE suppliers 
        SET support_email = 'supplier3@b2b.local', last_modified = NOW()
        WHERE business_name = 'Global Equipment Inc' AND support_email != 'supplier3@b2b.local';
    END IF;
END $$;

