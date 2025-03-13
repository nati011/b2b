-- v 0.1

-- Resources ----------------------------------------

    -- writers
CREATE OR REPLACE FUNCTION public.create_resource(
   r_name VARCHAR(255),
   r_action VARCHAR(255)
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id BIGINT;
BEGIN
    INSERT INTO public.resources (name, action)
    VALUES (r_name, r_action) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_resource_name(
    resource_id BIGINT,
    new_name TEXT
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.resources
    SET name = new_name
    WHERE id = resource_id
      AND is_deleted = FALSE;

    RETURN resource_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_resource_action(
    resource_id BIGINT,
    new_action TEXT
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.resources
    SET action = new_action
    WHERE id = resource_id
      AND is_deleted = FALSE;

    RETURN resource_id;
END;
$$;


CREATE OR REPLACE FUNCTION public.delete_resource(
   r_id BIGINT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.resources
    SET is_deleted = TRUE
    WHERE id = r_id;
END;
$$;

    --readers

CREATE OR REPLACE FUNCTION public.get_resources_by_id(
    resource_id BIGINT
)
RETURNS TABLE(id BIGINT, action VARCHAR(255), name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.action, r.name
    FROM public.resources r
    WHERE r.id = resource_id
      AND r.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_resources_by_name(
    resource_name VARCHAR(255)
)
RETURNS TABLE(id BIGINT, action VARCHAR(255), name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.action, r.name
    FROM public.resources r
    WHERE r.name = resource_name
      AND r.is_deleted = FALSE
    LIMIT 1; 
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_resources()
RETURNS TABLE(id BIGINT, action VARCHAR(255), name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.action, r.name
    FROM public.resources r
    WHERE r.is_deleted = FALSE;
END;
$$;

-- Roles ----------------------------------------

    -- writers
CREATE OR REPLACE FUNCTION public.create_role(
   r_name VARCHAR(255),
   r_desc VARCHAR(255)
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id BIGINT;
BEGIN
    INSERT INTO public.roles (name, description)
    VALUES (r_name, r_desc) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_role_name(
    role_id BIGINT,
    new_name TEXT
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
BEGIN
    
    UPDATE public.roles
    SET name = new_name
    WHERE id = role_id
      AND is_deleted = FALSE;

    RETURN role_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_role_desc(
    role_id BIGINT,
    new_desc TEXT
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.roles
    SET description = new_desc
    WHERE id = role_id
      AND is_deleted = FALSE;

    RETURN role_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.delete_role(
   r_id BIGINT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.roles
    SET is_deleted = TRUE
    WHERE id = r_id;
END;
$$;

    -- readers
CREATE OR REPLACE FUNCTION public.get_roles_by_id(
    role_id BIGINT
)
RETURNS TABLE(id BIGINT, description VARCHAR(255), name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.description, r.name
    FROM public.roles r
    WHERE r.id = role_id
      AND r.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_roles_by_name(
    role_name VARCHAR(255)
)
RETURNS TABLE(id BIGINT, action VARCHAR(255), name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.description, r.name
    FROM public.roles r
    WHERE r.name = role_name
      AND r.is_deleted = FALSE
    LIMIT 1; 
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_roles()
RETURNS TABLE(id BIGINT, description VARCHAR(255), name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.description, r.name
    FROM public.roles r
    WHERE r.is_deleted = FALSE;
END;
$$;

-- role_resources ----------------------------------------

    -- writers
CREATE OR REPLACE FUNCTION public.add_resource_to_role(
   role_identifier BIGINT,
   resource_identifier BIGINT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.role_resources (role_id, resource_id)
    VALUES (role_identifier, resource_identifier);
END;
$$;

CREATE OR REPLACE FUNCTION public.remove_resource_from_role(
   role_identifier BIGINT,
   resource_identifier BIGINT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
   UPDATE public.role_resources
   SET is_deleted = TRUE
   WHERE role_id = role_identifier 
       AND resource_id = resource_identifier
       AND is_deleted = FALSE;
END;
$$;

    -- readers
CREATE OR REPLACE FUNCTION public.get_all_resource_by_role(
    role_identifier INT
)
RETURNS TABLE(resource_id BIGINT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.resource_id
    FROM public.role_resources r
    WHERE r.role_id = role_identifier 
        AND r.is_deleted = FALSE;
END;
$$;

-- Users ----------------------------------------
    
    
    
    --readers
CREATE OR REPLACE FUNCTION public.get_users_by_id(
    user_id INT
)
RETURNS TABLE(id INT, fullname VARCHAR(255), email VARCHAR(255), phone VARCHAR(255), username VARCHAR(255), birthdate date, is_active boolean, external_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.fullname, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
    FROM public.users u
    WHERE u.id = user_id
      AND u.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_users_by_email(
    user_email VARCHAR(255)
)
RETURNS TABLE(id INT, fullname VARCHAR(255), email VARCHAR(255), phone VARCHAR(255), username VARCHAR(255), birthdate date, is_active boolean, external_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.fullname, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
    FROM public.users u
    WHERE u.email = user_email
      AND u.is_deleted = FALSE
    LIMIT 1;
END;
$$;


CREATE OR REPLACE FUNCTION public.get_users_by_phone(
    user_phone VARCHAR(255)
)
RETURNS TABLE(id INT, fullname VARCHAR(255), email VARCHAR(255), phone VARCHAR(255), username VARCHAR(255), birthdate date, is_active boolean, external_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.fullname, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
    FROM public.users u
    WHERE u.phone_number = user_phone
      AND u.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_users_by_username(
    user_username VARCHAR(255)
)
RETURNS TABLE(id INT, fullname VARCHAR(255), email VARCHAR(255), phone VARCHAR(255), username VARCHAR(255), birthdate date, is_active boolean, external_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.fullname, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
    FROM public.users u
    WHERE u.username = user_username
      AND u.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_users_by_active_status(
    user_active_status BOOLEAN
)
RETURNS TABLE(id INT, fullname VARCHAR(255), email VARCHAR(255), phone VARCHAR(255), username VARCHAR(255), birthdate date, is_active boolean, external_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.fullname, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
    FROM public.users u
    WHERE u.is_active = user_active_status
      AND u.is_deleted = FALSE
    LIMIT 1;
END;
$$;


CREATE OR REPLACE FUNCTION public.get_users_by_external_id(
    user_external_id VARCHAR(255)
)
RETURNS TABLE(id INT, fullname VARCHAR(255), email VARCHAR(255), phone VARCHAR(255), username VARCHAR(255), birthdate date, is_active boolean, external_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.fullname, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
    FROM public.users u
    WHERE u.external_id = user_external_id
      AND u.is_deleted = FALSE
    LIMIT 1;
END;
$$;


CREATE OR REPLACE FUNCTION public.get_all_users()
RETURNS TABLE(id INT, fullname VARCHAR(255), email VARCHAR(255), phone VARCHAR(255), username VARCHAR(255), birthdate date, is_active boolean, external_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.fullname, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
    FROM public.users u
    WHERE u.is_deleted = FALSE
    LIMIT 1;
END;
$$;

--writers

CREATE OR REPLACE FUNCTION public.create_user(
   u_fullname VARCHAR(255),
   u_email VARCHAR(255),
   u_phone VARCHAR(255),
   u_username VARCHAR(255),
   u_dob DATE,
   u_is_active BOOLEAN,
   u_external_id VARCHAR(255))
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.users 
	(fullname, 
	email, 
	phone_number, 
	username, 
	birth_date, 
	is_active, 
	external_id)
    VALUES 	
	(u_fullname, 
	u_email, 
	u_phone, 
	u_username, 
	u_dob, 
	u_is_active, 
	u_external_id)

    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.delete_user(
   u_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.users
    SET is_deleted = TRUE
    WHERE id = u_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_user_fullname(
    user_id INT,
    new_fullname VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.users
    SET fullname = new_fullname
    WHERE id = user_id
      AND is_deleted = FALSE;

    RETURN user_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_user_email(
    user_id INT,
    new_email VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.users
    SET email = new_email
    WHERE id = user_id
      AND is_deleted = FALSE;

    RETURN user_id;
END;
$$;


CREATE OR REPLACE FUNCTION public.update_user_dob(
    user_id INT,
    new_dob DATE
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.users
    SET birth_date = new_dob
    WHERE id = user_id
      AND is_deleted = FALSE;

    RETURN user_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_user_is_active_status(
    user_id INT,
    new_active_status BOOLEAN
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.users
    SET is_active = new_active_status
    WHERE id = user_id
      AND is_deleted = FALSE;

    RETURN user_id;
END;
$$;


CREATE OR REPLACE FUNCTION public.update_user_phone(
    user_id INT,
    new_phone VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.users
    SET phone_number = new_phone
    WHERE id = user_id
      AND is_deleted = FALSE;

    RETURN user_id;
END;
$$;


CREATE OR REPLACE FUNCTION public.update_user_name(
    user_id INT,
    new_username VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.users
    SET username = new_username
    WHERE id = user_id
      AND is_deleted = FALSE;

    RETURN user_id;
END;
$$;

-- user-roles ----------------------------------

    --writers
CREATE OR REPLACE FUNCTION public.add_role_to_user(
   user_identifier INT,
   role_identifier INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.user_roles (user_id, role_id)
    VALUES (role_identifier, user_identifier);
END;
$$;

CREATE OR REPLACE FUNCTION public.remove_role_from_user(
   user_identifier INT,
   role_identifier INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
   UPDATE public.user_roles
   SET is_deleted = TRUE
   WHERE role_id = role_identifier 
       AND user_id = user_identifier
       AND is_deleted = FALSE;
END;
$$;

    -- readers
CREATE OR REPLACE FUNCTION public.get_all_role_by_user(
    user_identifier INT
)
RETURNS TABLE(role_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.role_id
    FROM public.user_roles r
    WHERE r.user_id = user_identifier 
        AND r.is_deleted = FALSE;
END;
$$;

-- invoices ----------------------------------

    -- writers
CREATE OR REPLACE FUNCTION public.create_invoice(
   i_status VARCHAR(255),
   i_externalId VARCHAR(255),
   i_orderId int,
   i_subtotal DECIMAL(12,2),
   i_taxAmount DECIMAL(12,2)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id BIGINT;
BEGIN
    INSERT INTO public.invoices (status, external_id, order_id, subtotal, tax_amount)
    VALUES (i_status, i_externalId, i_orderId, i_subtotal, i_taxAmount) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;
    -- readers
CREATE OR REPLACE FUNCTION public.get_invoices_by_id(
    i_invoice_id INT
)
RETURNS TABLE(id INT, status VARCHAR(255), external_id VARCHAR(255), order_id INT, subtotal DECIMAL(12,2), tax_amount DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, i.status, i.external_id, i.order_id, i.subtotal, i.tax_amount
    FROM public.invoices i
    WHERE i.id = i_invoice_id
      AND i.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_invoices()
RETURNS TABLE(id INT, status VARCHAR(255), external_id VARCHAR(255), order_id INT, subtotal DECIMAL(12,2), tax_amount DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, i.status, i.external_id, i.order_id, i.subtotal, i.tax_amount
    FROM public.invoices i
    WHERE i.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_invoices_by_external_id(
     i_external_id VARCHAR(255)
)
RETURNS TABLE(id INT, status VARCHAR(255), external_id VARCHAR(255), order_id INT, subtotal DECIMAL(12,2), tax_amount DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, i.status, i.external_id, i.order_id, i.subtotal, i.tax_amount
    FROM public.invoices i
    WHERE i.is_deleted = FALSE 
    AND i.external_id = i_external_id;
END;
$$;


CREATE OR REPLACE FUNCTION public.get_invoices_by_status(
     i_status VARCHAR(255)
)
RETURNS TABLE(id INT, status VARCHAR(255), external_id VARCHAR(255), order_id INT, subtotal DECIMAL(12,2), tax_amount DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, i.status, i.external_id, i.order_id, i.subtotal, i.tax_amount
    FROM public.invoices i
    WHERE i.is_deleted = FALSE 
    AND i.status = i_status;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_invoices_by_order_id(
     i_order_id INT
)
RETURNS TABLE(id INT, status VARCHAR(255), external_id VARCHAR(255), order_id INT, subtotal DECIMAL(12,2), tax_amount DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, i.status, i.external_id, i.order_id, i.subtotal, i.tax_amount
    FROM public.invoices i
    WHERE i.is_deleted = FALSE 
    AND i.order_id = i_order_id
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_invoice_externalId(
    i_invoice_id INT,
    new_external_id VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.invoices
    SET external_id = new_external_id
    WHERE id = i_invoice_id
      AND is_deleted = FALSE;

    RETURN i_invoice_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_invoice_status(
    i_invoice_id INT,
    new_status VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.invoices
    SET status = new_status
    WHERE id = i_invoice_id
      AND is_deleted = FALSE;

    RETURN i_invoice_id;
END;
$$;


-- invoice line items--------------------------------------------

    -- writers
CREATE OR REPLACE FUNCTION public.create_invoice_line_item(
   i_product_name VARCHAR(255),
   i_qty INT,
   i_price DECIMAL(12,2),
   i_product_id INT,
   i_invoice_id INT
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.invoice_line_items (product_name, qty, price, product_id, invoice_id)
    VALUES (i_product_name, i_qty, i_price, i_product_id, i_invoice_id) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;
    
    -- readers
CREATE OR REPLACE FUNCTION public.get_invoice_line_item_by_invoice_id(
    i_invoice_id INT
)
RETURNS TABLE(id INT, product_name VARCHAR(255), qty INT, price DECIMAL(12,2), product_id INT, invoice_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, i.product_name, i.qty, i.price, i.product_id, i.invoice_id
    FROM public.invoice_line_items i
    WHERE i.invoice_id = i_invoice_id
      AND i.is_deleted = FALSE
    LIMIT 1;
END;
$$;

-- product -------------------------------------------------------------
    
    -- writers
CREATE OR REPLACE FUNCTION public.create_product(
  p_product_name VARCHAR(255),
  p_product_description VARCHAR(255),
  p_external_id VARCHAR(255),
  p_distributor_id INT
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.products (name, description, external_id, distributor_id)
    VALUES (p_product_name, p_product_description, p_external_id, p_distributor_id) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;


CREATE OR REPLACE FUNCTION public.update_product_name(
    i_product_id INT,
    new_name VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.products
    SET name = new_name
    WHERE id = i_product_id
      AND is_deleted = FALSE;

    RETURN i_product_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_product_external_Id(
    i_product_id INT,
    new_external_Id VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.products
    SET external_id = new_external_Id
    WHERE id = i_product_id
      AND is_deleted = FALSE;

    RETURN i_product_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_product_desc(
    i_product_id INT,
    new_desc VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.products
    SET description = new_desc
    WHERE id = i_product_id
      AND is_deleted = FALSE;

    RETURN i_product_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_product_status(
    i_product_id INT,
    new_is_active BOOLEAN
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.products
    SET is_active = new_is_active
    WHERE id = i_product_id
      AND is_deleted = FALSE;

    RETURN i_product_id;
END;
$$;


    -- readers
CREATE OR REPLACE FUNCTION public.get_products_by_id(
    p_product_id INT
)
RETURNS TABLE(id INT, product_name VARCHAR(255), product_description TEXT, external_id VARCHAR(255), is_active BOOLEAN, distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.description, p.external_id, p.is_active, p.distributor_id
    FROM public.products p
    WHERE p.id = p_product_id
      AND p.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_name(
    p_product_name VARCHAR(255)
)
RETURNS TABLE(id INT, product_name VARCHAR(255), product_description TEXT, external_id VARCHAR(255), is_active BOOLEAN, distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.description, p.external_id, p.is_active, p.distributor_id
    FROM public.products p
    WHERE p.name = p_product_name
      AND p.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_externalId(
    p_external_id VARCHAR(255)
)
RETURNS TABLE(id INT, product_name VARCHAR(255), product_description TEXT, external_id VARCHAR(255), is_active BOOLEAN, distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.description, p.external_id, p.is_active, p.distributor_id
    FROM public.products p
    WHERE p.external_id = p_external_id
      AND p.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_distributorId(
    p_distributor_id INT
)
RETURNS TABLE(id INT, product_name VARCHAR(255), product_description TEXT, external_id VARCHAR(255), is_active BOOLEAN, distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.description, p.external_id, p.is_active, p.distributor_id
    FROM public.products p
    WHERE p.distributor_id = p_distributor_id
      AND p.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_products()
RETURNS TABLE(id INT, product_name VARCHAR(255), product_description TEXT, external_id VARCHAR(255), is_active BOOLEAN, distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.description, p.external_id, p.is_active, p.distributor_id
    FROM public.products p
    WHERE p.is_deleted = FALSE;
END;
$$;

-- product_images ---------------------------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.add_image_to_product(
  i_url VARCHAR(255),
  i_blur_hash VARCHAR(255),
  i_product_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.p_images (url, blur_hash, product_id)
    VALUES (i_url, i_blur_hash, i_product_id);
END;
$$;

CREATE OR REPLACE FUNCTION public.remove_all_product_images(
    i_product_id INT
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.p_images
    SET is_deleted = TRUE
    WHERE product_id = i_product_id;

    RETURN i_product_id;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_images_by_productId(
    i_product_id INT
)
RETURNS TABLE(image_url VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT url
    FROM public.p_images i
    WHERE i.product_id = i_product_id
      AND i.is_deleted = FALSE;
END;
$$;

-- product_categories ---------------------------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.add_category_to_product(
  i_product_id INT,
  i_category_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.p_category (product_id, category_id)
    VALUES (i_product_id, i_category_id);
END;
$$;

CREATE OR REPLACE FUNCTION public.remove_all_categories_from_product(
    i_product_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.p_category
    SET is_deleted = TRUE
    WHERE product_id = i_product_id;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_categories_by_productId(
    p_product_id BIGINT
)
RETURNS TABLE(category_id BIGINT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.category_id
    FROM public.p_category p
    WHERE p.product_id = p_product_id
      AND p.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_categoryId(
    p_category_id BIGINT
)
RETURNS TABLE(product_id BIGINT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.product_id
    FROM public.p_category p
    WHERE p.category_id = p_category_id
      AND p.is_deleted = FALSE;
END;
$$;


-- product_price ---------------------------------------------------

    -- writer
CREATE OR REPLACE FUNCTION public.add_price_to_product(
  i_product_id INT,
  i_price INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.p_prices (price, product_id)
    VALUES (i_price, i_product_id);
END;
$$;

CREATE OR REPLACE FUNCTION public.update_product_price(
    i_product_id INT,
    i_product_price DECIMAL(12, 2)
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.p_prices
    SET price = i_product_price
    WHERE product_id = i_product_id
        AND is_deleted = FALSE;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_price_by_productId(
    p_product_id INT
)
RETURNS DECIMAL(12, 2)
LANGUAGE plpgsql
AS $$
DECLARE
    price DECIMAL(12, 2);
BEGIN
    SELECT p.price INTO price
    FROM public.p_prices p
    WHERE p.product_id = p_product_id
      AND p.is_deleted = FALSE
    LIMIT 1;

    RETURN COALESCE(price, 0);
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_price_range(
    p_min DECIMAL(12, 2),
    p_max DECIMAL(12, 2)
)
RETURNS TABLE(product_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.product_id
    FROM public.p_prices p
    WHERE p.price BETWEEN p_min AND p_max
      AND p.is_deleted = FALSE;
END;
$$;

-- product stock ----------------------------------------------------

    -- writer
CREATE OR REPLACE FUNCTION public.create_product_stock(
  i_quantity INT,
  i_product_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.p_stock (quantity, product_id)
    VALUES (i_quantity, i_product_id);
END;
$$;

CREATE OR REPLACE FUNCTION public.update_product_stock(
    i_product_id INT,
    i_quantity INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.p_stock
    SET quantity = i_quantity
    WHERE product_id = i_product_id
        AND is_deleted = FALSE;
END;
$$;
    
    -- reader
CREATE OR REPLACE FUNCTION public.get_stock_by_productId(
    p_product_id INT
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    quantity INT;
BEGIN
    SELECT p.quantity INTO quantity
    FROM public.p_stock p
    WHERE p.product_id = p_product_id
      AND p.is_deleted = FALSE
    LIMIT 1;

    RETURN COALESCE(quantity, 0);
END;
$$;

-- stock ledger ------------------------------------------------------

    -- writer
CREATE OR REPLACE FUNCTION public.stock_operation_ledger_entry(
  s_quantity INT,
  s_product_id INT,
  s_stock_operation VARCHAR(255),
  s_created_by_user_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.s_ledger (quantity, product_id, operation, created_by_user_id)
    VALUES (s_quantity, s_product_id, s_stock_operation, s_created_by_user_id);
END;
$$;

    -- reader

-- product attributes ------------------------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.create_product_attribute(
  i_name VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.p_attributes (name)
    VALUES (i_name) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;
    -- reader

-- product attribute-values ------------------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.create_product_attribute_value(
  i_name VARCHAR(255),
  i_product_id INT,
  i_attribute_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.p_attribute_values (name, product_id, attribute_id)
    VALUES (i_name, i_product_id, i_attribute_id);
END;
$$;
    -- reader