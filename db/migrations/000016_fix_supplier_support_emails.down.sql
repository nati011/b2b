-- Revert supplier support_email changes
-- Note: This reverts to the original seed values, but may not match actual user emails

UPDATE suppliers 
SET support_email = 'support@techsupply.local', last_modified = NOW()
WHERE business_name = 'TechSupply Co.';

UPDATE suppliers 
SET support_email = 'info@industrialsolutions.local', last_modified = NOW()
WHERE business_name = 'Industrial Solutions Ltd';

UPDATE suppliers 
SET support_email = 'contact@globalequip.local', last_modified = NOW()
WHERE business_name = 'Global Equipment Inc';









