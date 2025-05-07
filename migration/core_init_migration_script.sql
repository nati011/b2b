-- roles 

-- resources

-- superadmin user

-- create payment partners
SELECT public.create_payment_partner(
    'chapa', 
    'icon', 
    'ACTIVE', 
    'https://api.chapa.co', 
    'CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5'
);

SELECT public.create_payment_partner(
    'payOnDelivery', 
    'icon', 
    'ACTIVE', 
    'https://no.baseurl', 
    'no-secret'
);