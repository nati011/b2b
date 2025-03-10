-- v 0.1
-- Resources ----------------------------------------
-- writers
create or replace function public.create_resource (r_name VARCHAR(255), r_action VARCHAR(255)) RETURNS BIGINT LANGUAGE plpgsql as $$
    DECLARE
        new_id BIGINT;
    BEGIN
        INSERT INTO public.resources (name, action)
        VALUES (r_name, r_action) 
        RETURNING id INTO new_id;

        RETURN new_id;
    END;
    $$;

create or replace function public.update_resource_name (resource_id BIGINT, new_name TEXT) RETURNS BIGINT LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.resources
        SET name = new_name
        WHERE id = resource_id
        AND is_deleted = FALSE;

        RETURN resource_id;
    END;
    $$;

create or replace function public.update_resource_action (resource_id BIGINT, new_action TEXT) RETURNS BIGINT LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.resources
        SET action = new_action
        WHERE id = resource_id
        AND is_deleted = FALSE;

        RETURN resource_id;
    END;
    $$;

create or replace function public.delete_resource (r_id BIGINT) RETURNS VOID LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.resources
        SET is_deleted = TRUE
        WHERE id = r_id;
    END;
    $$;

--readers
create or replace function public.get_resources_by_id (resource_id BIGINT) RETURNS table (id BIGINT, action VARCHAR(255), name VARCHAR(255)) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.action, r.name
        FROM public.resources r
        WHERE r.id = resource_id
        AND r.is_deleted = FALSE
        LIMIT 1;
    END;
    $$;

create or replace function public.get_resources_by_name (resource_name VARCHAR(255)) RETURNS table (id BIGINT, action VARCHAR(255), name VARCHAR(255)) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.action, r.name
        FROM public.resources r
        WHERE r.name = resource_name
        AND r.is_deleted = FALSE
        LIMIT 1; 
    END;
    $$;

create or replace function public.get_all_resources () RETURNS table (id BIGINT, action VARCHAR(255), name VARCHAR(255)) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.action, r.name
        FROM public.resources r
        WHERE r.is_deleted = FALSE;
    END;
    $$;

-- Roles ----------------------------------------
-- writers
create or replace function public.create_role (r_name VARCHAR(255), r_desc VARCHAR(255)) RETURNS BIGINT LANGUAGE plpgsql as $$
    DECLARE
        new_id BIGINT;
    BEGIN
        INSERT INTO public.roles (name, description)
        VALUES (r_name, r_desc) 
        RETURNING id INTO new_id;

        RETURN new_id;
    END;
    $$;

create or replace function public.update_role_name (role_id BIGINT, new_name TEXT) RETURNS BIGINT LANGUAGE plpgsql as $$
    BEGIN
        
        UPDATE public.roles
        SET name = new_name
        WHERE id = role_id
        AND is_deleted = FALSE;

        RETURN role_id;
    END;
    $$;

create or replace function public.update_role_desc (role_id BIGINT, new_desc TEXT) RETURNS BIGINT LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.roles
        SET description = new_desc
        WHERE id = role_id
        AND is_deleted = FALSE;

        RETURN role_id;
    END;
    $$;

create or replace function public.delete_role (r_id BIGINT) RETURNS VOID LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.roles
        SET is_deleted = TRUE
        WHERE id = r_id;
    END;
    $$;

-- readers
create or replace function public.get_roles_by_id (role_id BIGINT) RETURNS table (
  id BIGINT,
  description VARCHAR(255),
  name VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.description, r.name
        FROM public.roles r
        WHERE r.id = role_id
        AND r.is_deleted = FALSE
        LIMIT 1;
    END;
    $$;

create or replace function public.get_roles_by_name (role_name VARCHAR(255)) RETURNS table (id BIGINT, action VARCHAR(255), name VARCHAR(255)) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.description, r.name
        FROM public.roles r
        WHERE r.name = role_name
        AND r.is_deleted = FALSE
        LIMIT 1; 
    END;
    $$;

create or replace function public.get_all_roles () RETURNS table (
  id BIGINT,
  description VARCHAR(255),
  name VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.description, r.name
        FROM public.roles r
        WHERE r.is_deleted = FALSE;
    END;
    $$;

-- role_resources ----------------------------------------
-- writers
create or replace function public.add_resource_to_role (
  role_identifier BIGINT,
  resource_identifier BIGINT
) RETURNS VOID LANGUAGE plpgsql as $$
    BEGIN
        INSERT INTO public.role_resources (role_id, resource_id)
        VALUES (role_identifier, resource_identifier);
    END;
    $$;

create or replace function public.remove_resource_from_role (
  role_identifier BIGINT,
  resource_identifier BIGINT
) RETURNS VOID LANGUAGE plpgsql as $$
    BEGIN
    UPDATE public.role_resources
    SET is_deleted = TRUE
    WHERE role_id = role_identifier 
        AND resource_id = resource_identifier
        AND is_deleted = FALSE;
    END;
    $$;

-- readers
create or replace function public.get_all_resource_by_role (role_identifier INT) RETURNS table (resource_id BIGINT) LANGUAGE plpgsql as $$
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
create or replace function public.get_users_by_id (user_id INT) RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT u.id, u.FirstName, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
        FROM public.users u
        WHERE u.id = user_id
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
    $$;

create or replace function public.get_users_by_email (user_email VARCHAR(255)) RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT u.id, u.FirstName, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
        FROM public.users u
        WHERE u.email = user_email
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
    $$;

create or replace function public.get_users_by_phone (user_phone VARCHAR(255)) RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT u.id, u.FirstName, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
        FROM public.users u
        WHERE u.phone_number = user_phone
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
    $$;

create or replace function public.get_users_by_username (user_username VARCHAR(255)) RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT u.id, u.FirstName, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
        FROM public.users u
        WHERE u.username = user_username
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
    $$;

create or replace function public.get_users_by_active_status (user_active_status BOOLEAN) RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT u.id, u.FirstName, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
        FROM public.users u
        WHERE u.is_active = user_active_status
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
    $$;

create or replace function public.get_users_by_external_id (user_external_id VARCHAR(255)) RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT u.id, u.FirstName, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
        FROM public.users u
        WHERE u.external_id = user_external_id
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
    $$;

create or replace function public.get_all () RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT u.id, u.FirstName, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
        FROM public.users u
        WHERE u.is_deleted = FALSE
        LIMIT 1;
    END;
    $$;

--writers
create or replace function public.create_user (
  u_FirstName VARCHAR(255),
  u_email VARCHAR(255),
  u_phone VARCHAR(255),
  u_username VARCHAR(255),
  u_dob DATE,
  u_external_id VARCHAR(255)
) RETURNS INT LANGUAGE plpgsql as $$
    DECLARE
        new_id INT;
    BEGIN
        INSERT INTO public.users 
        (FirstName, 
        email, 
        phone_number, 
        username, 
        birth_date,
        external_id)
        VALUES 	
        (u_FirstName, 
        u_email, 
        u_phone, 
        u_username, 
        u_dob, 
        u_external_id)

        RETURNING id INTO new_id;

        RETURN new_id;
    END;
    $$;

create or replace function public.delete_user (u_id INT) RETURNS VOID LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.users
        SET is_deleted = TRUE
        WHERE id = u_id;
    END;
    $$;

create or replace function public.update_user_FirstName (user_id INT, new_FirstName VARCHAR(255)) RETURNS INT LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.users
        SET FirstName = new_FirstName
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
    $$;

create or replace function public.update_user_email (user_id INT, new_email VARCHAR(255)) RETURNS INT LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.users
        SET email = new_email
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
    $$;

create or replace function public.update_user_dob (user_id INT, new_dob DATE) RETURNS INT LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.users
        SET birth_date = new_dob
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
    $$;

create or replace function public.update_user_is_active_status (user_id INT, new_active_status BOOLEAN) RETURNS INT LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.users
        SET is_active = new_active_status
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
    $$;

create or replace function public.update_user_phone (user_id INT, new_phone VARCHAR(255)) RETURNS INT LANGUAGE plpgsql as $$
    BEGIN
        UPDATE public.users
        SET phone_number = new_phone
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
    $$;

create or replace function public.update_user_name (user_id INT, new_username VARCHAR(255)) RETURNS INT LANGUAGE plpgsql as $$
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
create or replace function public.add_role_to_user (user_identifier INT, role_identifier INT) RETURNS VOID LANGUAGE plpgsql as $$
    BEGIN
        INSERT INTO public.user_roles (user_id, role_id)
        VALUES (role_identifier, user_identifier);
    END;
    $$;

create or replace function public.remove_role_from_user (user_identifier INT, role_identifier INT) RETURNS VOID LANGUAGE plpgsql as $$
    BEGIN
    UPDATE public.user_roles
    SET is_deleted = TRUE
    WHERE role_id = role_identifier 
        AND user_id = user_identifier
        AND is_deleted = FALSE;
    END;
    $$;

-- readers
create or replace function public.get_all_role_by_user (user_identifier INT) RETURNS table (role_id INT) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY
        SELECT r.role_id
        FROM public.user_roles r
        WHERE r.user_id = user_identifier 
            AND r.is_deleted = FALSE;
    END;
    $$;

-- Distributor ----------------------------------------
-- writers
create or replace function public.create_distributor () RETURNS BIGINT LANGUAGE plpgsql as $$
    DECLARE
        new_id BIGINT;

    BEGIN
        INSERT INTO public.distributors values(default)
        RETURNING id INTO new_id;
        RETURN new_id;
    END;
    $$;

create or replace function public.create_distributor_location (
  d_general_zone VARCHAR(255),
  d_region VARCHAR(255),
  d_woreda VARCHAR(255),
  d_business_id INT
) RETURNS bigint LANGUAGE plpgsql as $$
    DECLARE
        new_id BIGINT;

    BEGIN
        INSERT INTO public.db_locations(d_general_zone, d_region, d_woreda, d_business_id) 
        VALUES (general_zone, region, woreda, business_id,  distributor_id)
        RETURNING id INTO new_id;
        RETURN new_id;

    END;
    $$;

create or replace function public.delete_distributor(distributor_identifier INT) RETURNS VOID LANGUAGE plpgsql as $$

    BEGIN
        UPDATE public.distributors
        SET is_deleted = TRUE
        WHERE id =distributor_identifier
        AND is_deleted = FALSE;
    END;
    $$;


---readers
create or replace function public.get_distributor_by_id (distributor_id INT) RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY

        SELECT  t2.id, t3.FirstName, t3.email, t3.phone_number, t3.username, t3.birth_date, t3.external_id
        from public.distributor_users t1
        INNER JOIN public.distributors t2 on t1.distributor_id
        INNER JOIN public.users t3 on t2.user_id=t3.id
        WHERE t1.distributor_id = distributor_id AND t1.is_deleted=false
        LIMIT 1;
    END;
$$;

create or replace function public.get_all_distributors() RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql as $$
    BEGIN
        RETURN QUERY

        SELECT  t2.id, t3.FirstName, t3.email, t3.phone_number, t3.username, t3.birth_date, t3.external_id
        from public.distributor_users t1
        INNER JOIN public.distributors t2 on t1.distributor_id
        INNER JOIN public.users t3 on t2.user_id=t3.id
        WHERE  t1.is_deleted=false
        LIMIT 1;
    END;
$$;




-- Distributor Business ----------------------------------------
create or replace function public.create_distributor_business (d_name TEXT, d_tin VARCHAR(10), d_id INT) RETURNS bigint LANGUAGE plpgsql as $$
    DECLARE
        new_id BIGINT;

    BEGIN
        INSERT INTO public.distributor_business_info(d_name, d_tin, d_id) 
        VALUES (name, tin, distributor_id)
        RETURNING id INTO new_id;
        RETURN new_id;

    END;
    $$;

create or replace function public.create_distributor_business_location (
  d_name TEXT,
  d_tin VARCHAR(10),
  d_id INT,
  d_general_zone VARCHAR(255),
  d_region VARCHAR(255),
  d_woreda VARCHAR(255),

) RETURNS bigint LANGUAGE plpgsql as $$
    DECLARE
        new_business_id BIGINT;

    DECLARE new_location_id BIGINT;

    BEGIN
        new_business_id := create_distributor_business(
    d_name,
    d_tin,
    d_id);

    new_location_id := create_distributor_location (
  d_general_zone,
  d_region,
  d_woreda,
  new_business_id
) ;

    RETURN new_business_id;
    END;
    $$;

create or replace function public.update_distributor_business (db_name TEXT, db_tin VARCHAR(10), db_id INT) RETURNS bigint LANGUAGE plpgsql as $$
    DECLARE
        updated_id BIGINT;

    BEGIN
        UPDATE public.distributor_business_info
        SET name= db_name
        WHERE id=db_id;
        UPDATE public.distributor_business_info
        SET tin=db_tin
        WHERE id=db_id
        RETURNING id INTO updated_id;
        RETURN updated_id;
    END;
    $$;

create or replace function public.delete_distributor_business (distributor_id INT) RETURNS VOID LANGUAGE plpgsql as $$

    BEGIN
        UPDATE public.distributor_business_info
        SET is_deleted = TRUE
        WHERE distributor_id = user_id
        AND is_deleted = FALSE;
    END;
    $$;

-- readers
create or replace function public.get_distributor_business (distributor_id INT) RETURNS VOID LANGUAGE plpgsql as $$

    BEGIN
        CALL delete_user(user_id);
        UPDATE public.distributors
        SET is_deleted = TRUE
        WHERE user_id = user_id
        AND is_deleted = FALSE;
    END;
    $$;

-- Distributor User ----------------------------------------
-- writers
create or replace function public.create_distributor_user (
  u_FirstName VARCHAR(255),
  u_email VARCHAR(255),
  u_phone VARCHAR(255),
  u_username VARCHAR(255),
  u_dob DATE,
  u_external_id VARCHAR(255)
) RETURNS BIGINT LANGUAGE plpgsql as $$
    DECLARE new_user_id INT;
    DECLARE new_distributor_id INT;

    BEGIN
    new_user_id := create_user(
    u_FirstName,
    u_email,
    u_phone,
    u_username,
    u_dob,
    u_external_id );

    new_distributor_id := create_distributor();

    INSERT INTO public.distributor_users (user_id, distributor_id)
    VALUES 	
    (new_user_id, new_distributor_id)

    RETURN new_distributor_id;
END;
$$;
