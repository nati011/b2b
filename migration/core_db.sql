CREATE TABLE IF NOT EXISTS public."base"
(
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);
COMMENT ON TABLE public."base" IS 'stores universal fields.';

CREATE TABLE IF NOT EXISTS public."users"
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

COMMENT ON TABLE public."users" IS 'stores agent that interacts with the application.';

CREATE TABLE IF NOT EXISTS public."u_operations"
(
	id SERIAL PRIMARY KEY,
	value VARCHAR(255)
) INHERITS(public."base");

COMMENT ON TABLE public."u_operations" IS 'operation that can be performed on users';

CREATE TABLE public."user_audit"
(
    id SERIAL PRIMARY KEY,
    user_operation_id INT,
    new_record_id INT,
    old_record_id INT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	changed_by INT REFERENCES public."users"(id),
	FOREIGN KEY (user_operation_id) REFERENCES public."u_operations"(id) ON DELETE CASCADE,
	FOREIGN KEY (old_record_id) REFERENCES public."users"(id) ON DELETE CASCADE,
	FOREIGN KEY (new_record_id) REFERENCES public."users"(id) ON DELETE CASCADE
);

COMMENT ON TABLE public."user_audit" IS 'audit trail for users table';

CREATE TABLE IF NOT EXISTS public."r_operations"
(
	id SERIAL PRIMARY KEY,
	value VARCHAR(255)
) INHERITS(public."base");

COMMENT ON TABLE public."r_operations" IS 'operation that can be performed on roles';

CREATE TABLE IF NOT EXISTS public."roles"
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    description VARCHAR(255),
    user_id INT,
	FOREIGN KEY (user_id) REFERENCES public."users"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."roles" IS 'stores role definition of user agent.';

CREATE TABLE public."role_audit"
(
    id SERIAL PRIMARY KEY,
    role_operation_id INT,
    new_record_id INT,
    old_record_id INT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	changed_by INT REFERENCES public."users"(id),
	FOREIGN KEY (role_operation_id) REFERENCES public."r_operations"(id) ON DELETE CASCADE,
	FOREIGN KEY (old_record_id) REFERENCES public."roles"(id) ON DELETE CASCADE,
	FOREIGN KEY (new_record_id) REFERENCES public."roles"(id) ON DELETE CASCADE
);

COMMENT ON TABLE public."role_audit" IS 'audit trail for roles table';

CREATE TABLE IF NOT EXISTS public."resources"
(
    id SERIAL PRIMARY KEY,
    action VARCHAR(4),
    name VARCHAR(255),
    role_id INT,
    FOREIGN KEY (role_id) REFERENCES public."roles"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."resources" IS 'stores permissible resources for user agent role.';

CREATE TABLE IF NOT EXISTS public."role_resources"
(
    id SERIAL PRIMARY KEY,
	role_id INT,
    resource_id INT,
    FOREIGN KEY (role_id) REFERENCES public."roles"(id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES public."resources"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."role_resources" IS 'stores which resoures roles have permission to';

CREATE TABLE IF NOT EXISTS public."role_resources_operations"
(
	id SERIAL PRIMARY KEY,
	value VARCHAR(255)
) INHERITS(public."base");

COMMENT ON TABLE public."role_resources_operations" IS 'operation that can be performed on role_resources';

CREATE TABLE public."role_resources_audit"
(
    id SERIAL PRIMARY KEY,
    role_resources_operation_id INT NOT NULL,
    new_record_id INT,
    old_record_id INT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	changed_by INT REFERENCES public."users"(id),
	FOREIGN KEY (role_resources_operation_id) REFERENCES public."role_resources_operations"(id) ON DELETE CASCADE,
	FOREIGN KEY (old_record_id) REFERENCES public."role_resources"(id) ON DELETE CASCADE,
	FOREIGN KEY (new_record_id) REFERENCES public."role_resources"(id) ON DELETE CASCADE
);

COMMENT ON TABLE public."role_resources_audit" IS 'audit trail for role_resources table';

CREATE TABLE IF NOT EXISTS public."retailers"
(
  id SERIAL PRIMARY KEY
) INHERITS (public."base");

COMMENT ON TABLE public."retailers" IS 'stores retailer specific information(not user)';

CREATE TABLE IF NOT EXISTS public."retailer_users"
(
	user_id INT,
	retailer_id INT,
	FOREIGN KEY (user_id) REFERENCES public."users"(id) ON DELETE CASCADE,
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
	distributor_id INT,
	FOREIGN KEY (user_id) REFERENCES public."users"(id) ON DELETE CASCADE,
	FOREIGN KEY (distributor_id) REFERENCES public."distributors"(id) ON DELETE CASCADE
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

--------------------------------------------- revisit

CREATE TABLE IF NOT EXISTS public."retailer_business_info"
(
	id SERIAL PRIMARY KEY,
    name TEXT,
	tin VARCHAR(10) NOT NULL,
    retailer_id INT,
    FOREIGN KEY (retailer_id) REFERENCES public."retailers"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."retailer_business_info" IS 'stores busines information of retailers';

CREATE TABLE IF NOT EXISTS public."rb_locations"
(	
	lat FLOAT,
    long FLOAT,
	general_zone VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    woreda VARCHAR(255) NOT NULL,
	business_id INT,
    FOREIGN KEY (business_id) REFERENCES public."retailer_business_info"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."rb_locations" IS 'stores location information of businesses';

CREATE TABLE IF NOT EXISTS public."distributor_business_info"
(
	id SERIAL PRIMARY KEY,
    name TEXT,
	tin VARCHAR(10) NOT NULL,
    distributor_id INT,
    FOREIGN KEY (distributor_id) REFERENCES public."distributors"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."distributor_business_info" IS 'stores busines information of distributors';

CREATE TABLE IF NOT EXISTS public."db_locations"
(	
	general_zone VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    woreda VARCHAR(255) NOT NULL,
	business_id INT,
    FOREIGN KEY (business_id) REFERENCES public."distributor_business_info"(id) ON DELETE CASCADE
) INHERITS (public."base");

COMMENT ON TABLE public."rb_locations" IS 'stores location information of distributors';

--------------------------------------------- revisit

CREATE TABLE IF NOT EXISTS public."products"
(
    id SERIAL PRIMARY KEY,
	name VARCHAR(255),
    description TEXT,
    external_id VARCHAR(255),
	is_active BOOLEAN,
	distributor_id INT,
	FOREIGN KEY (distributor_id) REFERENCES public."distributors"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."products" IS 'stores products';

CREATE TABLE IF NOT EXISTS public."product_operations"
(
	id SERIAL PRIMARY KEY,
	value VARCHAR(255)
) INHERITS(public."base");

COMMENT ON TABLE public."product_operations" IS 'operation that can be performed on role_resources';

CREATE TABLE public."product_audit"
(
    id SERIAL PRIMARY KEY,
    product_operation_id INT,
    new_record_id INT,
    old_record_id INT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	changed_by INT REFERENCES public."users"(id),
	FOREIGN KEY (product_operation_id) REFERENCES public."product_operations"(id) ON DELETE CASCADE,
	FOREIGN KEY (old_record_id) REFERENCES public."products"(id) ON DELETE CASCADE,
	FOREIGN KEY (new_record_id) REFERENCES public."products"(id) ON DELETE CASCADE
);

COMMENT ON TABLE public."product_audit" IS 'audit trail for product table';

CREATE TABLE IF NOT EXISTS public."p_category"
(
   id SERIAL PRIMARY KEY,
   name VARCHAR(255),
   product_id INT,
   FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."p_category" IS 'category tags for products';

CREATE TABLE IF NOT EXISTS public."p_stock"
(
   quantity INT,
   product_id INT,
   FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."p_stock" IS 'quantiative data about products';

CREATE TABLE IF NOT EXISTS public."s_operations"
(
	id SERIAL PRIMARY KEY,
	value VARCHAR(255)
) INHERITS(public."base");

COMMENT ON TABLE public."s_operations" IS 'operation that can be performed on stocks';

CREATE TABLE IF NOT EXISTS public."s_ledger"
(
   id SERIAL PRIMARY KEY,
   quantity INT,
   product_id INT,
   stock_operation_id INT,
   created_by_user_id INT,
   created_on_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE,
   FOREIGN KEY (stock_operation_id) REFERENCES public."s_operations"(id) ON DELETE CASCADE,
   FOREIGN KEY (created_by_user_id) REFERENCES public."users"(id) ON DELETE CASCADE
);

COMMENT ON TABLE public."s_ledger" IS 'ledger for stock movement';

CREATE TABLE IF NOT EXISTS public."p_prices"
(
   price MONEY,
   product_id INT,
   FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."p_prices" IS 'stores product price information';

CREATE TABLE IF NOT EXISTS public."p_attributes"
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
) INHERITS(public."base");

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

COMMENT ON TABLE public."p_attributes" IS 'stores product attributes values(part of EAV)';

CREATE TABLE IF NOT EXISTS public."p_images"
(
    url VARCHAR(255),
    blur_hash_color_code VARCHAR(255),
    product_id INT,
    FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."p_attributes" IS 'stores images of products';

CREATE TABLE IF NOT EXISTS public."configurable_products"
(
    id SERIAL PRIMARY KEY,
	name VARCHAR(255),
    description TEXT,
    external_id VARCHAR(255),
	is_available BOOLEAN
) INHERITS(public."base");

COMMENT ON TABLE public."configurable_products" IS 'meta-product definition';

CREATE TABLE IF NOT EXISTS public."cp_attributes"
(
    id SERIAL PRIMARY KEY,
    product_attribute_id INT,
	configurable_product_id INT,
	FOREIGN KEY (product_attribute_id) REFERENCES public."p_attributes"(id) ON DELETE CASCADE,
	FOREIGN KEY (configurable_product_id) REFERENCES public."p_attributes"(id) ON DELETE CASCADE
	
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

CREATE TABLE IF NOT EXISTS public."o_statuses"
(
    id SERIAL PRIMARY KEY,
	value VARCHAR(255),
	description TEXT
) INHERITS(public."base");

COMMENT ON TABLE public."o_statuses" IS 'order states';

CREATE TABLE IF NOT EXISTS public."orders"
(
    id SERIAL PRIMARY KEY,
	retailer_id INT,
	status_id INT,
	total MONEY,
	FOREIGN KEY (retailer_id) REFERENCES public."users"(id) ON DELETE CASCADE,
	FOREIGN KEY (status_id) REFERENCES public."o_statuses"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."orders" IS 'stores orders';

CREATE TABLE IF NOT EXISTS public."order_operations"
(
	id SERIAL PRIMARY KEY,
	value VARCHAR(255)
) INHERITS(public."base");

COMMENT ON TABLE public."order_operations" IS 'operations that can be performed on order';

CREATE TABLE public."order_audit"
(
    id SERIAL PRIMARY KEY,
    order_operation_id INT,
    new_record_id INT,
    old_record_id INT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	changed_by INT REFERENCES public."users"(id),
	FOREIGN KEY (order_operation_id) REFERENCES public."order_operations"(id) ON DELETE CASCADE,
	FOREIGN KEY (old_record_id) REFERENCES public."orders"(id) ON DELETE CASCADE,
	FOREIGN KEY (new_record_id) REFERENCES public."orders"(id) ON DELETE CASCADE
);

COMMENT ON TABLE public."order_audit" IS 'audit trail for order table';


CREATE TABLE IF NOT EXISTS public."o_items"
(
	order_id INT,
	product_id INT,
	distributor_id INT,
	FOREIGN KEY (order_id) REFERENCES public."orders"(id) ON DELETE CASCADE,
	FOREIGN KEY (product_id) REFERENCES public."products"(id) ON DELETE CASCADE,
	FOREIGN KEY (distributor_id) REFERENCES public."distributors"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."o_items" IS 'stores order items';

CREATE TABLE IF NOT EXISTS public."i_statuses"
(
	id SERIAL PRIMARY KEY, 
	value VARCHAR(255),
	description TEXT
) INHERITS(public."base");

COMMENT ON TABLE public."i_statuses" IS 'stores invoice statuses';

CREATE TABLE IF NOT EXISTS public."payment_methods"
(
    id SERIAL PRIMARY KEY, 
	value VARCHAR(255),
	description TEXT
) INHERITS(public."base");

COMMENT ON TABLE public."payment_methods" IS 'stores payment options';

CREATE TABLE IF NOT EXISTS public."invoices"
(
    id SERIAL PRIMARY KEY,
	external_id VARCHAR(255),
	order_id INT,
	status_id INT,
	payment_method_id INT,
	FOREIGN KEY (order_id) REFERENCES public."orders"(id) ON DELETE CASCADE,
	FOREIGN KEY (status_id) REFERENCES public."i_statuses"(id) ON DELETE CASCADE,
	FOREIGN KEY (payment_method_id) REFERENCES public."payment_methods"(id) ON DELETE CASCADE
) INHERITS(public."base");

COMMENT ON TABLE public."invoices" IS 'stores invoices';

CREATE TABLE IF NOT EXISTS public."invoice_operations"
(
	id SERIAL PRIMARY KEY,
	value VARCHAR(255)
) INHERITS(public."base");

COMMENT ON TABLE public."invoice_operations" IS 'operations that can be performed on invoice';

CREATE TABLE public."invoice_audit"
(
    id SERIAL PRIMARY KEY,
    invoice_operation_id INT,
    new_record_id INT,
    old_record_id INT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	changed_by INT REFERENCES public."users"(id),
	FOREIGN KEY (invoice_operation_id) REFERENCES public."invoice_operations"(id) ON DELETE CASCADE,
	FOREIGN KEY (old_record_id) REFERENCES public."invoices"(id) ON DELETE CASCADE,
	FOREIGN KEY (new_record_id) REFERENCES public."invoices"(id) ON DELETE CASCADE
);

COMMENT ON TABLE public."invoice_audit" IS 'audit trail for invoice table';
