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


CREATE OR REPLACE FUNCTION public.get_all()
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
$$

-- user-roles ----------------------------------

    --writers
CREATE OR REPLACE FUNCTION public.add_role_to_user(
   user_identifier BIGINT,
   role_identifier BIGINT
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
   user_identifier BIGINT,
   role_identifier BIGINT
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
    user_identifier BIGINT
)
RETURNS TABLE(resource_id BIGINT)
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
