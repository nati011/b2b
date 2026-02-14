package middleware

// Action names for HTTP method mappings
// These are generic actions derived from HTTP methods
const (
	ActionView   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// Resource and path segment constants used in authorization middleware
// Note: These match constants defined in user module but are defined here
// to avoid import cycles. If user module constants change, these must be updated.
const (
	ResourceUsers         = "users"
	PathSegmentRoles      = "roles"
	PathSegmentActivate   = "activate"
	PathSegmentDeactivate = "deactivate"
	PathSegmentSuspend    = "suspend"
	ActionAssignRoles     = "assign_roles"
)
