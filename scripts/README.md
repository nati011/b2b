# Database Scripts

This directory contains scripts for managing the database.

## Scripts

### `clear_database.sql`
SQL script that truncates all tables in the database, removing all data while preserving the schema structure. It also resets all sequences to start from 1.

### `seed_dummy_data.sql`
SQL script that seeds the database with dummy data for testing purposes, including:
- Categories
- Test user and distributor
- Products (simple and configurable)
- Subscriptions

### `clear_and_seed.sh`
Bash script that clears the database and then seeds it with fresh data. Can be run manually:
```bash
./scripts/clear_and_seed.sh
```

### `seed_database.sh`
Bash script that seeds the database with dummy data (without clearing first).

## Automatic Database Initialization

The `docker-compose.yaml` includes a `db-init` service that automatically clears and seeds the database when you start the containers.

### Enable/Disable Auto-Seeding

To control automatic seeding, set the `AUTO_SEED` environment variable:

**Enable auto-seeding (default):**
```bash
docker-compose up
# or explicitly:
AUTO_SEED=true docker-compose up
```

**Disable auto-seeding:**
```bash
AUTO_SEED=false docker-compose up
```

### Manual Database Management

If you want to manage the database manually without the auto-init service:

1. **Remove the db-init service** from `docker-compose.yaml`, or
2. **Set AUTO_SEED=false** when starting containers

Then use the scripts manually:
```bash
# Clear and seed
./scripts/clear_and_seed.sh

# Or just seed (without clearing)
./scripts/seed_database.sh
```

## Usage Examples

### Start containers with auto-seeding (default)
```bash
docker-compose up -d
```

### Start containers without auto-seeding
```bash
AUTO_SEED=false docker-compose up -d
```

### Manually clear and seed after containers are running
```bash
./scripts/clear_and_seed.sh
```

### Stop containers
```bash
./scripts/stop_containers.sh
# or
docker-compose stop client admin
```
