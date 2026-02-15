# Credentials Reference

All default credentials used by the B2B / Efoyeta Store stack. Override where possible via `.env` (see below); do not commit real secrets to the repo.

---

## Postgres (database)

Used by: `db` container, `migrations`, `seed`, `backend`.

| Variable | Default | Override |
|----------|---------|----------|
| User | `postgres` | `POSTGRES_USER` in `.env` |
| Password | `postgres` | `POSTGRES_PASSWORD` in `.env` |
| Database | `b2b` | `POSTGRES_DB` in `.env` |
| Port | `5432` | (fixed in compose) |

To override: copy `.env.example` to `.env` and set the variables. Compose and the backend read them at runtime.

---

## PgAdmin

Used for the DB UI (e.g. http://localhost:5051/ or http://your-host/pgadmin/).

| Field | Value |
|-------|--------|
| Email | `admin@example.com` |
| Password | `admin` |

---

## App Basic Auth (API)

The backend uses HTTP Basic Auth for protected API routes. These come from config and/or the database.

### Bootstrap user (from config)

Created from `config/config.docker.compose.yaml.tpl` (or equivalent). Used for initial bootstrap.

| Username | Password |
|----------|----------|
| `superadmin` | `changeme` |

### Seeded users (after `db/seed`)

Created by `db/seed/001_seed_users.sql` and `db/seed/008_seed_credentials.sql`. Use **email** (or username) and password to call the API with Basic Auth.

| Email (login) | Password | Role |
|---------------|----------|------|
| `admin-admin@admin.local` | `changeme` | Officer / Super Administrator |
| `supplier1@b2b.local` | `password123` | Supplier |
| `supplier2@b2b.local` | `password123` | Supplier |
| `customer1@b2b.local` | `password123` | Customer |
| `customer2@b2b.local` | `password123` | Customer |

---

## Overriding credentials

- **Postgres:** Set `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` in `.env` in the project root. Used by compose and the backend (via config template).
- **PgAdmin:** Change `PGADMIN_DEFAULT_EMAIL` and `PGADMIN_DEFAULT_PASSWORD` in `docker-compose.yml` (or via env if you add support). If you already logged in once, you may need to remove the `pgadmin_data` volume and restart to apply a new email/password.
- **App users:** Change bootstrap user in config; change seeded users by editing `db/seed/001_seed_users.sql` and `db/seed/008_seed_credentials.sql` (use bcrypt for passwords) and re-running seed, or manage via the API after first login.

---

## Security notes

- `.env` is gitignored; do not commit real passwords.
- Defaults here are for local/dev only. Use strong, unique credentials in production.
- Change default Postgres and PgAdmin passwords before exposing the server.
