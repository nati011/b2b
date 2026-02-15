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

The `docker-compose.yml` stack runs **migrations** then **seed** after the database is healthy. So on `docker compose up -d`, the DB is migrated and seeded automatically.

### Manual clear and seed (scripts SQL)

To clear and reseed using the legacy SQL files in `scripts/` (e.g. `clear_database.sql`, `seed_dummy_data.sql`):

```bash
./scripts/clear_and_seed.sh
# or seed only (no clear):
./scripts/seed_database.sh
```

These require the **db** service to be running and use `docker-compose.yml` (or `docker compose`).

## Usage Examples

### Start the full stack
```bash
docker compose up -d
# or
docker compose -f docker-compose.yml up -d
```

### Manually clear and seed (scripts)
```bash
./scripts/clear_and_seed.sh
```

### Stop app (keep db and pgadmin running)
```bash
./scripts/stop_containers.sh
```

### Stop everything
```bash
docker compose down
```
