-- v 0.1
-- Resources ----------------------------------------
   
    -- writers
CREATE OR REPLACE FUNCTION public.create_resource(
   r_name VARCHAR(255),
   r_action VARCHAR(255),
   r_resource VARCHAR(255),
   r_scope VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.resources (name, action, resource, scope)
    VALUES (r_name, r_action, r_resource, r_scope) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_resource_name(
    resource_id INT,
    new_name VARCHAR(255)
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

CREATE OR REPLACE FUNCTION public.update_resource_resource(
    resource_id INT,
    new_resource VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.resources
    SET resource = new_resource
    WHERE id = resource_id
      AND is_deleted = FALSE;

    RETURN resource_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_resource_action(
    resource_id INT,
    new_action VARCHAR(255)
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

CREATE OR REPLACE FUNCTION public.update_resource_scope(
    resource_id INT,
    new_scope VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.resources
    SET scope = new_scope
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
              name VARCHAR(255),
              resource VARCHAR(255),
              scope VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.action, r.name, r.resource, r.scope
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
              name VARCHAR(255),
              resource VARCHAR(255),
              scope VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.action, r.name, r.resource, r.scope
    FROM public.resources r
    WHERE r.name = resource_name
      AND r.is_deleted = FALSE
    LIMIT 1; 
END;
$$;

CREATE OR REPLACE FUNCTION public.get_resources_by_scope(
    scope_name VARCHAR(255)
)
RETURNS TABLE(id INT, 
              action VARCHAR(255), 
              name VARCHAR(255),
              resource VARCHAR(255),
              scope VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.action, r.name, r.resource, r.scope
    FROM public.resources r
    WHERE r.scope = scope_name
      AND r.is_deleted = FALSE
    LIMIT 1; 
END;
$$;

CREATE OR REPLACE FUNCTION public.get_resources_by_resource(
    r_resource VARCHAR(255)
)
RETURNS TABLE(id INT, 
              action VARCHAR(255), 
              name VARCHAR(255),
              resource VARCHAR(255),
              scope VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT r.id, r.action, r.name, r.resource, r.scope
    FROM public.resources r
    WHERE r.resource = r_resource
      AND r.is_deleted = FALSE
    LIMIT 1; 
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_resources(
    r_limit INT,
    r_offset INT
)
RETURNS TABLE(
    id INT, 
    action VARCHAR(255), 
    name VARCHAR(255),
    resource VARCHAR(255),
    scope VARCHAR(255))
LANGUAGE plpgsql
AS $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.action, r.name, r.resource, r.scope
        FROM public.resources r
        WHERE r.is_deleted = FALSE
        LIMIT r_limit
        OFFSET r_offset;
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
    new_name VARCHAR(255)) 
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
    new_desc VARCHAR(255)) 
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
  name VARCHAR(255),
  description VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.name, r.description
        FROM public.roles r
        WHERE r.id = role_id
        AND r.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_roles_by_name (
    role_name VARCHAR(255)) 
RETURNS table (id INT, 
               name VARCHAR(255), 
               description VARCHAR(255)) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.name, r.description
        FROM public.roles r
        WHERE r.name = role_name
        AND r.is_deleted = FALSE
        LIMIT 1; 
    END;
$$;

create or replace function public.get_all_roles (
    r_limit INT,
    r_offset INT
) 
RETURNS table (
  id INT,
  name VARCHAR(255),
  description VARCHAR(255)
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT r.id, r.name, r.description
        FROM public.roles r
        WHERE r.is_deleted = FALSE
        LIMIT r_limit
        OFFSET r_offset;
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
RETURNS table (
resource_id INT,
name VARCHAR(255),
action VARCHAR(255),
resource VARCHAR(255)
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT r.resource_id, re.name, re.action, re.resource
        FROM public.role_resources r
        JOIN resources re on re.id = r.resource_id
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
  firstName VARCHAR(255),
  lastName VARCHAR(255),
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
               u.firstName,
               u.lastName,
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
  firstName VARCHAR(255),
  lastName VARCHAR(255),
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
               u.firstName,
               u.lastName,
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
  firstName VARCHAR(255),
  lastName VARCHAR(255),
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
               u.firstName,
               u.lastName, 
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
  firstName VARCHAR(255),
  lastName VARCHAR(255),
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
               u.firstName,
               u.lastName,
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
    user_active_status BOOLEAN,
    u_limit INT,
    u_offset INT
) 
RETURNS table (
  id INT,
  firstName VARCHAR(255),
  lastName VARCHAR(255),
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
               u.firstName,
               u.lastName,
               u.email, 
               u.phone_number, 
               u.username, 
               u.birth_date, 
               u.is_active, 
               u.external_id 
        FROM public.users u
        WHERE u.is_active = user_active_status
        AND u.is_deleted = FALSE
        LIMIT r_limit
        OFFSET r_offset;
    END;
$$;

create or replace function public.get_users_by_external_id (
    user_external_id VARCHAR(255)) 
RETURNS table (
  id INT,
  firstName VARCHAR(255),
  lastName VARCHAR(255),
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
        SELECT u.id, u.firstName, u.lastName, u.email, u.phone_number, u.username, u.birth_date, u.is_active, u.external_id 
        FROM public.users u
        WHERE u.external_id = user_external_id
        AND u.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_all () 
RETURNS table (
  id INT,
  firstName VARCHAR(255),
  lastName VARCHAR(255),
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
               u.firstName, 
               u.lastName,
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


CREATE OR REPLACE FUNCTION public.get_all_users(
    u_limit INT,
    u_offset INT
)
RETURNS TABLE(
    id INT, 
    firstName VARCHAR(255),
    lastName VARCHAR(255), 
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
           u.firstName, 
           u.lastName,
           u.email, 
           u.phone_number, 
           u.username, 
           u.birth_date, 
           u.is_active, 
           u.external_id 
    FROM public.users u
    WHERE u.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE function public.get_user_provider(id INT) 
RETURNS TABLE(user_id INT,
              provider_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT up.user_id, up.provider_id 
    FROM public.user_providers up 
    WHERE up.is_deleted = FALSE 
    AND up.user_id=id;

END;
$$;

CREATE OR REPLACE function public.get_user_provider_by_email(u_email VARCHAR(255)) 
RETURNS TABLE(user_id INT,
              provider_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT up.user_id, up.provider_id 
    FROM public.user_providers up
    JOIN public.users u ON u.id=up.user_id
    WHERE u.is_deleted = FALSE 
    AND u.email=u_email;
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

create or replace function public.create_user_provider(
u_id INT,
p_id VARCHAR(255)
) RETURNS INT LANGUAGE plpgsql 
AS $$
    DECLARE
        new_id INT;
    BEGIN
        INSERT INTO public.user_providers
        (user_id, 
         provider_id) VALUES(u_id, 
                             p_id);

        RETURN new_id;
    END;
$$;

CREATE OR REPLACE FUNCTION public.create_and_activate_user (
    u_firstname VARCHAR(255),
    u_lastname VARCHAR(255),
    u_email VARCHAR(255),
    u_phone VARCHAR(255),
    u_username VARCHAR(255),
    u_dob DATE,
    u_external_id VARCHAR(255)
) 
RETURNS INT 
LANGUAGE plpgsql 
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
        external_id,
        is_active
    )
    VALUES(
        u_firstname, 
        u_lastname,
        u_email, 
        u_phone, 
        u_username, 
        u_dob, 
        u_external_id,
        TRUE
    )
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.delete_user (
    u_id INT) 
    RETURNS VOID LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.users
        SET is_deleted = TRUE
        WHERE id = u_id;

        PERFORM public.remove_user_connection(u_id);
    END;
$$;

create or replace function public.update_user_FirstName (
    user_id INT, 
    new_FirstName VARCHAR(255)) 
    RETURNS INT LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.users
        SET firstName = new_FirstName
        WHERE id = user_id
        AND is_deleted = FALSE;

        RETURN user_id;
    END;
$$;

create or replace function public.update_user_lastName (
    user_id INT, 
    new_lastName VARCHAR(255)) 
    RETURNS INT LANGUAGE plpgsql 
AS $$
    BEGIN
        UPDATE public.users
        SET lastName = new_lastName
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
CREATE OR REPLACE FUNCTION public.create_distributor (
    d_name VARCHAR(255),
    d_tin VARCHAR(255),
    d_lat VARCHAR(255),
    d_long VARCHAR(255),
    d_generalZone VARCHAR(255),
    d_region VARCHAR(255),
    d_woreda VARCHAR(255)
) 
RETURNS INT 
LANGUAGE plpgsql 
AS $$
    DECLARE
        new_id INT;
        business_id INT;
    BEGIN
        INSERT INTO public.distributors DEFAULT VALUES
        RETURNING id INTO new_id;

        -- Business info
        INSERT INTO public.distributor_business_info(name, tin, distributor_id)
        VALUES(d_name, d_tin, new_id)
        RETURNING id INTO business_id;

        -- business Locations
        INSERT INTO public.db_locations(lat, long, general_zone, region, woreda, business_id)
        VALUES(d_lat, d_long, d_generalZone, d_region, d_woreda, business_id);

        RETURN new_id;
    END;
$$;

CREATE OR REPLACE FUNCTION public.update_distributor_name(
    r_id INT,
    r_name VARCHAR(255)
) 
RETURNS VOID
LANGUAGE plpgsql 
AS $$
    DECLARE new_id INT;
    BEGIN
        UPDATE public.distributor_business_info
            SET name = r_name
            WHERE id = r_id;
    END;
$$;


CREATE OR REPLACE FUNCTION public.update_distributor_tin(
    r_id INT,
    r_tin VARCHAR(255)
) 
RETURNS VOID 
LANGUAGE plpgsql 
AS $$
    DECLARE new_id INT;
    BEGIN
        UPDATE public.distributor_business_info
            SET tin = r_tin
            WHERE id = r_id;
    END;
$$;

CREATE OR REPLACE FUNCTION public.activate_distributors(
    d_id INT
) 
RETURNS VOID 
LANGUAGE plpgsql 
AS $$
BEGIN
    UPDATE public.distributors
    SET is_active = TRUE
    WHERE id = d_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.deactivate_distributors(
    d_id INT
) 
RETURNS VOID 
LANGUAGE plpgsql 
AS $$
BEGIN
    UPDATE public.distributors
    SET is_active = FALSE
    WHERE id = d_id;
END;
$$;

    -- readers
create or replace function public.get_distributor_by_id (
    d_distributor_id INT
) 
RETURNS TABLE (
  id INT,
  name VARCHAR(255),
  tin VARCHAR(255),
  lat VARCHAR(255),
  long VARCHAR(255),
  generalZone VARCHAR(255),
  region VARCHAR(255),
  woreda VARCHAR(255),
  is_active BOOLEAN,
  verdict VARCHAR(255)
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY

        SELECT  d.id, db.name, db.tin, db_loc.lat, db_loc.long, 
        db_loc.general_zone, db_loc.region, db_loc.woreda, d.is_active, dr.verdict
        FROM  public.distributors d
        JOIN public.distributor_business_info db 
        ON db.distributor_id = d.id
        JOIN public.db_locations db_loc 
        ON db_loc.business_id = db.id
        LEFT JOIN distributor_reviews dr
        ON dr.distributor_id = d.id
        WHERE d.id = d_distributor_id 
        AND d.is_deleted = FALSE
        LIMIT 1;
    END;
$$;

create or replace function public.get_distributor_by_user_id (
    d_user_id INT
) 
RETURNS TABLE (
  id INT,
  name VARCHAR(255),
  tin VARCHAR(255),
  lat VARCHAR(255),
  long VARCHAR(255),
  generalZone VARCHAR(255),
  region VARCHAR(255),
  woreda VARCHAR(255),
  is_active BOOLEAN,
  verdict VARCHAR(255)
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT d.id, db.name, db.tin, db_loc.lat, db_loc.long, 
            db_loc.general_zone, db_loc.region, db_loc.woreda, d.is_active, 
            COALESCE(dr.verdict, 'PENDING') AS verdict
        FROM public.distributors d
        JOIN public.distributor_business_info db ON db.distributor_id = d.id
        JOIN public.db_locations db_loc ON db_loc.business_id = db.id
        LEFT JOIN distributor_reviews dr ON dr.distributor_id = d.id
        JOIN distributor_users du ON du.distributor_id = d.id
        WHERE du.user_id = 92 AND d.is_deleted = FALSE;
    END;
$$;

CREATE OR REPLACE FUNCTION public.get_distributor_by_name (
    d_distributor_name VARCHAR(255),
    t_limit INT,
    t_offset INT
) 
RETURNS TABLE (
    id INT,
    name VARCHAR(255),
    tin VARCHAR(255),
    lat VARCHAR(255),
    long VARCHAR(255),
    generalZone VARCHAR(255),
    region VARCHAR(255),
    woreda VARCHAR(255),
    is_active BOOLEAN
) 
LANGUAGE plpgsql 
AS $$
BEGIN
    RETURN QUERY

    SELECT d.id, db.name, db.tin, db_loc.lat, db_loc.long, 
           db_loc.general_zone, db_loc.region, db_loc.woreda, d.is_active
    FROM public.distributors d
    JOIN public.distributor_business_info db 
        ON db.distributor_id = d.id
    JOIN public.db_locations db_loc 
        ON db_loc.business_id = db.id
    WHERE db.name ILIKE '%' || d_distributor_name || '%'
      AND d.is_deleted = FALSE
    LIMIT t_limit
    OFFSET t_offset;
END;
$$;

create or replace function public.get_distributor_by_tin (
    distributor_tin VARCHAR(255)
) 
RETURNS TABLE (
  id INT,
  name VARCHAR(255),
  tin VARCHAR(255),
  lat VARCHAR(255),
  long VARCHAR(255),
  generalZone VARCHAR(255),
  region VARCHAR(255),
  woreda VARCHAR(255),
  is_active BOOLEAN
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY

       SELECT  d.id, db.name, db.tin, db_loc.lat, db_loc.long, 
        db_loc.general_zone, db_loc.region, db_loc.woreda, d.is_active
        FROM  public.distributors d
        JOIN public.distributor_business_info db 
        ON db.distributor_id = d.id
        JOIN public.db_locations db_loc 
        ON db_loc.business_id = db.id
        WHERE db.tin = distributor_tin 
        AND d.is_deleted = FALSE;
    END;
$$;

CREATE OR REPLACE FUNCTION public.get_distributor_by_status (
    d_status BOOLEAN,
    t_limit INT,
    t_offset INT
) 
RETURNS TABLE (
    id INT,
    name VARCHAR(255),
    tin VARCHAR(255),
    lat VARCHAR(255),
    long VARCHAR(255),
    generalZone VARCHAR(255),
    region VARCHAR(255),
    woreda VARCHAR(255),
    is_active BOOLEAN,
    verdict VARCHAR(255),
    total_count BIGINT
) 
LANGUAGE plpgsql 
AS $$
BEGIN
    RETURN QUERY

        SELECT d.id, db.name, db.tin, db_loc.lat, db_loc.long, 
        db_loc.general_zone, db_loc.region, db_loc.woreda, d.is_active, dr.verdict, COUNT(*) OVER() AS total_count
        FROM  public.distributors d
        JOIN public.distributor_business_info db 
        ON db.distributor_id = d.id
        JOIN public.db_locations db_loc 
        ON db_loc.business_id = db.id
        JOIN distributor_reviews dr
        ON dr.distributor_id = d.id
    WHERE d.is_active = d_status
      AND d.is_deleted = FALSE
    LIMIT t_limit
    OFFSET t_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_distributor_by_approval_status (
    d_status VARCHAR(255),
    t_limit INT,
    t_offset INT
) 
RETURNS TABLE (
    id INT,
    name VARCHAR(255),
    tin VARCHAR(255),
    lat VARCHAR(255),
    long VARCHAR(255),
    generalZone VARCHAR(255),
    region VARCHAR(255),
    woreda VARCHAR(255),
    is_active BOOLEAN,
    verdict VARCHAR(255),
    total_count BIGINT
) 
LANGUAGE plpgsql 
AS $$
BEGIN
    RETURN QUERY

        SELECT d.id, db.name, db.tin, db_loc.lat, db_loc.long, 
        db_loc.general_zone, db_loc.region, db_loc.woreda, d.is_active, dr.verdict, COUNT(*) OVER() AS total_count
        FROM  public.distributors d
        JOIN public.distributor_business_info db 
        ON db.distributor_id = d.id
        JOIN public.db_locations db_loc 
        ON db_loc.business_id = db.id
        JOIN distributor_reviews dr
        ON dr.distributor_id = d.id
    WHERE dr.verdict = d_status
      AND d.is_deleted = FALSE
    LIMIT t_limit
    OFFSET t_offset;
END;
$$;


create or replace function public.get_all_distributors (
    t_limit INT,
    t_offset INT
) 
RETURNS TABLE (
  id INT,
  name VARCHAR(255),
  tin VARCHAR(255),
  lat VARCHAR(255),
  long VARCHAR(255),
  generalZone VARCHAR(255),
  region VARCHAR(255),
  woreda VARCHAR(255),
  is_active BOOLEAN,
  verdict VARCHAR(255),
  total_count BIGINT
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY

        SELECT d.id, db.name, db.tin, db_loc.lat, db_loc.long, 
        db_loc.general_zone, db_loc.region, db_loc.woreda, d.is_active, dr.verdict, COUNT(*) OVER() AS total_count
        FROM  public.distributors d
        JOIN public.distributor_business_info db 
        ON db.distributor_id = d.id
        JOIN public.db_locations db_loc 
        ON db_loc.business_id = db.id
        JOIN distributor_reviews dr
        ON dr.distributor_id = d.id
        WHERE d.is_deleted = FALSE
        LIMIT t_limit
        OFFSET t_offset;
    END;
$$;

-- distributor user agent ---------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.create_distributor_user (
    d_id INT,
    d_user_id INT
) 
RETURNS INT 
LANGUAGE plpgsql 
AS $$
DECLARE
    new_id INT;
BEGIN   
    INSERT INTO public.distributor_users(user_id, distributor_id)
    VALUES(d_user_id, d_id);

    RETURN d_user_id;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_all_distributor_users (
    d_id INT,
    t_limit INT,
    t_offset INT
) 
RETURNS TABLE (
  id INT
) 
LANGUAGE plpgsql 
AS $$
BEGIN   
    RETURN QUERY
    SELECT user_id
    FROM public.distributor_users du
    WHERE du.distributor_id = d_id
     AND du.is_deleted = FALSE
    LIMIT t_limit
    OFFSET t_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_distributor_user_details (
    d_id INT,
    t_limit INT,
    t_offset INT
) 
RETURNS TABLE (
  id INT,
  firstName VARCHAR(255),
  lastName VARCHAR(255),
  email VARCHAR(255),
  phone VARCHAR(255),
  username VARCHAR(255),
  total_count BIGINT
) LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY
        SELECT u.id, 
               u.firstName,
               u.lastName,
               u.email, 
               u.phone_number, 
               u.username, 
               COUNT(*) OVER() AS total_count
    FROM public.distributor_users du
    JOIN users u ON du.user_id = u.id 
    WHERE du.distributor_id = d_id
    AND du.is_deleted = FALSE
    LIMIT t_limit
    OFFSET t_offset;
END;
$$;


-- distributor reveiw ---------------------------------

    -- writer
CREATE OR REPLACE FUNCTION public.approve_distributor_review (
    d_distributor_id INT,
    d_comment VARCHAR(255),
    d_reviewed_by VARCHAR(255)
) 
RETURNS VOID
LANGUAGE plpgsql 
AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM public.distributor_reviews WHERE distributor_id = d_distributor_id) THEN
        UPDATE public.distributor_reviews
        SET verdict = 'APPROVED',
            comment = d_comment,
            reviewed_by = d_reviewed_by
        WHERE distributor_id = d_distributor_id;
    ELSE
        INSERT INTO public.distributor_reviews(distributor_id,
                                               verdict,
                                               comment,
                                               reviewed_by)
        VALUES(d_distributor_id,
               'APPROVED',
               d_comment,
               d_reviewed_by);
    END IF;
END;
$$;

CREATE OR REPLACE FUNCTION public.reject_distributor_review (
    d_distributor_id INT,
    d_comment VARCHAR(255),
    d_reviewed_by VARCHAR(255)
) 
RETURNS VOID
LANGUAGE plpgsql 
AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM public.distributor_reviews WHERE distributor_id = d_distributor_id) THEN
        UPDATE public.distributor_reviews
        SET verdict = 'REJECTED',
            comment = d_comment,
            reviewed_by = d_reviewed_by
        WHERE distributor_id = d_distributor_id;
    ELSE
        INSERT INTO public.distributor_reviews(distributor_id,
                                               verdict,
                                               comment,
                                               reviewed_by)
        VALUES(d_distributor_id,
               'REJECTED',
               d_comment,
               d_reviewed_by);
    END IF;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_approval_status_by_distributor_id (
    d_id INT
) 
RETURNS VARCHAR(255)
LANGUAGE plpgsql 
AS $$
DECLARE
    approval_status VARCHAR(255);
BEGIN   
    SELECT verdict
    INTO approval_status
    FROM public.distributor_reviews dr
    WHERE dr.distributor_id = d_id
    AND dr.is_deleted = FALSE
    LIMIT 1;

    RETURN approval_status;
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

CREATE OR REPLACE FUNCTION public.get_retailer_by_user_id (
    r_user_id INT
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
    SELECT  r.id, 
            rb.name, 
            rb.tin, 
            rb_loc.lat, 
            rb_loc.long, 
            rb_loc.general_zone AS generalZone,
            rb_loc.region, 
            rb_loc.woreda
    FROM public.retailers r
    JOIN public.retailer_business_info rb 
        ON rb.retailer_id = r.id
    JOIN public.rb_locations rb_loc 
        ON rb_loc.business_id = rb.id
    JOIN public.retailer_users rb_rus
        ON rb_rus.retailer_id = r.id
    WHERE rb_rus.user_id = r_user_id
        AND r.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_retailer_by_name (
    r_retailer_name VARCHAR(255)
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

    SELECT r.id, rb.name, rb.tin, rb_loc.lat, rb_loc.long, 
           rb_loc.general_zone, rb_loc.region, rb_loc.woreda
    FROM public.retailers r
    JOIN public.retailer_business_info rb 
        ON rb.retailer_id = r.id
    JOIN public.rb_locations rb_loc 
        ON rb_loc.business_id = rb.id
    WHERE rb.name ILIKE '%' || r_retailer_name || '%'
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

create or replace function public.get_all_retailers (
    t_limit INT,
    t_offset INT
) 
RETURNS TABLE (
  id INT,
  name VARCHAR(255),
  tin VARCHAR(255),
  lat VARCHAR(255),
  long VARCHAR(255),
  generalZone VARCHAR(255),
  region VARCHAR(255),
  woreda VARCHAR(255),
  total_count BIGINT
) 
LANGUAGE plpgsql 
AS $$
    BEGIN
        RETURN QUERY

        SELECT r.id, rb.name, rb.tin, rb_loc.lat, rb_loc.long, 
        rb_loc.general_zone, rb_loc.region, rb_loc.woreda, COUNT(*) OVER() AS total_count
        FROM  public.retailers r
        JOIN public.retailer_business_info rb 
        ON rb.retailer_id = r.id
        JOIN public.rb_locations rb_loc 
        ON rb_loc.business_id = rb.id
        LIMIT t_limit
        OFFSET t_offset;
    END;
$$;

-- retailer user agent ---------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.create_retailer_user (
    r_id INT,
    r_user_id INT
) 
RETURNS INT 
LANGUAGE plpgsql 
AS $$
BEGIN   
    INSERT INTO public.retailer_users(user_id, retailer_id)
    VALUES(r_user_id, r_id);
    
    RETURN r_user_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.remove_user_connection (
    r_user_id INT
) 
RETURNS INT 
LANGUAGE plpgsql 
AS $$
BEGIN   
    UPDATE public.retailer_users
    SET is_deleted = TRUE
    WHERE user_id = r_user_id;

    UPDATE public.distributor_users
    SET is_deleted = TRUE
    WHERE user_id = r_user_id;

    UPDATE public.user_providers
    SET is_deleted = TRUE
    WHERE user_id = r_user_id;

    RETURN r_user_id;
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
    WHERE ru.retailer_id = r_id
    AND ru.is_deleted = FALSE;
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
RETURNS TABLE(
    id INT, 
    status VARCHAR(255), 
    external_id VARCHAR(255), 
    order_id INT, 
    subtotal DECIMAL(12,2), 
    tax_amount DECIMAL(12,2),
    created_date TIMESTAMP
    )
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, i.status, i.external_id, i.order_id, i.subtotal, i.tax_amount, i.created_date
    FROM public.invoices i
    WHERE i.id = i_invoice_id
      AND i.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_invoices(
    t_limit INT,
    t_offset INT
)
RETURNS TABLE(
    id INT, 
    status VARCHAR(255), 
    external_id VARCHAR(255), 
    order_id INT, 
    subtotal DECIMAL(12,2), 
    tax_amount DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT i.id, i.status, i.external_id, i.order_id, i.subtotal, i.tax_amount
    FROM public.invoices i
    WHERE i.is_deleted = FALSE
    LIMIT t_limit
    OFFSET t_offset;
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
    i_status VARCHAR(255),
    t_limit INT,
    t_offset INT
     
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
    AND i.status = i_status
    LIMIT t_limit
    OFFSET t_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_invoices_by_order_id(
    i_order_id INT
)
RETURNS TABLE(
    id INT, 
    status VARCHAR(255), 
    external_id VARCHAR(255), 
    order_id INT, 
    subtotal DECIMAL(12,2), 
    tax_amount DECIMAL(12,2),
    created_date TIMESTAMP  
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT 
        i.id, 
        i.status, 
        i.external_id, 
        i.order_id, 
        i.subtotal, 
        i.tax_amount,
        i.created_date 
    FROM public.invoices i
    WHERE i.is_deleted = FALSE 
    AND i.order_id = i_order_id;
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
    RETURNING id INTO new_id;  -- Assuming 'id' is the primary key column

    RETURN new_id;  -- Return the new ID
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
      AND i.is_deleted = FALSE;
END;
$$;

-- product -------------------------------------------------------------
    
    -- writers
CREATE OR REPLACE FUNCTION public.create_product(
  p_product_name VARCHAR(255),
  p_product_description VARCHAR(255),
  p_external_id VARCHAR(255),
  p_distributor_id INT,
  p_price DECIMAL(12,2)
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
                                 distributor_id,
                                 price)
    VALUES (p_product_name, 
            p_product_description, 
            p_external_id, 
            p_distributor_id,
            p_price) 
    RETURNING id INTO new_id;
    
    PERFORM public.create_product_stock(new_id); 

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


CREATE OR REPLACE FUNCTION public.update_product_price(
    i_product_id INT,
    i_price DECIMAL(12,2)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.products
    SET price = i_price
    WHERE id = i_product_id
      AND is_deleted = FALSE;

    RETURN i_product_id;
END;
$$;


    -- readers
CREATE OR REPLACE FUNCTION public.get_products_by_id(
    p_product_id INT
)
RETURNS TABLE(
    id INT, 
    product_name VARCHAR(255), 
    product_description VARCHAR(255), 
    external_id VARCHAR(255), 
    is_active BOOLEAN, 
    distributor_id INT,
    quantity INT,
    available_quantity INT,
    reserved_quantity INT,
    price DECIMAL(12,2)
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id,
           ps.quantity,
           (ps.quantity - ps.reserved_quantity) AS available_quantity,
           ps.reserved_quantity,
           p.price
    FROM public.products p
    JOIN p_stock ps 
      ON ps.product_id = p.id
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
              product_description VARCHAR(255), 
              external_id VARCHAR(255), 
              is_active BOOLEAN, 
              distributor_id INT,
              quantity INT,
              available_quantity INT,
              reserved_quantity INT,
              price DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id,
           ps.quantity,
           (ps.quantity - ps.reserved_quantity) AS available_quantity,
           ps.reserved_quantity,
           p.price
    FROM public.products p
    JOIN p_stock ps 
      ON  ps.product_id = p.id
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
              product_description VARCHAR(255), 
              external_id VARCHAR(255), 
              is_active BOOLEAN, 
              distributor_id INT,
              quantity INT,
              available_quantity INT,
              reserved_quantity INT,
              price DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id,
           ps.quantity,
           (ps.quantity - ps.reserved_quantity) AS available_quantity,
           ps.reserved_quantity,
           p.price
    FROM public.products p
     JOIN p_stock ps 
      ON  ps.product_id = p.id
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
              product_description VARCHAR(255), 
              external_id VARCHAR(255), 
              is_active BOOLEAN, 
              distributor_id INT,
              quantity INT,
              available_quantity INT,
              reserved_quantity INT,
              price DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
      SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id,
           ps.quantity,
           (ps.quantity - ps.reserved_quantity) AS available_quantity,
           ps.reserved_quantity,
           p.price
    FROM public.products p
     JOIN p_stock ps 
      ON  ps.product_id = p.id
    WHERE p.distributor_id = p_distributor_id
      AND p.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_products_paginated(
    p_limit INT DEFAULT 10,
    p_offset INT DEFAULT 0,
    p_search VARCHAR(255) DEFAULT NULL
)
RETURNS TABLE(
    id INT,
    product_name VARCHAR(255),
    product_description VARCHAR(255),
    external_id VARCHAR(255),
    is_active BOOLEAN,
    distributor_id INT,
    quantity INT,
    available_quantity INT,
    reserved_quantity INT,
    price DECIMAL(12,2),
    total_count BIGINT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.name,
        p.description,
        p.external_id,
        p.is_active,
        p.distributor_id,
        ps.quantity,
        (ps.quantity - ps.reserved_quantity) AS available_quantity,
        ps.reserved_quantity,
        p.price,
        COUNT(*) OVER() AS total_count
    FROM public.products p
    JOIN p_stock ps ON ps.product_id = p.id
    WHERE p.is_deleted = FALSE
        AND (p_search IS NULL OR 
             p.name ILIKE '%' || p_search || '%' OR 
             p.description ILIKE '%' || p_search || '%' OR
             p.external_id ILIKE '%' || p_search || '%')
    ORDER BY p.created_date DESC
    LIMIT p_limit
    OFFSET p_offset;
END
$$;

CREATE OR REPLACE FUNCTION public.get_all_products()
RETURNS TABLE(id INT, 
              product_name VARCHAR(255), 
              product_description VARCHAR(255), 
              external_id VARCHAR(255), 
              is_active BOOLEAN, 
              distributor_id INT,
              quantity INT,
              available_quantity INT,
              reserved_quantity INT,
              price DECIMAL(12,2))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
     SELECT p.id, 
           p.name, 
           p.description, 
           p.external_id, 
           p.is_active, 
           p.distributor_id,
           ps.quantity,
           (ps.quantity - ps.reserved_quantity) AS available_quantity,
           ps.reserved_quantity,
           p.price
    FROM public.products p
    JOIN p_stock ps 
    ON  ps.product_id = p.id
    WHERE p.is_deleted = FALSE;
   
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_price_range(
    p_min DECIMAL(12, 2),
    p_max DECIMAL(12, 2)
)
RETURNS TABLE(
    id INT, 
    product_name VARCHAR(255), 
    product_description VARCHAR(255), 
    external_id VARCHAR(255), 
    is_active BOOLEAN, 
    distributor_id INT,
    quantity INT,
    available_quantity INT,
    reserved_quantity INT,
    price DECIMAL(12, 2)
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
      SELECT p.id, 
             p.name, 
             p.description, 
             p.external_id, 
             p.is_active, 
             p.distributor_id,
             ps.quantity,
             (ps.quantity - ps.reserved_quantity) AS available_quantity,
             ps.reserved_quantity,
             p.price
      FROM public.products p
      JOIN p_stock ps 
        ON ps.product_id = p.id
      WHERE p.price BETWEEN p_min AND p_max
        AND p.is_deleted = FALSE
      LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_products_by_categoryIds(
   p_category_ids INT[]
)
RETURNS TABLE(
    id INT, 
    product_name VARCHAR(255), 
    product_description VARCHAR(255), 
    external_id VARCHAR(255), 
    is_active BOOLEAN, 
    distributor_id INT,
    quantity INT,
    available_quantity INT,
    reserved_quantity INT,
    price DECIMAL(12, 2)
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
      SELECT p.id, 
             p.name, 
             p.description, 
             p.external_id, 
             p.is_active, 
             p.distributor_id,
             ps.quantity,
             (ps.quantity - ps.reserved_quantity) AS available_quantity,
             ps.reserved_quantity,
             p.price
      FROM public.products p
      JOIN p_category pc 
        ON pc.product_id = p.id
      JOIN p_stock ps 
        ON ps.product_id = p.id
      WHERE pc.category_id = ANY(p_category_ids)
        AND p.is_deleted = FALSE
      LIMIT 1;
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
RETURNS TABLE(image_url VARCHAR(255), image_blur_hash VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT url,blur_hash
    FROM public.p_images i
    WHERE i.product_id = i_product_id
      AND i.is_deleted = FALSE;
END;
$$;

-- configurable product -------------------------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.create_configurable_product(
  p_product_name VARCHAR(255),
  p_product_description VARCHAR(255),
  p_external_id VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.configurable_products (name, 
                                              description, 
                                              external_id)
    VALUES (p_product_name, 
            p_product_description, 
            p_external_id) 
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_configurable_product_name(
    i_configurable_product_id INT,
    new_name VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.configurable_products
    SET name = new_name
    WHERE id = i_configurable_product_id
      AND is_deleted = FALSE;

    RETURN i_configurable_product_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_configurable_product_desc(
    i_configurable_product_id INT,
    new_desc VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.configurable_products
    SET description = new_desc
    WHERE id = i_configurable_product_id
      AND is_deleted = FALSE;

    RETURN i_configurable_product_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_configurable_product_externalId(
    i_configurable_product_id INT,
    new_external_id VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.configurable_products
    SET external_id = new_external_id
    WHERE id = i_configurable_product_id
      AND is_deleted = FALSE;

    RETURN i_configurable_product_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_configurable_product_isAvailable_status(
    i_configurable_product_id INT,
    is_available_status BOOLEAN
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.configurable_products
    SET is_available = is_available_status
    WHERE id = i_configurable_product_id
      AND is_deleted = FALSE;

    RETURN i_configurable_product_id;
END;
$$;
    
    -- reader
CREATE OR REPLACE FUNCTION public.get_configurable_products_by_id(
    cp_product_id INT
)
RETURNS TABLE(cp_id INT, 
              cp_name VARCHAR(255), 
              cp_description VARCHAR(255), 
              cp_external_id VARCHAR(255), 
              cp_is_available BOOLEAN)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT cp.id, 
           cp.name, 
           cp.description, 
           cp.external_id, 
           cp.is_available
    FROM public.configurable_products cp
    WHERE cp.id = cp_product_id
      AND cp.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_configurable_products()
RETURNS TABLE(id INT, 
              name VARCHAR(255), 
              description VARCHAR(255), 
              external_id VARCHAR(255), 
              is_available BOOLEAN)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT cp.id, 
           cp.name, 
           cp.description, 
           cp.external_id, 
           cp.is_available
    FROM public.configurable_products cp
    WHERE cp.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_configurable_products_by_name(
    cp_name VARCHAR(255) 
)
RETURNS TABLE(id INT, 
              name VARCHAR(255), 
              description VARCHAR(255), 
              external_id VARCHAR(255), 
              is_available BOOLEAN)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
     SELECT cp.id, 
           cp.name, 
           cp.description, 
           cp.external_id, 
           cp.is_available
    FROM public.configurable_products cp
    WHERE cp.name = cp_name
        AND cp.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_configurable_products_by_ext_id(
    cp_external_id VARCHAR(255) 
)
RETURNS TABLE(id INT, 
              name VARCHAR(255), 
              description VARCHAR(255), 
              external_id VARCHAR(255), 
              is_available BOOLEAN)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT cp.id, 
           cp.name, 
           cp.description, 
           cp.external_id, 
           cp.is_available
    FROM public.configurable_products cp
    WHERE cp.external_id = cp_external_id
        AND cp.is_deleted = FALSE;
END;
$$;

-- configurable product attribute ---------------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.add_attribute_to_configurable_product(
  p_attribute_id INT,
  cp_product_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.cp_attributes (product_attribute_id, configurable_product_id)
    VALUES (p_attribute_id, cp_product_id);
END;
$$;

CREATE OR REPLACE FUNCTION public.remove_all_configurable_product_attributes(
    cp_product_id INT
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
     UPDATE public.cp_attributes
    SET is_deleted = TRUE
    WHERE configurable_product_id = cp_product_id;

    RETURN cp_product_id;
END;
$$;
    -- reader
CREATE OR REPLACE FUNCTION public.get_all_configurable_product_attributes_values(
  cp_product_id INT
)
RETURNS TABLE(cp_attribute_id VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
     RETURN QUERY
    SELECT DISTINCT pa.name
    FROM public.cp_attributes c
    JOIN public.p_attributes pa
    ON pa.id = c.product_attribute_id
    WHERE c.configurable_product_id = cp_product_id
      AND c.is_deleted = FALSE;
END;
$$;

-- configurable product member ------------------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.add_product_to_configurable_product(
  cp_p_id INT,
  p_product_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.cp_members (cp_id, product_id)
    VALUES (cp_p_id, p_product_id);
END;
$$;

CREATE OR REPLACE FUNCTION public.remove_all_configurable_product_members(
    cp_product_id INT
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
     UPDATE public.cp_members
    SET is_deleted = TRUE
    WHERE cp_id = cp_product_id;

    RETURN cp_product_id;
END;
$$;
    -- reader
CREATE OR REPLACE FUNCTION public.get_all_configurable_product_members(
  cp_product_id INT
)
RETURNS TABLE(p_product_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
     RETURN QUERY
    SELECT c.product_id
    FROM public.cp_members c
    WHERE c.cp_id = cp_product_id
      AND c.is_deleted = FALSE;
END;
$$;

-- configurable product image ---------------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.add_image_to_configurable_product(
  i_url VARCHAR(255),
  i_blur_hash VARCHAR(255),
  i_configurable_product_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.cp_images (url, blur_hash, configurable_product_id)
    VALUES (i_url, i_blur_hash, i_configurable_product_id);
END;
$$;

CREATE OR REPLACE FUNCTION public.remove_all_configurable_product_images(
    i_configurable_product_id INT
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.cp_images
    SET is_deleted = TRUE
    WHERE configurable_product_id = i_configurable_product_id;

    RETURN i_configurable_product_id;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_images_by_cp_Id(
    i_configurable_product_id INT
)
RETURNS TABLE(image_url VARCHAR(255), image_blur_hash VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT url,blur_hash
    FROM public.cp_images i
    WHERE i.configurable_product_id = i_configurable_product_id
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


CREATE OR REPLACE FUNCTION public.update_category(
    c_id INT,
    c_name VARCHAR(255)
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.category
    SET name = c_name
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

-- -- product stock ----------------------------------------------------

    -- writer
CREATE OR REPLACE FUNCTION public.create_product_stock(
  i_product_id INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.p_stock (product_id)
    VALUES (i_product_id);
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


CREATE OR REPLACE FUNCTION public.reserve_product_stock(
    i_product_id INT,
    new_reserved_quantity INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.p_stock
    SET reserved_quantity = reserved_quantity + new_reserved_quantity
    WHERE product_id = i_product_id
        AND is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.free_reserved_product_stock(
    i_product_id INT,
    new_freed_quantity INT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.p_stock
    SET reserved_quantity = reserved_quantity - new_freed_quantity
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
CREATE OR REPLACE FUNCTION public.get_stock_ledger_entries()
RETURNS TABLE( 
  s_id INT,
  s_quantity INT,
  s_product_id INT,
  s_stock_operation VARCHAR(255),
  s_created_on TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT id,
           quantity,
           product_id,
           operation,
           created_on_date
    FROM public.s_ledger s;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_stock_ledger_entries_by_product_id(
    p_id INT
)
RETURNS TABLE( 
  s_id INT,
  s_quantity INT,
  s_product_id INT,
  s_stock_operation VARCHAR(255),
  s_created_on TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT id,
           quantity,
           product_id,
           operation,
           created_on_date
    FROM public.s_ledger s
    WHERE s.product_id = p_id;
END;
$$;

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
CREATE OR REPLACE FUNCTION public.get_attributes_values_by_productId(
    p_product_id INT
)
RETURNS TABLE(p_attribute_name VARCHAR(255), p_attribute_value VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.name, av.name 
    FROM public.p_attributes p
    JOIN public.p_attribute_values av
    ON av.attribute_id = p.id
    WHERE p.id = p_product_id
      AND p.is_deleted = FALSE;
END;
$$;


CREATE OR REPLACE FUNCTION public.get_attribute_id_by_name(
    p_attribute_name VARCHAR(255)
)
RETURNS TABLE(attribute_id INT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT pa.id
    FROM public.p_attributes pa
    WHERE pa.name = p_attribute_name
      AND pa.is_deleted = FALSE;
END;
$$;

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
CREATE OR REPLACE FUNCTION public.get_attributes_values_by_productId(
    p_product_id INT
)
RETURNS TABLE(p_attribute_name VARCHAR(255), p_attribute_value VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.name, av.name 
    FROM public.p_attribute_values av
    JOIN public.p_attributes p
    ON av.attribute_id = p.id
    WHERE av.product_id = p_product_id
      AND p.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_attributes_values_by_attribute_id(
    p_attribute_id INT
)
RETURNS TABLE(p_attribute_name VARCHAR(255), p_attribute_value VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.name, av.name 
    FROM public.p_attributes p
    RIGHT JOIN public.p_attribute_values av
    ON av.attribute_id = p.id
    WHERE av.attribute_id = p_attribute_id
      AND p.is_deleted = FALSE;
END;
$$;

--- Order ---------------------------------------------

    -- writer
CREATE OR REPLACE FUNCTION public.create_order(
  o_retailer_id INT,
  o_status VARCHAR(255),
  o_total DECIMAL(12,2),
  o_payment_status VARCHAR(255),
  o_delivery_status VARCHAR(255),
  o_confirmation_status VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.orders (retailer_id, 
                               status, 
                               total,
                               payment_status,
                               delivery_status,
                               confirmation_status)
    VALUES (o_retailer_id, 
            o_status, 
            o_total, 
            o_payment_status,
            o_delivery_status,
            o_confirmation_status)
    RETURNING id INTO new_id;
    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_order_status(
    order_id INT,
    new_status VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.orders
    SET status = new_status
    WHERE id = order_id
      AND is_deleted = FALSE;

    RETURN order_id;
END;
$$;


CREATE OR REPLACE FUNCTION public.update_order_payment_status(
    order_id INT,
    new_payment_status VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.orders
    SET payment_status = new_payment_status
    WHERE id = order_id
      AND is_deleted = FALSE;

    RETURN order_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_order_delivery_status(
    order_id INT,
    new_delivery_status VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.orders
    SET delivery_status = new_delivery_status
    WHERE id = order_id
      AND is_deleted = FALSE;

    RETURN order_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_order_confirmation_status(
    order_id INT,
    new_confirmation_status VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.orders
    SET confirmation_status = new_confirmation_status
    WHERE id = order_id
      AND is_deleted = FALSE;

    RETURN order_id;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_orders_by_id(
    o_order_id INT
)
RETURNS TABLE(id INT, 
              retailer_id INT,
              retailer_name VARCHAR(255), 
              status VARCHAR(255),
              total DECIMAL(12,2),
              delivery_status VARCHAR(255),
              payment_status VARCHAR(255),
              confirmation_status VARCHAR(255),
              created_date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT o.id, 
           o.retailer_id, 
           r.name,
           o.status, 
           o.total,
           o.payment_status,
           o.delivery_status,
           o.confirmation_status,
           o.created_date
    FROM public.orders o
    JOIN public.retailer_business_info r
    ON r.retailer_id = o.retailer_id
    WHERE o.id = o_order_id
      AND o.is_deleted = FALSE
    LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_orders_by_retailer_id(
    o_retailer_id INT,
    o_limit INT,
    o_offset INT
)
RETURNS TABLE(
    id INT, 
    retailer_id INT,
    retailer_name VARCHAR(255),
    status VARCHAR(255),
    total DECIMAL(12,2),
    delivery_status VARCHAR(255),
    payment_status VARCHAR(255),
    created_date TIMESTAMP,
    confirmation_status VARCHAR(255),
    payment_method VARCHAR(255),
    total_count BIGINT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT o.id, 
           o.retailer_id,
           r.name AS retailer_name,
           o.status, 
           o.total,
           o.delivery_status,
           o.payment_status,
           o.created_date,
           o.confirmation_status,
           pp.payment_method,
           COUNT(*) OVER() AS total_count
    FROM public.orders o
    JOIN public.retailer_business_info r
        ON r.retailer_id = o.retailer_id
    JOIN public.payments p
        ON p.order_id = o.id
    JOIN public.payment_partners pp
        ON pp.id = p.partner_id
    WHERE o.retailer_id = o_retailer_id
      AND o.is_deleted = FALSE
    LIMIT o_limit
    OFFSET o_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_orders_by_status(
    o_status VARCHAR(255),
    o_limit INT,
    o_offset INT
)
RETURNS TABLE(id INT, 
              retailer_id INT,
              retailer_name VARCHAR(255),
              status VARCHAR(255),
              total DECIMAL(12,2),
              payment_status VARCHAR(255),
              delivery_status VARCHAR(255),
              payment_method VARCHAR(255),
              created_date TIMESTAMP,
              total_count BIGINT,
              confirmation_status VARCHAR(255)
              )
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT o.id, 
        o.retailer_id,
        r.name,
        o.status, 
        o.total,
        o.payment_status,
        o.delivery_status,
        par.payment_method,
        o.created_date,
        COUNT(*) OVER() AS total_count,
        o.confirmation_status
    FROM public.orders o
    JOIN public.retailer_business_info r
    ON r.retailer_id = o.retailer_id
    JOIN public.payments p 
    ON p.order_id=o.id
    JOIN payment_partners par
    ON par.id = p.partner_id
    WHERE o.status = o_status
    AND o.is_deleted = FALSE
    LIMIT o_limit
    OFFSET o_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_orders(
    t_limit INT,
    t_offset INT
)
RETURNS TABLE(id INT, 
              retailer_id INT,
              retailer_name VARCHAR(255),
              status VARCHAR(255),
              total DECIMAL(12,2),
              payment_status VARCHAR(255),
              delivery_status VARCHAR(255),
              payment_method VARCHAR(255),
              created_date TIMESTAMP,
              total_count BIGINT,
              confirmation_status VARCHAR(255)
    )
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY

SELECT o.id, 
        o.retailer_id,
        r.name,
        o.status, 
        o.total,
        o.payment_status,
        o.delivery_status,
        par.payment_method,
        o.created_date,
        COUNT(*) OVER() AS total_count,
        o.confirmation_status
        FROM public.orders o
        JOIN public.retailer_business_info r
        ON r.retailer_id = o.retailer_id
        JOIN public.payments p 
        ON p.order_id=o.id
        JOIN payment_partners par
        ON par.id = p.partner_id
        WHERE o.is_deleted = FALSE
        ORDER BY o.created_date DESC
        LIMIT t_limit
        OFFSET t_offset;
END;
$$;

-- Order Item ------------------------------------------
    
    -- writer
CREATE OR REPLACE FUNCTION public.create_order_item(
  o_order_id INT,
  o_product_id INT,
  o_quantity INT,
  o_price DECIMAL(12,2)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO 
    public.o_items (order_id, 
                    product_id, 
                    quantity,
                    price)
    VALUES (o_order_id, 
            o_product_id, 
            o_quantity,
            o_price);
    
    RETURN o_order_id;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_order_items_by_order_id(
    oi_order_id INT
)
RETURNS TABLE(o_order_id INT,
              o_product_id INT,
              name VARCHAR(255),                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      
              o_quantity INT,
              o_price DECIMAL(12,2)
              )
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT oi.order_id, 
            oi.product_id, 
            p.name,
            oi.quantity, 
            oi.price
        FROM public.o_items oi
        JOIN public.products p
        ON p.id=oi.product_id
    WHERE oi.order_id = oi_order_id
      AND oi.is_deleted = FALSE;
END;
$$;

--- transaction -------------------

    -- writer
CREATE OR REPLACE FUNCTION public.record_transaction(
    t_amount DECIMAL(12,2),
    t_partner_id INT,
    t_status VARCHAR(255),
    t_tx_ref VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.transactions ( amount,
                                      partner_id,
                                      status,
                                      tx_ref)
    VALUES (t_amount,
            t_partner_id,
            t_status,
            t_tx_ref)
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$; 


CREATE OR REPLACE FUNCTION public.update_transaction_status(
    t_id INT,
    t_status VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    UPDATE public.transactions
    SET status = t_status
    WHERE id = t_id
      AND is_deleted = FALSE;

    RETURN t_id;
END;
$$; 

CREATE OR REPLACE FUNCTION public.update_transaction_by_transaction_ref(
    t_tx_ref VARCHAR(255),
    t_status VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    UPDATE public.transactions
    SET status = t_status
    WHERE tx_ref = t_tx_ref
      AND is_deleted = FALSE
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$; 

    -- reader
CREATE OR REPLACE FUNCTION public.get_transaction_by_id(
    t_id INT
)
RETURNS TABLE(id INT,
              amount DECIMAL(12,2),
              partner_id INT,
              tx_ref VARCHAR(255),
              status VARCHAR(255),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.amount, t.partner_id, t.tx_ref, t.status, t.date
    FROM public.transactions t
    WHERE t.id = t_id
      AND t.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_transactions(
    t_limit INT,
    t_offset INT
)
RETURNS TABLE(id INT,
              amount DECIMAL(12,2),
              partner_id INT,
              tx_ref VARCHAR(255),
              status VARCHAR(255),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.amount, t.partner_id, t.tx_ref, t.status, t.date
    FROM public.transactions t
    WHERE t.is_deleted = FALSE
    LIMIT t_limit
    OFFSET t_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_transactions_by_partner_id(
    t_partner_id INT,
    t_limit INT,
    t_offset INT
)
RETURNS TABLE(id INT,
              amount DECIMAL(12,2),
              partner_id INT,
              tx_ref VARCHAR(255),
              status VARCHAR(255),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.amount, t.partner_id, t.tx_ref, t.status, t.date
    FROM public.transactions t
    WHERE t.partner_id = t_partner_id
      AND t.is_deleted = FALSE
       LIMIT t_limit
    OFFSET t_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_transactions_by_date(
    t_date TIMESTAMP,
    t_limit INT,
    t_offset INT
)
RETURNS TABLE(id INT,
              amount DECIMAL(12,2),
              partner_id INT,
              tx_ref VARCHAR(255),
              status VARCHAR(255),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.amount, t.partner_id, t.tx_ref, t.status, t.date
    FROM public.transactions t
    WHERE t.date = t_date
      AND t.is_deleted = FALSE
       LIMIT t_limit
    OFFSET t_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_transactions_by_transaction_ref(
    transaction_ref VARCHAR(255),
    t_limit INT,
    t_offset INT
)
RETURNS TABLE(id INT,
              amount DECIMAL(12,2),
              partner_id INT,
              tx_ref VARCHAR(255),
              status VARCHAR(255),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.amount, t.partner_id, t.tx_ref, t.status, t.date
    FROM public.transactions t
    WHERE t.tx_ref = transaction_ref
      AND t.is_deleted = FALSE
       LIMIT t_limit
    OFFSET t_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_transactions_by_status(
    t_status TIMESTAMP,
    t_limit INT,
    t_offset INT
)
RETURNS TABLE(id INT,
              amount DECIMAL(12,2),
              partner_id INT,
              tx_ref VARCHAR(255),
              status VARCHAR(255),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.amount, t.partner_id, t.tx_ref, t.status, t.date
    FROM public.transactions t
    WHERE t.status = t_status
      AND t.is_deleted = FALSE
       LIMIT t_limit
    OFFSET t_offset;
END;
$$;

--- payment_partner --------------------------------------

    -- reader
CREATE OR REPLACE FUNCTION public.get_payment_partner_by_id(
    p_id INT
)
RETURNS TABLE(id int,
              name VARCHAR(255),
              icon VARCHAR(255),
              status VARCHAR(255),
              base_url VARCHAR(255),
              payment_method VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.icon, p.status, p.base_url, p.payment_method
    FROM public.payment_partners p
    WHERE p.id = p_id
      AND p.is_deleted = FALSE
      LIMIT 1;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_all_payment_partners(
p_limit INT,
p_offset INT
)
RETURNS TABLE(id int,
              name VARCHAR(255),
              icon VARCHAR(255),
              status VARCHAR(255),
              base_url VARCHAR(255),
              payment_method VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.icon, p.status, p.base_url, p.payment_method
    FROM public.payment_partners p
    WHERE p.is_deleted = FALSE
    LIMIT p_limit
    OFFSET p_offset;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_payment_partner_secret(
    p_id INT
)RETURNS TABLE(
              name VARCHAR(255),
              base_url VARCHAR(255),
              secret VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.name, p.base_url,p.secret
    FROM public.payment_partners p
    WHERE p.id = p_id
      AND p.is_deleted = FALSE;
END;
$$;


CREATE OR REPLACE FUNCTION public.get_payment_partner_by_name(
    p_name VARCHAR(255)
)
RETURNS TABLE(id int,
              name VARCHAR(255),
              icon VARCHAR(255),
              status VARCHAR(255),
              base_url VARCHAR(255),
              payment_method VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.icon, p.status, p.base_url, p.payment_method
    FROM public.payment_partners p
    WHERE p.name = p_name
      AND p.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_payment_partner_by_status(
    p_status VARCHAR(255)
)
RETURNS TABLE(id int,
              name VARCHAR(255),
              icon VARCHAR(255),
              status VARCHAR(255),
              base_url VARCHAR(255),
              payment_method VARCHAR(255))
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT p.id, p.name, p.icon, p.status, p.base_url, p.payment_method
    FROM public.payment_partners p
    WHERE p.status = p_status
      AND p.is_deleted = FALSE;
END;
$$;

    -- writer

CREATE OR REPLACE FUNCTION public.create_payment_partner(
    p_name VARCHAR(255),
    p_icon VARCHAR(255),
    p_status VARCHAR(255),
    p_base_url VARCHAR(255),
    p_secret VARCHAR(255),
    p_payment_method VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO 
    public.payment_partners (name, 
                            icon, 
                            status, 
                            base_url,
                            secret, 
                            payment_method)
    VALUES (p_name, 
            p_icon, 
            p_status,
            p_base_url,
            p_secret,
            p_payment_method
            )
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;  


CREATE OR REPLACE FUNCTION public.update_payment_partner_status(
    p_id INT,
    new_status VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.payment_partners
    SET status = new_status
    WHERE id = p_id
      AND is_deleted = FALSE;

    RETURN p_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_payment_partner_name(
    p_id INT,
    new_name VARCHAR(255)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.payment_partners
    SET name = new_name
    WHERE id = p_id
      AND is_deleted = FALSE;

    RETURN p_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_payment_partner_secret(
    p_id INT,
    p_base_url VARCHAR(255),
    p_secret VARCHAR(255)
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.payment_partners
    SET base_url = p_base_url,
        secret = p_secret
    WHERE id = p_id
      AND is_deleted = FALSE;
END;
$$;

-- payment ----------------------------------------

    -- writer

CREATE OR REPLACE FUNCTION public.create_payment(
    p_order_id INT,
    p_partner_id INT,
    p_transaction_ref VARCHAR(255),
    p_amount DECIMAL(12, 2)
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.payments(order_id, 
                                partner_id, 
                                transaction_ref, 
                                amount)
    VALUES (p_order_id, 
            p_partner_id, 
            p_transaction_ref, 
            p_amount)
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

    -- reader
CREATE OR REPLACE FUNCTION public.get_payment_by_id(
    p_id INT
)
RETURNS TABLE(id INT,
              order_id INT,
              partner_id INT,
              transaction_ref VARCHAR(255),
              amount DECIMAL(12, 2),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, 
           t.order_id, 
           t.partner_id,
           t.transaction_ref,
           t.amount,
           t.created_date
    FROM public.payments t
    WHERE t.id = p_id
      AND t.is_deleted = FALSE;
END;
$$;


CREATE OR REPLACE FUNCTION public.get_payment_by_order_id(
    p_order_id INT
)
RETURNS TABLE(id INT,
              order_id INT,
              partner_id INT,
              transaction_ref VARCHAR(255),
              amount DECIMAL(12, 2),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, 
           t.order_id, 
           t.partner_id,
           t.transaction_ref,
           t.amount,
           t.created_date
    FROM public.payments t
    WHERE t.order_id = p_order_id
      AND t.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_payment_by_tx_ref(
    p_tx_ref VARCHAR(255)
)
RETURNS TABLE(id INT,
              order_id INT,
              partner_id INT,
              transaction_ref VARCHAR(255),
              amount DECIMAL(12, 2),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, 
           t.order_id, 
           t.partner_id,
           t.transaction_ref,
           t.amount,
           t.created_date
    FROM public.payments t
    WHERE t.transaction_ref = p_tx_ref
      AND t.is_deleted = FALSE;
END;
$$; 

CREATE OR REPLACE FUNCTION public.get_all_payment()
RETURNS TABLE(id INT,
              order_id INT,
              partner_id INT,
              transaction_ref VARCHAR(255),
              amount DECIMAL(12, 2),
              date TIMESTAMP)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, 
           t.order_id, 
           t.partner_id,
           t.transaction_ref,
           t.amount,
           t.created_date
    FROM public.payments t
    WHERE t.is_deleted = FALSE;
END;
$$; 

-- email_templates ---------------------------

    -- Write

CREATE OR REPLACE FUNCTION public.create_email_template(
    template_name VARCHAR(255),
    template_html TEXT
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    new_id INT;
BEGIN
    INSERT INTO public.email_templates (name, html)
    VALUES (template_name, template_html)
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

CREATE OR REPLACE FUNCTION public.update_template(
   id INT,
   name VARCHAR(255),
   new_html TEXT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE public.email_templates
    SET html = new_html,
        last_modified = CURRENT_TIMESTAMP
    WHERE name = template_name;
    COMMIT;
END;
$$;

    -- Read
CREATE OR REPLACE FUNCTION public.get_all_templates()
RETURNS TABLE(id int,
              name VARCHAR(255),
              html TEXT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, 
           t.name, 
           t.html
    FROM public.email_templates t
    WHERE t.is_deleted = FALSE;
END;
$$;

CREATE OR REPLACE FUNCTION public.get_template_by_id(
   template_id INT
)
RETURNS TABLE(id int,
              name VARCHAR(255),
              html TEXT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, 
           t.name, 
           t.html
    FROM public.email_templates t
    WHERE t.id = template_id
      AND t.is_deleted = FALSE;
END;
$$;


CREATE OR REPLACE FUNCTION public.get_template_by_name(
   template_name VARCHAR(255)
)
RETURNS TABLE(id int,
              name VARCHAR(255),
              html TEXT)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, 
           t.name, 
           t.html
    FROM public.email_templates t
    WHERE t.name = template_name
      AND t.is_deleted = FALSE;
END;
$$;


-- order_expiry_configuration ---------------
    
    -- Read

CREATE OR REPLACE FUNCTION public.get_order_expiry_duration_config()
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    expiry_duration INT;
BEGIN
    SELECT duration_in_minutes INTO expiry_duration
    FROM public."order_expiry_duration_config"
    LIMIT 1;

    RETURN COALESCE(expiry_duration, 0);
END;
$$;
    
    -- Write

CREATE OR REPLACE FUNCTION public.set_order_expiry_duration_config(new_duration INT)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM public."order_expiry_duration_config") THEN
        UPDATE public."order_expiry_duration_config"
        SET duration_in_minutes = new_duration;
    ELSE
        INSERT INTO public."order_expiry_duration_config" (duration_in_minutes)
        VALUES (new_duration);
    END IF;
END;
$$;