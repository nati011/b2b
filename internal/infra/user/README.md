# User Module

User entity and lifecycle management.

## Overview

The User module manages system users and their profiles. It handles user lifecycle, role assignments, and provides the foundation for access control across the system. **This module is authentication scheme agnostic** - it works with any authentication mechanism (Basic Auth, OAuth2, etc.) through the middleware layer.

## Domain Boundaries

**Owns**:
- User accounts and profiles
- User status and lifecycle (active, inactive, suspended)
- User types (self-service, officer)
- User-role assignments

**References** (via IDs, not domain objects):
- Role IDs (for role assignments, roles owned by authz/role)
- Permission IDs (for permission checks, permissions owned by authz/permission)

**Does NOT Own**:
- Authentication credentials (handled by auth middleware layer)
- Authentication mechanisms (Basic Auth, OAuth2, etc. - handled by middleware)
- Clients (owned by Client domain - users may be linked to clients)
- Accounts (owned by Portfolio)
- Transactions (owned by Settlement)
- Roles and permissions (owned by authz modules)

**Interaction Patterns**:
- Publishes events: `UserCreated`, `UserUpdated`, `UserStatusChanged`, `RoleAssigned`
- Consumes events: `RoleCreated` (to validate role assignments)
- Used by all domains for access control
- Receives authenticated user context from middleware layer (via `httputil.UserFromContext()`)

## Features

- User account management
- User status management (active, inactive, suspended)
- User types (self-service, officer)
- Role assignment and management
- User profile management

## Authentication Scheme Agnostic

The user module is designed to work with any authentication scheme:
- **Basic Auth**: Credentials configured in config.yaml, users created automatically during bootstrap
- **OAuth2**: Tokens validated by middleware, user loaded from database
- **Other schemes**: Any authentication mechanism that places a `*domain.User` in the request context via `httputil.WithUser()`

The module does not implement authentication logic - it only manages user entities and their lifecycle. Authentication is handled by the middleware layer (`pkg/http/middleware/auth.go`), which validates credentials/tokens and loads users into the request context.

