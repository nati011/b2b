# Docker Compose config: backend connects to service name "db"
# Password is injected at runtime from POSTGRES_* env vars
app:
  name: "marketplace"
  version: "0.5.0-Snapshot"
  env: "development"

server:
  host: "0.0.0.0"
  port: 8090
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 60s

logging:
  level: "info"
  format: "json"
  source_location: false
  sampling_rate: 1.0
  async: false
  sanitize: false

db:
  driver: postgres
  host: db
  port: 5432
  user: ${POSTGRES_USER}
  password: ${POSTGRES_PASSWORD}
  database: ${POSTGRES_DB}
  sslmode: disable

auth:
  mode: basic
  basic:
    realm: "B2B"
    email_domain: "@admin.local"
    user_prefix: "admin-"
    user_name: "Super Administrator"
    users:
      - username: "superadmin"
        password: "changeme"

resources:
  path: "config/resources.yaml"

roles:
  path: "config/roles.yaml"
