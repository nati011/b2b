# Database Seeding Scripts

This directory contains scripts for seeding the database with dummy data for testing purposes.

## Files

- `seed_dummy_data.sql` - SQL script containing all the seed data
- `seed_database.sh` - Bash script to easily run the seed SQL file

## Usage

### Option 1: Using the Bash Script (Recommended)

```bash
./scripts/seed_database.sh
```

This script will:
- Check if Docker is running
- Check if the database container is running (and start it if needed)
- Run the SQL seed file
- Display verification statistics

### Option 2: Using Docker Compose Directly

```bash
docker compose exec db psql -U postgres -d b2b < scripts/seed_dummy_data.sql
```

### Option 3: Using psql directly (if you have direct database access)

```bash
psql -U postgres -d b2b -f scripts/seed_dummy_data.sql
```

## What Gets Created

The seed script creates:

1. **Categories** (8 categories):
   - Electronics
   - Clothing
   - Home & Kitchen
   - Sports & Outdoors
   - Beauty & Personal Care
   - Books
   - Toys & Games
   - Automotive

2. **Test User and Distributor**:
   - Test distributor user account
   - Active distributor with business info
   - Active subscription (Year Plan)

3. **Simple Products** (6 products):
   - Wireless Bluetooth Headphones (Electronics)
   - Premium Cotton T-Shirt (Clothing)
   - Automatic Coffee Maker (Home & Kitchen)
   - Professional Running Shoes (Sports & Outdoors)
   - Hydrating Face Moisturizer (Beauty & Personal Care)
   - Complete Guide to Web Development (Books)

4. **Configurable Products** (2 configurable products):
   - **Smartphone Pro** with 3 variants (64GB, 128GB, 256GB)
   - **Gaming Laptop** with 2 variants (8GB RAM, 16GB RAM)

All products include:
- Images (using Unsplash placeholders)
- Categories
- Stock quantities
- Active status

## Verification

After running the seed script, you can verify the data with:

```sql
-- Count categories
SELECT COUNT(*) as total_categories FROM public.category WHERE is_deleted = false;

-- Count active products
SELECT COUNT(*) as total_products FROM public.products WHERE is_deleted = false AND is_active = true;

-- Count configurable products
SELECT COUNT(*) as total_configurable_products FROM public.configurable_products WHERE is_deleted = false AND is_available = true;

-- Check distributor subscription
SELECT d.id, d.is_active, ds.status 
FROM public.distributors d 
JOIN public.distributor_subscriptions ds ON d.id = ds.distributor_id 
WHERE d.is_deleted = false AND ds.status = 'active';
```

## Notes

- The script uses `BEGIN` and `COMMIT` transactions, so all data is created atomically
- If categories or subscription plans already exist, they won't be duplicated
- The script is idempotent - you can run it multiple times safely
- All products are created with active status and stock quantities
- Images use Unsplash placeholder URLs

## Troubleshooting

If you encounter errors:

1. **Docker not running**: Make sure Docker is running
2. **Database container not running**: The script will try to start it automatically
3. **Permission denied**: Make sure the script is executable: `chmod +x scripts/seed_database.sh`
4. **SQL errors**: Check the error messages for specific issues


