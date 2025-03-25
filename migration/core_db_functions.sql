-- v 0.1
-- Resources ----------------------------------------
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

    --readers

CREATE OR REPLACE FUNCTION public.get_resources_by_id(
    resource_id INT
)
RETURNS TABLE(id INT, 
              action VARCHAR(255), 
              name VARCHAR(255))
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
RETURNS TABLE(id INT, 
              action VARCHAR(255), 
              name VARCHAR(255))
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
RETURNS TABLE(
    id INT, action VARCHAR(255), 
    name VARCHAR(255))
LANGUAGE plpgsql
AS $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.action, r.name
        FROM public.resources r
        WHERE r.is_deleted = FALSE;
    END;
    $$;

create or replace function public.update_resource_action (
    resource_id INT, 
    new_action TEXT) 
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

create or replace function public.delete_resource (
    r_id INT) 
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
create or replace function public.get_resources_by_id (
    resource_id INT) 
RETURNS table (id INT, 
               action VARCHAR(255), 
               name VARCHAR(255)) 
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

create or replace function public.get_resources_by_name (
    resource_name VARCHAR(255)) 
RETURNS table (id INT, 
               action VARCHAR(255), 
               name VARCHAR(255)) 
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

create or replace function public.get_all_resources () 
RETURNS table (id INT, 
               action VARCHAR(255), 
               name VARCHAR(255)) 
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
create or replace function public.create_role (
    r_name VARCHAR(255), 
    r_desc VARCHAR(255)) 
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

create or replace function public.update_role_name (
    role_id INT, 
    new_name TEXT) 
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

create or replace function public.update_role_desc (
    role_id INT, 
    new_desc TEXT) 
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

create or replace function public.delete_role (
    r_id INT) 
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
create or replace function public.get_roles_by_id (
    role_id INT) 
RETURNS table (
  id INT,
  description VARCHAR(255),
  name VARCHAR(255)
) LANGUAGE plpgsql 
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

create or replace function public.get_roles_by_name (
    role_name VARCHAR(255)) 
RETURNS table (id INT, 
               action VARCHAR(255), 
               name VARCHAR(255)) 
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

create or replace function public.get_all_roles () 
RETURNS table (
  id INT,
  description VARCHAR(255),
  name VARCHAR(255)
) LANGUAGE plpgsql 
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
create or replace function public.add_resource_to_role (
  role_identifier INT,
  resource_identifier INT
) RETURNS VOID LANGUAGE plpgsql 
AS $$
    BEGIN
        INSERT INTO public.role_resources (role_id, resource_id)
        VALUES (role_identifier, resource_identifier);
    END;
$$;

create or replace function public.remove_resource_from_role (
  role_identifier INT,
  resource_identifier INT
) RETURNS VOID LANGUAGE plpgsql 
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
create or replace function public.get_all_resource_by_role (
    role_identifier INT) 
RETURNS table (resource_id INT) 
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
create or replace function public.get_users_by_id (
    user_id INT) 
RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT u.id, 
               u.FirstName, 
               u.email, 
               u.phone_number, 
               u.username, 
               u.birth_date, 
               u.is_active, 
               u.external_id 
        FROM public.users u
        WHERE u.id = user_id
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_users_by_email (
    user_email VARCHAR(255)) 
RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT u.id, 
               u.FirstName, 
               u.email, 
               u.phone_number, 
               u.username, 
               u.birth_date, 
               u.is_active, 
               u.external_id 
        FROM public.users u
        WHERE u.email = user_email
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_users_by_phone (
    user_phone VARCHAR(255)) 
RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT u.id, 
               u.FirstName, 
               u.email, 
               u.phone_number, 
               u.username, 
               u.birth_date, 
               u.is_active, 
               u.external_id 
        FROM public.users u
        WHERE u.phone_number = user_phone
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_users_by_username (
    user_username VARCHAR(255)) 
RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT u.id, 
               u.FirstName, 
               u.email, 
               u.phone_number, 
               u.username, 
               u.birth_date, 
               u.is_active, 
               u.external_id 
        FROM public.users u
        WHERE u.username = user_username
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_users_by_active_status (
    user_active_status BOOLEAN) 
RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT u.id, 
               u.FirstName, 
               u.email, 
               u.phone_number, 
               u.username, 
               u.birth_date, 
               u.is_active, 
               u.external_id 
        FROM public.users u
        WHERE u.is_active = user_active_status
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_users_by_external_id (
    user_external_id VARCHAR(255)) 
RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT u.id, u.FirstName, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
        FROM public.users u
        WHERE u.external_id = user_external_id
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_all () 
RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT u.id, 
               u.FirstName, 
               u.email, 
               u.phone_number, 
               u.username, 
               u.birth_date, 
               u.is_active, 
               u.external_id 
        FROM public.users u
        WHERE u.is_deleted = FALSE
        LIMIT 1;
    END;
$$;


CREATE OR REPLACE FUNCTION public.get_all_users()
RETURNS TABLE(
    id INT, 
    fullname VARCHAR(255), 
    email VARCHAR(255), 
    phone VARCHAR(255), 
    username VARCHAR(255), 
    birthdate date, 
    is_active boolean, 
    external_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, 
           u.fullname, 
           u.email, 
           u.phone_number, 
           u.username, 
           u.birth_date, 
           u.is_active, 
           u.external_id 
    FROM public.users u
    WHERE u.is_deleted = FALSE
    LIMIT 1;
END;
$$;

--writers
create or replace function public.create_user (
  u_firstname VARCHAR(255),
  u_lastname VARCHAR(255),
  u_email VARCHAR(255),
  u_phone VARCHAR(255),
  u_username VARCHAR(255),
  u_dob DATE,
  u_external_id VARCHAR(255)
) RETURNS INT LANGUAGE plpgsql 
AS $$
    DECLARE
        new_id INT;
    BEGIN
        INSERT INTO public.users 
        (
        firstName, 
        lastName,
        email, 
        phone_number, 
        username, 
        birth_date,
        external_id)
        VALUES 	
        (
        u_firstname, 
        u_lastname,
        u_email, 
        u_phone, 
        u_username, 
        u_dob, 
        u_external_id)

        RETURNING id INTO new_id;

        RETURN new_id;
    END;
$$;

create or replace function public.delete_user (
    u_id INT) 
    RETURNS VOID LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.users
        SET is_deleted = TRUE
        WHERE id = u_id;
    END;
    $$;

create or replace function public.update_user_FirstName (
    user_id INT, 
    new_FirstName VARCHAR(255)) 
    RETURNS INT LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.users
        SET FirstName = new_FirstName
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
$$;

create or replace function public.update_user_email (
    user_id INT, 
    new_email VARCHAR(255)) 
    RETURNS INT LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.users
        SET email = new_email
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
$$;

create or replace function public.update_user_dob (
    user_id INT, 
    new_dob DATE) 
    RETURNS INT LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.users
        SET birth_date = new_dob
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
$$;

create or replace function public.update_user_is_active_status (
    user_id INT, 
    new_active_status BOOLEAN) 
    RETURNS INT LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.users
        SET is_active = new_active_status
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
$$;

create or replace function public.update_user_phone (
    user_id INT, 
    new_phone VARCHAR(255)) 
    RETURNS INT LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.users
        SET phone_number = new_phone
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
$$;

create or replace function public.update_user_name (
    user_id INT, 
    new_username VARCHAR(255)) 
    RETURNS INT LANGUAGE plpgsql 
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
create or replace function public.add_role_to_user (
    user_identifier INT, 
    role_identifier INT) 
    RETURNS VOID LANGUAGE plpgsql 
AS $$
    BEGIN
        INSERT INTO public.user_roles (user_id, role_id)
        VALUES (role_identifier, user_identifier);
    END;
$$;

create or replace function public.remove_role_from_user (
    user_identifier INT, 
    role_identifier INT) 
    RETURNS VOID LANGUAGE plpgsql 
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
create or replace function public.get_all_role_by_user (
    user_identifier INT) 
RETURNS table (role_id INT) LANGUAGE plpgsql 
AS $$
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
create or replace function public.create_distributor () 
RETURNS INT LANGUAGE plpgsql 
AS $$
    DECLARE
        new_id INT;

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
) RETURNS INT LANGUAGE plpgsql 
AS $$
    DECLARE
        new_id INT;

    BEGIN
        INSERT INTO public.db_locations(
            d_general_zone, 
            d_region, 
            d_woreda, 
            d_business_id) 
        VALUES (general_zone, 
                region, 
                woreda, 
                business_id, 
                distributor_id)
        RETURNING id INTO new_id;
        RETURN new_id;

    END;
$$;

create or replace function public.delete_distributor(
    distributor_identifier INT) 
    RETURNS VOID LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.distributors
        SET is_deleted = TRUE
        WHERE id =distributor_identifier
        AND is_deleted = FALSE;
    END;
    $$;


---readers
create or replace function public.get_distributor_by_id (
    distributor_id INT) 
RETURNS table (
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

        SELECT  t2.id, 
                t3.FirstName, 
                t3.email, 
                t3.phone_number, 
                t3.username, 
                t3.birth_date, 
                t3.external_id
        from public.distributor_users t1
        INNER JOIN public.distributors t2 
        on t1.distributor_id
        INNER JOIN public.users t3 
        on t2.user_id=t3.id
        WHERE t1.distributor_id = distributor_id 
        AND t1.is_deleted=false
        LIMIT 1;
    END;
$$;

create or replace function public.get_all_distributors() 
RETURNS table (
  id INT,
  FirstName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  birthdate date,
  is_active boolean,
  external_id VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY

        SELECT  t2.id, 
                t3.FirstName, 
                t3.email, 
                t3.phone_number, 
                t3.username, 
                t3.birth_date, 
                t3.external_id
        from public.distributor_users t1
        INNER JOIN public.distributors t2 
        on t1.distributor_id
        INNER JOIN public.users t3 
        on t2.user_id=t3.id
        WHERE  t1.is_deleted=false
        LIMIT 1;
    END;
$$;




-- Distributor Business ----------------------------------------
create or replace function public.create_distributor_business (
    d_name TEXT, 
    d_tin VARCHAR(10), 
    d_id INT) 
RETURNS INT LANGUAGE plpgsql 
AS $$
    DECLARE
        new_id INT;

    BEGIN
        INSERT INTO public.distributor_business_info(
            d_name, 
            d_tin, 
            d_id) 
        VALUES (name, 
                tin, 
                distributor_id)
        RETURNING id INTO new_id;
        RETURN new_id;

    END;
$$;


create or replace function public.update_distributor_location (
    d_id INT,
    d_general_zone VARCHAR(255),
    d_region VARCHAR(255),
    d_woreda VARCHAR(255),
    d_business_id INT
) RETURNS INT LANGUAGE plpgsql as $$
    DECLARE
        new_id INT;

    BEGIN
        INSERT INTO public.db_locations(
            d_general_zone, 
            d_region, d_woreda, 
            d_business_id) 
        VALUES (general_zone, 
                region, 
                woreda, 
                business_id,  
                distributor_id)
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
  d_woreda VARCHAR(255)

) RETURNS INT LANGUAGE plpgsql as $$
    DECLARE
        new_business_id INT;

    DECLARE new_location_id INT;

    BEGIN
        new_business_id := create_distributor_business(d_name, 
                                                       d_tin, 
                                                       d_id);
        new_location_id := create_distributor_location (d_general_zone, 
                                                        d_region, 
                                                        d_woreda, 
                                                        new_business_id) ;

    RETURN new_business_id;
    END;
$$;

CREATE or REPLACE FUNCTION public.update_distributor_business (
    db_name TEXT, 
    db_tin VARCHAR(10), 
    db_id INT) 
RETURNS INT LANGUAGE plpgsql 
AS $$
    DECLARE
        updated_id INT;

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

CREATE OR REPLACE FUNCTION public.update_distributor_business_location(
    d_id INT,
    d_name TEXT,
    d_tin VARCHAR(10),
    d_general_zone VARCHAR(255),
    d_region VARCHAR(255),
    d_woreda VARCHAR(255)

) RETURNS INT 
LANGUAGE plpgsql 
AS $$
    DECLARE
        updated_business_id INT;

    DECLARE updated_location_id INT;

    BEGIN
        updated_business_id := 
        update_distributor_business(d_name, d_tin, d_id);

    UPDATE public.db_locations
            SET general_zone= d_general_zone
            WHERE business_id=d_id;
            UPDATE public.db_locations
            SET region=d_region
            WHERE id=db_id;
            UPDATE public.db_locations
            SET woreda=d_woreda
            WHERE id=db_id
            RETURNING id INTO updated_location_id;

    RETURN updated_business_id;
    END;
$$;

create or replace function public.delete_distributor_business (
    distributor_id INT) 
RETURNS VOID 
LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.distributor_business_info
        SET is_deleted = TRUE
        WHERE distributor_id = user_id
        AND is_deleted = FALSE;
    END;
$$;

-- readers
CREATE or REPLACE function public.get_distributor_business (distributor_id INT) 
RETURNS TABLE(
    id INT,
    name TEXT,
    tin VARCHAR(10),
    general_zone VARCHAR(255),
    region VARCHAR(255),
    woreda VARCHAR(255)) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        SELECT id, name, tin FROM public.distributor_business_info t1
        JOIN public.db_locations t2 on t1.id
        WHERE distributor_id = distributor_id 
        AND is_deleted = FALSE;
    END;
$$;

CREATE or REPLACE function public.get_business_by_id (business_id INT) 
RETURNS TABLE(
    id INT,
    name TEXT,
    tin VARCHAR(10),
    general_zone VARCHAR(255),
    region VARCHAR(255),
    woreda VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        SELECT id, name, tin FROM public.distributor_business_info t1
        JOIN public.db_locations t2 on t1.id
        WHERE id = business_id
        AND is_deleted = FALSE;
    END;
$$;

CREATE or REPLACE function public.get_all_businesses () 
RETURNS TABLE(
    id INT,
    name TEXT,
    tin VARCHAR(10),
    general_zone VARCHAR(255),
    region VARCHAR(255),
    woreda VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        SELECT id, name, tin FROM public.distributor_business_info t1
        JOIN public.db_locations t2 on t1.id
        WHERE is_deleted = FALSE;
    END;
$$;



-- Distributor User ----------------------------------------
    
    -- writers
create or replace function public.create_distributor_user (
  u_firstname VARCHAR(255),
  u_lastname VARCHAR(255),
  u_email VARCHAR(255),
  u_phone VARCHAR(255),
  u_username VARCHAR(255),
  u_dob DATE,
  u_external_id VARCHAR(255)
) RETURNS INT 
LANGUAGE plpgsql 
AS $$
    DECLARE new_user_id INT;
    DECLARE new_distributor_id INT;

    BEGIN
    new_user_id := create_user(
    u_firstname,
    u_lastname,
    u_email,
    u_phone,
    u_username,
    u_dob,
    u_external_id );

    new_distributor_id := create_distributor();

    INSERT INTO public.distributor_users (user_id, distributor_id)
    VALUES 	
    (new_user_id, new_distributor_id);

    RETURN new_distributor_id;
END;
$$;

CREATE or REPLACE FUNCTION public.update_distributor_business (
    db_name TEXT, 
    db_tin VARCHAR(10), 
    db_id INT) 
RETURNS INT LANGUAGE plpgsql 
AS $$
    DECLARE
        updated_id INT;

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

CREATE OR REPLACE FUNCTION public.update_distributor_business_location(
    d_id INT,
    d_name TEXT,
    d_tin VARCHAR(10),
    d_general_zone VARCHAR(255),
    d_region VARCHAR(255),
    d_woreda VARCHAR(255)
) RETURNS INT 
LANGUAGE plpgsql 
AS $$
    DECLARE
        updated_business_id INT;

    DECLARE updated_location_id INT;

    BEGIN
        updated_business_id := update_distributor_business(
    d_name,
    d_tin,
    d_id);


    UPDATE public.db_locations
            SET general_zone= d_general_zone
            WHERE business_id=d_id;
            UPDATE public.db_locations
            SET region=d_region
            WHERE id=db_id;
            UPDATE public.db_locations
            SET woreda=d_woreda
            WHERE id=db_id
            RETURNING id INTO updated_location_id;

    RETURN updated_business_id;
    END;
$$;


-- Retailer ----------------------------------------
    
    -- writers
CREATE OR REPLACE FUNCTION public.create_retailer (
    r_name VARCHAR(255),
    r_tin VARCHAR(255),
    r_lat VARCHAR(255),
    r_long VARCHAR(255),
    r_generalZone VARCHAR(255),
    r_region VARCHAR(255),
    r_woreda VARCHAR(255)
) 
RETURNS INT 
LANGUAGE plpgsql 
AS $$
    DECLARE
        new_id INT;
        business_id INT;
    BEGIN
        INSERT INTO public.retailers DEFAULT VALUES
        RETURNING id INTO new_id;

        -- Business info
        INSERT INTO public.retailer_business_info(name, tin, retailer_id)
        VALUES(r_name, r_tin, new_id)
        RETURNING id INTO business_id;

        -- business Locations
        INSERT INTO public.rb_locations(lat, long, general_zone, region, woreda, business_id)
        VALUES(r_lat, r_long, r_generalZone, r_region, r_woreda, business_id);  -- Added business_id

        RETURN new_id;
    END;
$$;

CREATE OR REPLACE FUNCTION public.update_retailer_name(
    r_id INT,
    r_name VARCHAR(255)
) 
RETURNS VOID
LANGUAGE plpgsql 
AS $$
    DECLARE new_id INT;
    BEGIN
        UPDATE public.retailer_business_info
            SET name = r_name
            WHERE id = r_id;
    END;
$$;


CREATE OR REPLACE FUNCTION public.update_retailer_tin(
    r_id INT,
    r_tin VARCHAR(255)
) 
RETURNS VOID 
LANGUAGE plpgsql 
AS $$
    DECLARE new_id INT;
    BEGIN
        UPDATE public.retailer_business_info
            SET tin = r_tin
            WHERE id = r_id;
    END;
$$;

    -- readers
create or replace function public.get_retailer_by_id (
    r_retailer_id INT
) 
RETURNS TABLE (
  id INT,
  name VARCHAR(255),
  tin VARCHAR(255),
  lat VARCHAR(255),
  long VARCHAR(255),
  generalZone VARCHAR(255),
  region VARCHAR(255),
  woreda VARCHAR(255)
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY

        SELECT  r.id, rb.name, rb.tin, rb_loc.lat, rb_loc.long, 
        rb_loc.general_zone, rb_loc.region, rb_loc.woreda
        FROM  public.retailers r
        JOIN public.retailer_business_info rb 
        ON rb.retailer_id = r.id
        JOIN public.rb_locations rb_loc 
        ON rb_loc.business_id = rb.id
        WHERE r.id = r_retailer_id 
        AND r.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_retailer_by_name (
    retailer_name VARCHAR(255)
) 
RETURNS TABLE (
  id INT,
  name VARCHAR(255),
  tin VARCHAR(255),
  lat VARCHAR(255),
  long VARCHAR(255),
  generalZone VARCHAR(255),
  region VARCHAR(255),
  woreda VARCHAR(255)
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY

       SELECT  r.id, rb.name, rb.tin, rb_loc.lat, rb_loc.long, 
        rb_loc.general_zone, rb_loc.region, rb_loc.woreda
        FROM  public.retailers r
        JOIN public.retailer_business_info rb 
        ON rb.retailer_id = r.id
        JOIN public.rb_locations rb_loc 
        ON rb_loc.business_id = rb.id
        WHERE rb.name = retailer_name 
        AND r.is_deleted = FALSE;
    END;
$$;

create or replace function public.get_retailer_by_tin (
    retailer_tin VARCHAR(255)
) 
RETURNS TABLE (
  id INT,
  name VARCHAR(255),
  tin VARCHAR(255),
  lat VARCHAR(255),
  long VARCHAR(255),
  generalZone VARCHAR(255),
  region VARCHAR(255),
  woreda VARCHAR(255)
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY

       SELECT  r.id, rb.name, rb.tin, rb_loc.lat, rb_loc.long, 
        rb_loc.general_zone, rb_loc.region, rb_loc.woreda
        FROM  public.retailers r
        JOIN public.retailer_business_info rb 
        ON rb.retailer_id = r.id
        JOIN public.rb_locations rb_loc 
        ON rb_loc.business_id = rb.id
        WHERE rb.tin = retailer_tin 
        AND r.is_deleted = FALSE;
    END;
$$;

create or replace function public.get_all_retailers () 
RETURNS TABLE (
  id INT,
  name VARCHAR(255),
  tin VARCHAR(255),
  lat VARCHAR(255),
  long VARCHAR(255),
  generalZone VARCHAR(255),
  region VARCHAR(255),
  woreda VARCHAR(255)
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY

        SELECT r.id, rb.name, rb.tin, rb_loc.lat, rb_loc.long, 
        rb_loc.general_zone, rb_loc.region, rb_loc.woreda
        FROM  public.retailers r
        JOIN public.retailer_business_info rb 
        ON rb.retailer_id = r.id
        JOIN public.rb_locations rb_loc 
        ON rb_loc.business_id = rb.id;
    END;
$$;

-- retailer user agent ---------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.create_retailer_user (
    r_id INT,
    r_user_id INT
) 
RETURNS VOID 
LANGUAGE plpgsql 
AS $$
BEGIN   
    INSERT INTO public.retailer_users(user_id, retailer_id)
    VALUES(r_user_id, r_id);
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_all_retailer_users (
    r_id INT
) 
RETURNS TABLE (
  id INT
) 
LANGUAGE plpgsql 
AS $$
BEGIN   
    RETURN QUERY
    SELECT user_id
    FROM public.retailer_users ru
    WHERE ru.retailer_id = r_id;
END;
$$;

-- invoices --------------------------------------------

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
    new_id INT;
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
RETURNS TABLE(id INT, 
              status VARCHAR(255), 
              external_id VARCHAR(255), 
              order_id INT, 
              subtotal DECIMAL(12,2), 
              tax_amount DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, 
           i.status, 
           i.external_id, 
           i.order_id, 
           i.subtotal, 
           i.tax_amount
    FROM public.invoices i
    WHERE i.is_deleted = FALSE 
    AND i.external_id = i_external_id;
END;
$$;


CREATE OR REPLACE FUNCTION public.get_invoices_by_status(
     i_status VARCHAR(255)
)
RETURNS TABLE(id INT, 
              status VARCHAR(255), 
              external_id VARCHAR(255), 
              order_id INT, 
              subtotal DECIMAL(12,2), 
              tax_amount DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, 
           i.status, 
           i.external_id, 
           i.order_id, 
           i.subtotal, 
           i.tax_amount
    FROM public.invoices i
    WHERE i.is_deleted = FALSE 
    AND i.status = i_status;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_invoices_by_order_id(
     i_order_id INT
)
RETURNS TABLE(id INT, 
              status VARCHAR(255), 
              external_id VARCHAR(255), 
              order_id INT, 
              subtotal DECIMAL(12,2), 
              tax_amount DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, 
           i.status, 
           i.external_id, 
           i.order_id, 
           i.subtotal, 
           i.tax_amount
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
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.invoice_line_items (product_name, 
                                           qty, 
                                           price, 
                                           product_id, 
                                           invoice_id)
    VALUES (i_product_name, 
            i_qty, 
            i_price, 
            i_product_id, 
            i_invoice_id) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;
    
    -- readers
CREATE OR REPLACE FUNCTION public.get_invoice_line_item_by_invoice_id(
    i_invoice_id INT
)
RETURNS TABLE(id INT, 
              product_name VARCHAR(255), 
              qty INT, 
              price DECIMAL(12,2), 
              product_id INT, 
              invoice_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, 
           i.product_name, 
           i.qty, 
           i.price, 
           i.product_id, 
           i.invoice_id
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
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.products (name, 
                                 description, 
                                 external_id, 
                                 distributor_id)
    VALUES (p_product_name, 
            p_product_description, 
            p_external_id, 
            p_distributor_id) 
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
RETURNS TABLE(id INT, 
              product_name VARCHAR(255), 
              product_description TEXT, 
              external_id VARCHAR(255), 
              is_active BOOLEAN, 
              distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id
    FROM public.products p
    WHERE p.id = p_product_id
      AND p.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_name(
    p_product_name VARCHAR(255)
)
RETURNS TABLE(id INT, 
              product_name VARCHAR(255), 
              product_description TEXT, 
              external_id VARCHAR(255), 
              is_active BOOLEAN, 
              distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id
    FROM public.products p
    WHERE p.name = p_product_name
      AND p.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_externalId(
    p_external_id VARCHAR(255)
)
RETURNS TABLE(id INT, 
              product_name VARCHAR(255), 
              product_description TEXT, 
              external_id VARCHAR(255), 
              is_active BOOLEAN, 
              distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id
    FROM public.products p
    WHERE p.external_id = p_external_id
      AND p.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_distributorId(
    p_distributor_id INT
)
RETURNS TABLE(id INT, 
              product_name VARCHAR(255), 
              product_description TEXT, 
              external_id VARCHAR(255), 
              is_active BOOLEAN, 
              distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id
    FROM public.products p
    WHERE p.distributor_id = p_distributor_id
      AND p.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_products()
RETURNS TABLE(id INT, 
              product_name VARCHAR(255), 
              product_description TEXT, 
              external_id VARCHAR(255), 
              is_active BOOLEAN, 
              distributor_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id
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

-- category -------------------------------------------------------------

    -- writer
CREATE OR REPLACE FUNCTION public.create_category(
  c_name VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO 
    public.category (name)
    VALUES (c_name)
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.remove_category(
    c_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.category
    SET is_deleted = TRUE
    WHERE id = c_id;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_category(
    c_id INT
)
RETURNS TABLE(id INT, name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT c.id, c.name
    FROM public.category c
    WHERE c.id = c_id
      AND c.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_category()
RETURNS TABLE(id INT, name VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT c.id, c.name
    FROM public.category c
    WHERE c.is_deleted = FALSE;
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
    p_product_id INT
)
RETURNS TABLE(category_id INT)
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
    p_category_id INT
)
RETURNS TABLE(product_id INT)
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
    INSERT INTO public.s_ledger (quantity, 
                                 product_id, 
                                 operation, 
                                 created_by_user_id)
    VALUES (s_quantity, 
            s_product_id, 
            s_stock_operation, 
            s_created_by_user_id);
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
    INSERT INTO public.p_attribute_values (name, 
                                           product_id, 
                                           attribute_id)
    VALUES (i_name, 
            i_product_id, 
            i_attribute_id);
END;
$$;
    -- reader

