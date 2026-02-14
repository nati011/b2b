package middleware

import (
	"context"
	"net/http"
	"strings"

	resourceDomain "marketplace/internal/infra/authz/resource/domain"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/logger"
)

const (
	// UUIDLength is the standard UUID string length (with hyphens)
	UUIDLength = 36
	// UUIDHyphenCount is the number of hyphens in a standard UUID
	UUIDHyphenCount = 4
	// MaxSegmentLengthForID is the maximum length for a segment to be considered an ID
	MaxSegmentLengthForID = 50
	// MaxQualifierLength is the maximum length for a segment to be considered a qualifier
	MaxQualifierLength = 20
	// MinResourceNameLength is the minimum length for a resource name
	MinResourceNameLength = 3
)

// ResourceCatalog defines the interface for validating resources and actions.
type ResourceCatalog interface {
	ResolveAction(ctx context.Context, code, action string) (*resourceDomain.Resource, *resourceDomain.Action, error)
}

// PermissionChecker defines the interface for checking user permissions.
type PermissionChecker interface {
	HasPermission(ctx context.Context, userID, resource, action string) (bool, error)
}

// AuthorizationMiddleware handles authorization by checking if the authenticated user
// has the required permission for the requested resource and action.
type AuthorizationMiddleware struct {
	resourceCatalog   ResourceCatalog
	permissionChecker PermissionChecker
	publicRoutes      []string
	allowedRoutes     []string // Routes accessible to all authenticated users (bypass permission checks)
}

// NewAuthorizationMiddleware creates a new authorization middleware.
func NewAuthorizationMiddleware(resourceCatalog ResourceCatalog, permissionChecker PermissionChecker, publicRoutes []string, allowedAuthenticatedRoutes []string) *AuthorizationMiddleware {
	if publicRoutes == nil {
		publicRoutes = []string{}
	}
	if allowedAuthenticatedRoutes == nil {
		allowedAuthenticatedRoutes = []string{}
	}
	return &AuthorizationMiddleware{
		resourceCatalog:   resourceCatalog,
		permissionChecker: permissionChecker,
		publicRoutes:      publicRoutes,
		allowedRoutes:     allowedAuthenticatedRoutes,
	}
}

// Handle processes the request to authorize the user.
func (m *AuthorizationMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if route is public
		if m.isPublicRoute(r.Method, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get user from context (set by auth middleware)
		user := httputil.UserFromContext(r.Context())
		if user == nil {
			// No user in context - this means authentication middleware didn't set one
			// For now, we allow unauthenticated requests (auth middleware handles auth)
			// In strict mode, we could return 401 here
			next.ServeHTTP(w, r)
			return
		}

		// Check if route is in allowed list (accessible to all authenticated users)
		if m.isAllowedRoute(r.Method, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Derive required permission from route
		resource, action, err := derivePermission(r.Method, r.URL.Path)
		if err != nil {
			// Could not derive permission - allow request (might be a route we don't protect)
			next.ServeHTTP(w, r)
			return
		}

		// Validate that resource and action exist in catalog
		ctx := r.Context()
		// Log what we're checking for debugging (using Info level so it shows up)
		if resource != "" && action != "" {
			logger.InfoContext(ctx, "Checking authorization", "resource", resource, "action", action, "path", r.URL.Path, "method", r.Method)
		}
		_, _, err = m.resourceCatalog.ResolveAction(ctx, resource, action)
		if err != nil {
			// Log the error for debugging
			logger.ErrorContext(ctx, "Resource/action resolution failed", "resource", resource, "action", action, "error", err, "path", r.URL.Path)
			// Resource/action doesn't exist in catalog - return 404
			httputil.JSON(w, http.StatusNotFound, map[string]string{
				"error": "resource or action not found",
			})
			return
		}

		// Check if user has the required permission
		hasPerm, err := m.permissionChecker.HasPermission(ctx, user.ID, resource, action)
		if err != nil || !hasPerm {
			httputil.JSON(w, http.StatusForbidden, map[string]string{
				"error": "insufficient permissions",
			})
			return
		}

		// User has permission - continue to handler
		next.ServeHTTP(w, r)
	})
}

// isPublicRoute checks if the given method and path match any public route.
// Public routes can be specified as:
// - Exact path: "/users" (matches any method on /users)
// - Method-specific: "GET /users/me" (matches only GET /users/me)
// - Path prefix: "/public" (matches /public and /public/* for any method)
func (m *AuthorizationMiddleware) isPublicRoute(method, path string) bool {
	for _, publicRoute := range m.publicRoutes {
		// Check for method-specific routes (e.g., "GET /users/me")
		if strings.HasPrefix(publicRoute, method+" ") {
			routePath := strings.TrimPrefix(publicRoute, method+" ")
			if path == routePath || strings.HasPrefix(path, routePath+"/") {
				return true
			}
		}
		// Check for path-only routes (matches any method)
		if !strings.Contains(publicRoute, " ") {
			if path == publicRoute || strings.HasPrefix(path, publicRoute+"/") {
				return true
			}
		}
	}
	return false
}

// isAllowedRoute checks if the given method and path match any route in the allowed list.
// Allowed routes require authentication but bypass permission checks.
func (m *AuthorizationMiddleware) isAllowedRoute(method, path string) bool {
	routeKey := method + " " + path
	for _, allowedRoute := range m.allowedRoutes {
		if routeKey == allowedRoute {
			return true
		}
		// Support method-specific routes (e.g., "GET /users/me")
		if strings.HasPrefix(allowedRoute, method+" ") {
			routePath := strings.TrimPrefix(allowedRoute, method+" ")
			if path == routePath || strings.HasPrefix(path, routePath+"/") {
				return true
			}
		}
	}
	return false
}

// derivePermission extracts the resource and action from the HTTP method and path.
// Returns resource, action, and error.
func derivePermission(method, path string) (resource, action string, err error) {
	segments := normalizePath(path)
	if len(segments) == 0 {
		return "", "", nil
	}

	// Try multi-segment resource first
	if resource, action, ok := handleMultiSegmentResource(method, segments); ok {
		return resource, action, nil
	}

	// Handle single-segment resources
	return handleSingleSegmentResource(method, segments)
}

// normalizePath normalizes the path and returns segments
func normalizePath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

// isLikelyQualifier checks if a segment is likely a qualifier (not an ID or action).
// Qualifiers are typically short, lowercase words that modify the resource type.
func isLikelyQualifier(segment string) bool {
	// Known actions are not qualifiers
	if isKnownAction(segment) {
		return false
	}

	// UUIDs are not qualifiers (basic check: contains hyphens and is long enough)
	// UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if len(segment) == UUIDLength && strings.Count(segment, "-") == UUIDHyphenCount {
		return false
	}

	// Very long segments are likely IDs, not qualifiers
	if len(segment) > MaxSegmentLengthForID {
		return false
	}

	// Qualifiers are typically short, lowercase alphanumeric words (may contain hyphens)
	// They should not start with numbers (IDs often do)
	if len(segment) > 0 && segment[0] >= '0' && segment[0] <= '9' {
		return false
	}

	// Simple heuristic: if it's a short word with only letters/hyphens, it might be a qualifier
	// This is a conservative check - we'll only treat it as a qualifier if it's clearly not an ID
	return len(segment) <= MaxQualifierLength && isAlphanumericOrHyphen(segment)
}

// isAlphanumericOrHyphen checks if a string contains only alphanumeric characters and hyphens
func isAlphanumericOrHyphen(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return true
}

// isLikelyID checks if a segment is likely an ID (UUID or numeric)
func isLikelyID(segment string) bool {
	// UUID format check
	if len(segment) == UUIDLength && strings.Count(segment, "-") == UUIDHyphenCount {
		return true
	}
	// Very long segments are likely IDs
	if len(segment) > MaxSegmentLengthForID {
		return true
	}
	// Segments starting with numbers are likely IDs
	if len(segment) > 0 && segment[0] >= '0' && segment[0] <= '9' {
		return true
	}
	return false
}

// isKnownAction checks if a segment is a known action verb.
// This is shared between qualifier detection and action extraction.
func isKnownAction(segment string) bool {
	knownActions := map[string]bool{
		"approve":      true,
		"reject":       true,
		"withdraw":     true,
		"disburse":     true,
		"repayments":   true,
		"repayment":    true,
		"close":        true,
		"write-off":    true,
		"activate":     true,
		"deactivate":   true,
		"schedule":     true,
		"transactions": true,
		"charges":      true,
		"summary":      true,
		"template":     true,
		"versions":     true,
		"collaterals":  true,
		"assign":       true,
		"unassign":     true,
		"terminate":    true,
		"maintenance":  true,
		"dispose":      true,
	}
	return knownActions[segment]
}

// handleMultiSegmentResource handles multi-segment resource paths generically.
// It converts path segments like ["loan", "accounts"] to resource code "loan-accounts".
// For paths like ["loan", "accounts", "{id}", "approve"], it extracts the action from the last segment.
// It automatically detects qualifiers in the third position and rearranges them.
// It also handles namespace prefixes where the first segment is a domain/organizational prefix.
func handleMultiSegmentResource(method string, segments []string) (resource, action string, ok bool) {
	if len(segments) < 2 {
		return "", "", false
	}

	// If the second segment is an ID (UUID), this is not a multi-segment resource
	// It's a single-segment resource with an ID parameter (e.g., /clients/{id})
	// Let handleSingleSegmentResource handle it instead
	if isLikelyID(segments[1]) {
		return "", "", false
	}

	// Check if first segment is likely a namespace prefix (domain/organizational word)
	// and second segment is likely a resource name (plural noun)
	// If so, use just the second segment as the resource code
	if isLikelyNamespacePrefix(segments[0]) && isLikelyResourceName(segments[1]) {
		resource = segments[1]
		// Extract action from remaining segments or HTTP method
		action = extractActionFromSegments(method, segments[2:])
		if action == "" {
			return "", "", false
		}
		return resource, action, true
	}

	// Check if third segment is a qualifier and handle it
	if resource, action, ok := handleQualifierResource(method, segments); ok {
		return resource, action, true
	}

	// Default pattern: Convert first two segments to hyphenated resource code
	// e.g., ["loan", "accounts"] -> "loan-accounts"
	resource = segments[0] + "-" + segments[1]

	// Extract action from remaining segments or HTTP method
	action = extractActionFromSegments(method, segments[2:])
	if action == "" {
		return "", "", false
	}

	return resource, action, true
}

// isLikelyNamespacePrefix checks if a segment is likely a namespace/domain prefix
// rather than part of a resource name. These are typically organizational or domain-level words.
func isLikelyNamespacePrefix(segment string) bool {
	// Common namespace prefixes that indicate organizational/domain grouping
	namespacePrefixes := map[string]bool{
		"organization": true,
		"portfolio":    true,
		"core":         true,
		"domain":       true,
		"system":       true,
		"accounting":   true,
	}
	return namespacePrefixes[segment]
}

// isLikelyResourceName checks if a segment is likely a resource name (plural noun)
// rather than a namespace prefix. Resource names are typically plural nouns.
func isLikelyResourceName(segment string) bool {
	// Resource names are typically plural (end with 's', 'es', 'ies', etc.)
	// and are common resource types. We use a simple heuristic: if it ends with
	// a plural suffix and is not too short, it's likely a resource name.
	if len(segment) < MinResourceNameLength {
		return false
	}

	// Special cases for uncountable nouns that can be used as plurals
	uncountableNouns := map[string]bool{
		"stuff":       true,
		"staff":       true,
		"equipment":   true,
		"furniture":   true,
		"information": true,
	}
	if uncountableNouns[segment] {
		return true
	}

	// Check if it ends with common plural suffixes
	pluralSuffixes := []string{"s", "es", "ies"}
	for _, suffix := range pluralSuffixes {
		if strings.HasSuffix(segment, suffix) {
			// Additional check: ensure it's not just the suffix (e.g., "s" alone)
			// and has a reasonable length for a resource name
			return len(segment) > len(suffix)
		}
	}

	return false
}

// handleQualifierResource handles resource paths with qualifiers in the third position.
// Returns resource, action, and ok flag if a qualifier pattern was matched.
func handleQualifierResource(method string, segments []string) (resource, action string, ok bool) {
	if len(segments) < 3 || !isLikelyQualifier(segments[2]) {
		return "", "", false
	}

	// If there's a next segment and it's a known action, this segment is definitely a qualifier
	// e.g., /deposit/accounts/fixed/activate -> "fixed" is qualifier, "activate" is action
	if len(segments) >= 4 && isKnownAction(segments[3]) {
		resource = segments[2] + "-" + segments[0] + "-" + segments[1]
		action = extractActionFromSegments(method, segments[3:])
		if action != "" {
			return resource, action, true
		}
	}

	// If no next segment or next segment looks like an ID, this might be a qualifier
	// But be conservative: only treat as qualifier if next segment is clearly an ID
	if len(segments) == 3 || (len(segments) >= 4 && isLikelyID(segments[3])) {
		resource = segments[2] + "-" + segments[0] + "-" + segments[1]
		action = extractActionFromSegments(method, segments[3:])
		if action != "" {
			return resource, action, true
		}
	}

	return "", "", false
}

// extractActionFromSegments extracts the action from path segments.
// It handles patterns like:
// - [] (no segments) -> use HTTP method mapping
// - ["{id}"] -> use HTTP method mapping (GET -> view, POST -> create, etc.)
// - ["{id}", "action"] -> use the action verb (approve, reject, etc.)
// - ["action"] -> use the action verb if it's a known action, otherwise treat as ID
func extractActionFromSegments(method string, remainingSegments []string) string {
	if len(remainingSegments) == 0 {
		// No additional segments, use HTTP method mapping
		// GET /loan/accounts -> view (list)
		// POST /loan/accounts -> create
		return mapHTTPMethodToAction(method)
	}

	// Check if the last segment is a known action verb
	lastSegment := remainingSegments[len(remainingSegments)-1]
	if isKnownAction(lastSegment) {
		// Map action names to their canonical action codes
		actionMap := map[string]string{
			"repayments":   "repayment", // Normalize plural to singular
			"schedule":     ActionView,  // GET endpoint returning view data
			"transactions": ActionView,  // GET endpoint returning view data
			"charges":      ActionView,  // GET endpoint returning view data
			"summary":      ActionView,  // GET endpoint returning view data
		}
		if mappedAction, ok := actionMap[lastSegment]; ok {
			return mappedAction
		}
		// Most actions map to themselves
		return lastSegment
	}

	// If we have segments but the last one isn't a known action,
	// it's likely an ID, so use HTTP method mapping
	// GET /loan/accounts/{id} -> view (get one)
	// PUT /loan/accounts/{id} -> update
	// DELETE /loan/accounts/{id} -> delete
	return mapHTTPMethodToAction(method)
}

// handleSingleSegmentResource handles single-segment resource paths
func handleSingleSegmentResource(method string, segments []string) (resource, action string, err error) {
	resource = segments[0]

	// Map singular resource names to their plural forms (e.g., product -> products)
	resource = mapResourceToPlural(resource)

	// Handle special cases first
	if len(segments) >= 2 {
		if action := handleSpecialCases(method, segments); action != "" {
			return resource, action, nil
		}
	}

	// Map HTTP method to action for standard CRUD operations
	action = mapHTTPMethodToAction(method)
	if action == "" {
		return "", "", nil
	}

	return resource, action, nil
}

// handleSpecialCases processes special path patterns and returns the action if matched.
func handleSpecialCases(method string, segments []string) string {
	// Check for nested resources with special actions
	if segments[0] == ResourceUsers && segments[1] == PathSegmentRoles {
		return handleUserRolesAction(method)
	}

	// Check for deprecate action
	if len(segments) >= 2 && segments[len(segments)-1] == "deprecate" {
		if segments[0] == "resources" {
			return "deprecate"
		}
	}

	// Check for custom action verbs (activate, deactivate, suspend)
	if len(segments) >= 2 {
		customAction := segments[1]
		if customAction == PathSegmentActivate || customAction == PathSegmentDeactivate || customAction == PathSegmentSuspend {
			// These are user management actions - map to update
			return ActionUpdate
		}
	}

	return ""
}

// handleUserRolesAction maps HTTP methods to actions for user roles endpoints.
func handleUserRolesAction(method string) string {
	switch method {
	case http.MethodPost, http.MethodDelete:
		return ActionAssignRoles
	case http.MethodGet:
		return ActionView
	default:
		return ""
	}
}

// mapResourceToPlural maps singular resource names to their plural forms.
// This handles cases where routes use singular forms (e.g., /product) but
// resource codes use plural forms (e.g., products).
func mapResourceToPlural(resource string) string {
	singularToPlural := map[string]string{
		"product":  "products",
		"customer": "customers", // In case it's used
		"supplier": "suppliers", // In case it's used
		"order":    "orders",    // In case it's used
	}
	if plural, ok := singularToPlural[resource]; ok {
		return plural
	}
	return resource
}

// mapHTTPMethodToAction maps standard HTTP methods to actions.
func mapHTTPMethodToAction(method string) string {
	switch method {
	case http.MethodGet:
		return ActionView
	case http.MethodPost:
		return ActionCreate
	case http.MethodPut, http.MethodPatch:
		return ActionUpdate
	case http.MethodDelete:
		return ActionDelete
	default:
		return ""
	}
}
