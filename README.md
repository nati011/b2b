## Efoyeta Store
A complete Marketplace solution for Ethiopia.

## Setup and run

### Prerequisites
- Docker and Docker Compose
- Git (to clone the repo)

### 1. Clone and enter the project
```bash
git clone <repo-url> b2b && cd b2b
```

### 2. (Optional) Set Postgres credentials
Defaults are `postgres` / `postgres` / `b2b`. To override, create a `.env` in the project root:
```bash
cp .env.example .env
# Edit .env and set POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB if needed
```

### 3. Start everything
From the project root:
```bash
docker compose up -d --build
```
This will:
- Start Postgres and wait until it is healthy
- Run migrations, then seed data
- Start backend, then frontend, then nginx

First run can take a few minutes (builds + migrations + seed). Check status:
```bash
docker compose ps -a
```

### 4. URLs
| Service   | URL                    |
|----------|-------------------------|
| App      | http://localhost/       |
| Backend  | http://localhost:8090/  |
| PgAdmin  | http://localhost:5051/  (admin@local.local / admin) |

### 5. Useful commands
```bash
# View logs (all or a service)
docker compose logs -f
docker compose logs -f backend

# Stop everything
docker compose down

# Stop and remove volumes (resets DB)
docker compose down -v

# Rebuild and start after code/config changes
docker compose up -d --build
```

### 6. Deploy on a server
1. Copy the repo to the server (e.g. `/root/b2b`).
2. Create `.env` with your `POSTGRES_PASSWORD` (and optionally user/db).
3. Run from the project directory:
   ```bash
   docker compose up -d --build
   ```
4. Ensure firewall allows ports 80 (nginx), and optionally 8090, 5432, 5051 if you need direct access.
5. For production, put nginx or a reverse proxy in front and use TLS; the app is served on port 80 by the compose nginx service.
