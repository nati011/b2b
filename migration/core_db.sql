
CREATE TABLE IF NOT EXISTS public."base"
(
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS public."user"
(
    id SERIAL PRIMARY KEY,
    fullName VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    birth_date DATE,
    is_active BOOLEAN DEFAULT FALSE,
    external_id VARCHAR(255)
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."role"
(
    id SERIAL PRIMARY KEY,
    r_name VARCHAR(255),
    r_desc VARCHAR(255),
    user_id INT,
	FOREIGN KEY (user_id) REFERENCES public."user"(id)
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."resource"
(
    id SERIAL PRIMARY KEY,
    r_action VARCHAR(4),
    r_name VARCHAR(255),
    role_id INT,
    FOREIGN KEY (role_id) REFERENCES public."role"(id)
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."role_resource"
(
    role_id INT,
    resource_id INT,
    FOREIGN KEY (role_id) REFERENCES public."role"(id),
    FOREIGN KEY (resource_id) REFERENCES public."resource"(id)
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."retailer"
(
  user_id INT,
  FOREIGN KEY (user_id) REFERENCES public."user"(id)
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."distributor"
(
  id SERIAL PRIMARY KEY
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."admin"
(
	id SERIAL PRIMARY KEY
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."distributor_user"
(
	user_id INT,
	distributor_id INT,
	FOREIGN KEY (user_id) REFERENCES public."user"(id),
	FOREIGN KEY (distributor_id) REFERENCES public."admin"(id)
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."store"
(
 id SERIAL PRIMARY KEY,
 distributor_id INT,
 FOREIGN KEY (distributor_id) REFERENCES public."distributor"(id)
) INHERITS (public."base");


CREATE TABLE IF NOT EXISTS public."admin_user"
(
	user_id INT,
	admin_id INT,
	FOREIGN KEY (user_id) REFERENCES public."user"(id),
	FOREIGN KEY (admin_id) REFERENCES public."admin"(id)
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."business"
(
	id SERIAL PRIMARY KEY,
    name TEXT,
	tin VARCHAR(10) NOT NULL,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES public."user"(id)
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."b_location"
(	
	lat FLOAT,
    long FLOAT,
	general_zone VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    woreda VARCHAR(255) NOT NULL,
	business_id INT,
    FOREIGN KEY (business_id) REFERENCES public."business"(id)
) INHERITS (public."base");

CREATE TABLE IF NOT EXISTS public."product"
(
    id SERIAL PRIMARY KEY,
	name VARCHAR(255),
    description TEXT,
    external_id VARCHAR(255),
	is_active BOOLEAN
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."p_stock"
(
   quantity INT,
   product_id INT,
   FOREIGN KEY (product_id) REFERENCES public."product"(id)
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."p_price"
(
   price MONEY,
   product_id INT,
   FOREIGN KEY (product_id) REFERENCES public."product"(id)
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."p_attribute"
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."p_attribute_values"
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
	product_id INT,
	attribute_id INT,
	FOREIGN KEY (attribute_id) REFERENCES public."p_attribute"(id),
	FOREIGN KEY (product_id) REFERENCES public."product"(id)
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."p_image"
(
    url VARCHAR(255),
    blur_hash_color_code VARCHAR(255),
    product_id INT,
    FOREIGN KEY (product_id) REFERENCES public."product"(id)
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."configurable_product"
(
    id SERIAL PRIMARY KEY,
	name VARCHAR(255),
    description TEXT,
    external_id VARCHAR(255),
	is_available BOOLEAN
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."cp_attribute"
(
    id SERIAL PRIMARY KEY,
    p_attribute_id INT,
	FOREIGN KEY (p_attribute_id) REFERENCES public."p_attribute"(id)
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."cp_member"
(
    id SERIAL PRIMARY KEY,
    cp_attribute_id INT,
	product_id INT,
	FOREIGN KEY (cp_attribute_id) REFERENCES public."cp_attribute"(id),
	FOREIGN KEY (product_id) REFERENCES public."product"(id)
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."o_status"
(
    id SERIAL PRIMARY KEY,
	name VARCHAR(255),
	value VARCHAR(255)
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."order"
(
    id SERIAL PRIMARY KEY,
	retailer_id INT,
	product_id INT,
	status_id INT,
	FOREIGN KEY (retailer_id) REFERENCES public."user"(id),
	FOREIGN KEY (product_id) REFERENCES public."product"(id),
	FOREIGN KEY (status_id) REFERENCES public."o_status"(id)
) INHERITS(public."base");

CREATE TABLE IF NOT EXISTS public."o_item"
(
	order_id INT,
	product_id INT,
	distributor_id INT,
	FOREIGN KEY (order_id) REFERENCES public."order"(id),
	FOREIGN KEY (product_id) REFERENCES public."product"(id),
	FOREIGN KEY (distributor_id) REFERENCES public."user"(id)
) INHERITS(public."base");

