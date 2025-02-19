-- v 0.1

-- Resources

    -- writers
CREATE OR REPLACE FUNCTION public.create_resource(
   r_name VARCHAR(255),
   r_action VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.resources (name, action)
    VALUES (r_name, r_action) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_resource_name(
    resource_id INT,
    new_name TEXT
)
RETURNS INT
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
    resource_id INT,
    new_action TEXT
)
RETURNS INT
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
   r_id INT
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

CREATE OR REPLACE FUNCTION public.get_resources_by_id(
    resource_id INT
)
RETURNS TABLE(id INT, action VARCHAR(255), name VARCHAR(255))
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
RETURNS TABLE(id INT, action VARCHAR(255), name VARCHAR(255))
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
RETURNS TABLE(id INT, action VARCHAR(255), name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.action, r.name
    FROM public.resources r
    WHERE r.is_deleted = FALSE;
END;
$$;

-- Roles

    -- writers
CREATE OR REPLACE FUNCTION public.create_role(
   r_name VARCHAR(255),
   r_desc VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.roles (name, description)
    VALUES (r_name, r_desc) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_role_name(
    role_id INT,
    new_name TEXT
)
RETURNS INT
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
    role_id INT,
    new_desc TEXT
)
RETURNS INT
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
   r_id INT
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

CREATE OR REPLACE FUNCTION public.get_roles_by_id(
    role_id INT
)
RETURNS TABLE(id INT, description VARCHAR(255), name VARCHAR(255))
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
RETURNS TABLE(id INT, action VARCHAR(255), name VARCHAR(255))
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
RETURNS TABLE(id INT, description VARCHAR(255), name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.description, r.name
    FROM public.roles r
    WHERE r.is_deleted = FALSE;
END;
$$;


-- Users