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

CREATE OR REPLACE FUNCTION public.update_action(
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

CREATE OR REPLACE FUNCTION public.update_name(
   r_name VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.resources (name)
    VALUES (r_name) 
    RETURNING id INTO new_id;
    RETURN new_id;
END;
$$;
-- readers
CREATE OR REPLACE FUNCTION public.get_resources_by_id(
   resource_id VARCHAR(255)
)
RETURNS SETOF public.resources
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM public.resources r
    WHERE r.id = resource_id
      AND t.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_resources_by_name(
   resource_name VARCHAR(255)
)
RETURNS SETOF public.resources
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM public.resources r
    WHERE r.name = resource_name
      AND t.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_resources()
RETURNS SETOF public.resources
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM public.resources r
    WHERE t.is_deleted = FALSE;
END;
$$;
-- Roles

-- Users