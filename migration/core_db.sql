-- version 0.1

CREATE TABLE IF NOT EXISTS public."base"
(
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);
COMMENT ON TABLE public."base" IS 'stores universal fields.';

CREATE TABLE IF NOT EXISTS public."roles" 
(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  description VARCHAR(255)
) INHERITS (public."base");

COMMENT ON TABLE public."roles" IS 'stores role definition of user agent.';

CREATE TABLE IF NOT EXISTS public."users" 
(
  id SERIAL PRIMARY KEY,
  firstName VARCHAR(255) NOT NULL,
  lastName VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  phone_number VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  birth_date DATE,
  is_active BOOLEAN DEFAULT FALSE,
  external_id VARCHAR(255)
) INHERITS (public."base");

COMMENT ON TABLE public."users" IS 'stores agent that interacts with the application.';

CREATE TABLE IF NOT EXISTS public."user_roles"
(
	user_id INT,
	role_id INT,
	FOREIGN KEY (role_id) REFERENCES public."roles"(id) ON DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES public."users"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."user_roles" IS 'stores role user mappings';


CREATE TABLE IF NOT EXISTS public."user_providers"
(
    user_id INT,
    provider_id VARCHAR(255),
    PRIMARY KEY(user_id, provider_id),
    FOREIGN KEY (user_id) REFERENCES public."users"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."user_providers" IS 'stores user provider mappings';

CREATE TABLE IF NOT EXISTS public."resources" 
(
  id SERIAL PRIMARY KEY,
  action VARCHAR(255),
  name VARCHAR(255)
) INHERITS (public."base");

COMMENT ON TABLE public."resources" IS 'stores permissible resources for user agent role.';

CREATE TABLE IF NOT EXISTS public."role_resources"
(
	role_id INT,
	resource_id INT,
	FOREIGN KEY (role_id) REFERENCES public."roles"(id) ON DELETE CASCADE,
	FOREIGN KEY(resource_id) REFERENCES public."resources"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."role_resources" IS 'stores role resource mappings';


CREATE TABLE IF NOT EXISTS public."retailers" 
(
  id SERIAL PRIMARY KEY
) INHERITS (public."base");

COMMENT ON TABLE public."retailers" IS 'stores retailer specific information(not user)';

CREATE TABLE IF NOT EXISTS public."retailer_users"
(
  user_id INT,
  retailer_id INT,
	FOREIGN KEY (retailer_id) REFERENCES public."retailers"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."retailer_users" IS 'stores retailer agents(always on the consumtion end of the application';

CREATE TABLE IF NOT EXISTS public."distributors"
(
  id SERIAL PRIMARY KEY
) INHERITS (public."base");

COMMENT ON TABLE public."distributors" IS 'stores distributor specific information(not user)';

CREATE TABLE IF NOT EXISTS public."distributor_users"
(
  user_id INT,
  distributor_id INT
	-- FOREIGN KEY (user_id) REFERENCES public."users"(id) ON DELETE CASCADE,
	-- FOREIGN KEY (distributor_id) REFERENCES public."distributors"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."distributor_users" IS 'stores distributor agents(always on the supply end of the application';

CREATE TABLE IF NOT EXISTS public."admins"
(
	id SERIAL PRIMARY KEY
) INHERITS (public."base");

COMMENT ON TABLE public."admins" IS 'stores admin specific information(not user)';

CREATE TABLE IF NOT EXISTS public."admin_users"
(
  user_id INT,
  admin_id INT,
	FOREIGN KEY (user_id) REFERENCES public."users"(id) ON DELETE CASCADE,
	FOREIGN KEY (admin_id) REFERENCES public."admins"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."admin_users" IS 'stores admin agents(always on the adminstration end of the application';

CREATE TABLE IF NOT EXISTS public."retailer_business_info" 
(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  tin VARCHAR(10) NOT NULL,
  retailer_id INT UNIQUE REFERENCES public."retailers" (id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."retailer_business_info" IS 'stores busines information of retailers';

CREATE TABLE IF NOT EXISTS public."rb_locations"
(	
  lat VARCHAR(255),
  long VARCHAR(255),
  general_zone VARCHAR(255) NOT NULL,
  region VARCHAR(255) NOT NULL,
  woreda VARCHAR(255) NOT NULL,
  business_id INT,
  FOREIGN KEY(business_id) REFERENCES public."retailer_business_info" (id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."rb_locations" IS 'stores location information of businesses';

CREATE TABLE IF NOT EXISTS public."distributor_business_info" 
(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  tin VARCHAR(10) NOT NULL,
  distributor_id INT 
  -- UNIQUE REFERENCES public."distributors" (id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."distributor_business_info" IS 'stores busines information of distributors';

CREATE TABLE IF NOT EXISTS public."db_locations" 
(
  lat VARCHAR(255),
  long VARCHAR(255),
  general_zone VARCHAR(255) NOT NULL,
  region VARCHAR(255) NOT NULL,
  woreda VARCHAR(255) NOT NULL,
  business_id INT,
  FOREIGN KEY(business_id) REFERENCES public."distributor_business_info" (id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."rb_locations" IS 'stores location information of distributors';

CREATE TABLE IF NOT EXISTS public."category" 
(
  id SERIAL PRIMARY KEY, 
  name VARCHAR(255)
) INHERITS (public."base");

COMMENT ON TABLE public."category" IS 'category tags for products';

CREATE TABLE IF NOT EXISTS public."products" 
(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  description VARCHAR(255),
  external_id VARCHAR(255),
  is_active BOOLEAN DEFAULT FALSE,
  distributor_id INT
  -- ,
  -- FOREIGN KEY (distributor_id) REFERENCES public."distributors" (id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."products" IS 'stores products';

CREATE TABLE IF NOT EXISTS public."p_category" 
(
  product_id INT,
  category_id INT
) INHERITS (public."base");

COMMENT ON TABLE public."p_category" IS 'many to many relationship for category and product';

CREATE TABLE IF NOT EXISTS public."p_stock"
(
  quantity INT,
  product_id INT,
   FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."p_stock" IS 'quantiative data about products';

CREATE TABLE IF NOT EXISTS public."s_ledger"
(
  id SERIAL PRIMARY KEY,
  quantity INT,
  product_id INT,
  operation VARCHAR(255),
  created_by_user_id INT,
  created_on_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE
);

COMMENT ON TABLE public."s_ledger" IS 'ledger for stock movement';

CREATE TABLE IF NOT EXISTS public."p_prices" 
(
  price DECIMAL(12,2),
  product_id INT,
  FOREIGN KEY (product_id) REFERENCES public."products" (id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."p_prices" IS 'stores product price information';

CREATE TABLE IF NOT EXISTS public."p_attributes" 
(
  id SERIAL PRIMARY KEY, 
  name VARCHAR(255)
) INHERITS (public."base");

COMMENT ON TABLE public."p_attributes" IS 'stores product attributes(part of EAV)';

CREATE TABLE IF NOT EXISTS public."p_attribute_values"
(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  product_id INT,
  attribute_id INT,
	FOREIGN KEY (attribute_id) REFERENCES public."p_attributes"(id) ON DELETE CASCADE,
	FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."p_attribute_values" IS 'stores product attributes values(part of EAV)';

CREATE TABLE IF NOT EXISTS public."p_images"
(
  url VARCHAR(255),
  blur_hash VARCHAR(255),
  product_id INT,
    FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."p_images" IS 'stores images of products';

CREATE TABLE IF NOT EXISTS public."configurable_products"
(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  description VARCHAR(255),
  external_id VARCHAR(255),
  is_available BOOLEAN DEFAULT FALSE
) INHERITS(public."base");

COMMENT ON TABLE public."configurable_products" IS 'meta-product definition';

CREATE TABLE IF NOT EXISTS public."cp_attributes"
(
  id SERIAL PRIMARY KEY,
  product_attribute_id INT,
  configurable_product_id INT,
	FOREIGN KEY (product_attribute_id) REFERENCES public."p_attributes"(id) ON DELETE CASCADE,
	FOREIGN KEY (configurable_product_id) REFERENCES public."configurable_products"(id) ON DELETE CASCADE
	
) INHERITS(public."base");

COMMENT ON TABLE public."cp_attributes" IS 'meta-product definition criteria(part of EAV)';

CREATE TABLE IF NOT EXISTS public."cp_members"
(
  id SERIAL PRIMARY KEY,
  cp_id INT,
  product_id INT,
	FOREIGN KEY (cp_id) REFERENCES public."configurable_products"(id) ON DELETE CASCADE,
	FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."cp_members" IS 'configurable product attribute values(part of EAV)';

CREATE TABLE IF NOT EXISTS public."cp_images"
(
  url VARCHAR(255),
  blur_hash VARCHAR(255),
  configurable_product_id INT,
  FOREIGN KEY (configurable_product_id) REFERENCES public."configurable_products"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."p_attributes" IS 'stores images of products';

CREATE TABLE IF NOT EXISTS public."orders"
(
  id SERIAL PRIMARY KEY,
  retailer_id INT,
  status VARCHAR(255),
  payment_status VARCHAR(255),
  delivery_status VARCHAR(255),
  total DECIMAL(12,2),
	FOREIGN KEY (retailer_id) REFERENCES public."retailers"(id) ON DELETE CASCADE

) INHERITS(public."base");

COMMENT ON TABLE public."orders" IS 'stores orders';

CREATE TABLE IF NOT EXISTS public."o_items" 
(
  order_id INT,
  product_id INT,
  quantity INT,
  price DECIMAL(12,2),
  FOREIGN KEY (order_id) REFERENCES public."orders" (id) ON DELETE CASCADE,
  FOREIGN KEY (product_id) REFERENCES public."products" (id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."o_items" IS 'stores order items';


CREATE TABLE IF NOT EXISTS public."invoices" 
(
  id SERIAL PRIMARY KEY,
  status VARCHAR(255),
  external_id VARCHAR(255),
  order_id INT,
  subtotal DECIMAL(12,2),
  tax_amount DECIMAL(12,2)
) INHERITS (public."base");

COMMENT ON TABLE public."invoices" IS 'stores invoices';


CREATE TABLE IF NOT EXISTS public."invoice_line_items" 
(
  id SERIAL PRIMARY KEY,
  product_name VARCHAR(255),
  qty INT,
  price DECIMAL(12,2),
  product_id INT,
  invoice_id INT,
  FOREIGN KEY (invoice_id) REFERENCES public."invoices" (id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."invoices" IS 'stores invoice line items';

CREATE TABLE IF NOT EXISTS public."transactions"
(
  id SERIAL PRIMARY KEY,
  user_id INT,
  amount DECIMAL(12,2),
  partner_id INT,
  date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_Id) REFERENCES public."users" (id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."invoices" IS 'stores transactions';

CREATE TABLE IF NOT EXISTS public."payment_partners"
(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  icon VARCHAR(255),
  status VARCHAR(255),
  baseurl VARCHAR(255),
  secret TEXT
) INHERITS (public."base");

COMMENT ON TABLE public."invoices" IS 'stores payment processing partners';