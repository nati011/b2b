# Configuration Service (Hot-Reloadable)

Single configuration management service with real-time updates, no server restarts.

## Goals
- Single source of truth for runtime configuration
- Merge multiple sources (file, env, flags, remote)
- Validate against schema before apply
- Hot reload with watchers and change notifications
- Backward-compatible schema evolution

## Directory Structure
```
internal/config/
├── adapters/           # Integration adapters for consumers (HTTP, gRPC, app)
├── registry/           # Central registry of config sections & defaults
├── schema/             # JSON/YAML schemas (validation)
├── sources/
│   ├── env/            # Environment variables loader
│   ├── file/           # YAML/JSON/TOML file loader (configs/)
│   ├── flags/          # Command-line flags loader
│   └── remote/         # Remote KV/Secrets (Consul/Etcd/Vault) loader
└── watcher/            # File/remote watchers, debounce, events
```

## Concepts
- Providers: load config fragments from sources
- Merger: deterministic merge (priority: flags > env > remote > file > defaults)
- Validator: validate merged config against schemas
- Snapshot: immutable view of current config
- Watcher: emits change events; consumers subscribe

## Update Flow
1) Load defaults → merge file/env/flags/remote → validate
2) On change event → reload changed sections only → validate → publish
3) Consumers receive typed section updates via adapters

## Consumers
- Application services (DB pool, HTTP server, posting engine)
- Middleware (CORS, rate limits)
- Feature toggles

## Non-Goals
- Persisting secrets (use remote source)
- Business logic in this layer