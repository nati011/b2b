# Database Seed Data

This directory contains SQL seed files for populating the database with development/test data.

## Files

The seed files should be executed in the following order:

1. `001_seed_users.sql` - Creates test users (admin, suppliers, customers)
2. `008_seed_credentials.sql` - Creates authentication credentials for seeded users (password: password123)
3. `009_seed_user_roles.sql` - Assigns roles to users based on their user_type (must run after roles are bootstrapped by the application)
4. `002_seed_suppliers.sql` - Creates supplier records (8 suppliers)
5. `003_seed_customers.sql` - Creates customer records (10 customers)
6. `004_seed_categories.sql` - Creates product categories (10 categories)
7. `005_seed_products.sql` - Creates product records (24 products across different categories)
8. `006_seed_product_categories.sql` - Maps products to categories
9. `007_seed_orders.sql` - Creates sample orders and order items (5 orders)

**Note:** Roles are configured in `config/roles.yaml` and automatically bootstrapped by the application during initialization (see `internal/app/bootstrap.go`). The `009_seed_user_roles.sql` file assigns these roles to seeded users and must be run after the application has bootstrapped the roles.

## Usage

### Manual Execution

Execute the seed files in order using psql:

```bash
psql -h localhost -U postgres -d b2b -f db/seed/001_seed_users.sql
psql -h localhost -U postgres -d b2b -f db/seed/008_seed_credentials.sql
psql -h localhost -U postgres -d b2b -f db/seed/002_seed_suppliers.sql
psql -h localhost -U postgres -d b2b -f db/seed/003_seed_customers.sql
psql -h localhost -U postgres -d b2b -f db/seed/004_seed_categories.sql
psql -h localhost -U postgres -d b2b -f db/seed/005_seed_products.sql
psql -h localhost -U postgres -d b2b -f db/seed/006_seed_product_categories.sql
psql -h localhost -U postgres -d b2b -f db/seed/007_seed_orders.sql
```

### Using a Script

You can create a script to run all seeds:

```bash
#!/bin/bash
for file in db/seed/*.sql; do
  echo "Seeding $file..."
  psql -h localhost -U postgres -d b2b -f "$file"
done
```

### Docker Execution

If running in Docker:

```bash
docker exec -i b2b-db psql -U postgres -d b2b < db/seed/001_seed_users.sql
# ... repeat for other files
```

## Data Overview

### Users
- 1 admin user
- 2 supplier users
- 2 customer users

### Credentials
- Authentication credentials for all seeded users
- Default password for all users: `password123`
- Credentials are created for:
  - `customer1@b2b.local` (Customer One)
  - `customer2@b2b.local` (Customer Two)
  - `supplier1@b2b.local` (Supplier One)
  - `supplier2@b2b.local` (Supplier Two)

### Suppliers
- 8 active suppliers with business names and contact information

### Customers
- 10 active customers with addresses in Addis Ababa

### Categories
- 10 product categories (Electronics, Office Supplies, Tools, Furniture, etc.)

### Products
- 24 products across 8 categories
- Products include electronics, office supplies, tools, furniture, safety equipment, industrial items, construction materials, and cleaning supplies
- Each product has realistic attributes, pricing, and inventory levels

### Orders
- 5 sample orders with various statuses (pending, confirmed, shipped, delivered)
- Orders include order items linking to products

## Notes

- **Roles**: Roles are configured in `config/roles.yaml` and automatically bootstrapped by the application during startup. Do not seed roles via SQL.
- **Credentials**: The `008_seed_credentials.sql` file must be executed after `001_seed_users.sql` since credentials reference user IDs. All passwords are hashed using bcrypt.
- All seed files use `ON CONFLICT DO NOTHING` to allow safe re-execution
- Timestamps are set to realistic values (some orders are from days ago)
- Prices are in ETB (Ethiopian Birr)
- All data is suitable for development/testing purposes only

## Default Login Credentials

After seeding, you can log in with any of these accounts:

- **Email:** `customer1@b2b.local` | **Password:** `password123`
- **Email:** `customer2@b2b.local` | **Password:** `password123`
- **Email:** `supplier1@b2b.local` | **Password:** `password123`
- **Email:** `supplier2@b2b.local` | **Password:** `password123`

