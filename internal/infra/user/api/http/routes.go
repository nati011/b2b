package user

import (
	"net/http"
	"strings"

	httputil "marketplace/pkg/http"
	"marketplace/pkg/http/middleware"
)

// HTTP route path
const RouteUsers = "/users"

// Resource code
const ResourceUsers = "users"

// Path segments used in user routes
const (
	PathSegmentMe         = "me"
	PathSegmentRoles      = "roles"
	PathSegmentActivate   = "activate"
	PathSegmentDeactivate = "deactivate"
	PathSegmentSuspend    = "suspend"
)

// Action names used in user context
const (
	ActionView        = "view"
	ActionCreate      = "create"
	ActionUpdate      = "update"
	ActionDelete      = "delete"
	ActionAssignRoles = "assign_roles"
)

// publicRoutesProvider implements middleware.PublicRoutesProvider
type publicRoutesProvider struct{}

func (p *publicRoutesProvider) PublicRoutes() []string {
	return []string{
		"POST " + RouteUsers, // Allow self-service user registration without authentication
	}
}

// authenticatedRoutesProvider implements middleware.AuthenticatedRoutesProvider
type authenticatedRoutesProvider struct{}

func (p *authenticatedRoutesProvider) AuthenticatedRoutes() []string {
	return []string{
		"GET " + RouteUsers + "/" + PathSegmentMe, // Allow authenticated users to get their own information
	}
}

// init registers this module's public and authenticated routes with the middleware registry
func init() {
	middleware.RegisterPublicRoutesProvider(&publicRoutesProvider{})
	middleware.RegisterAuthenticatedRoutesProvider(&authenticatedRoutesProvider{})
}

// @resource code=users service=user-management desc="Core user directory records"
// RegisterHTTPRoutes wires all user HTTP routes into the provided mux.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *UserHandler) {
	// @action name=view desc="Read user records"
	mux.HandleFunc(RouteUsers, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.ListUsers(w, r)
			return
		}
		// @action name=create desc="Create new users"
		handler.CreateUser(w, r)
	})

	// @action name=view desc="Get current authenticated user"
	mux.HandleFunc(RouteUsers+"/"+PathSegmentMe, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetMe(w, r)
			return
		}
		httputil.MethodNotAllowed(w)
	})

	// @action name=view desc="Read a specific user record"
	// @action name=update desc="Update user records"
	// @action name=delete desc="Delete user records"
	// @action name=assign_roles desc="Assign roles to users"
	// @action name=activate desc="Activate user accounts"
	// @action name=deactivate desc="Deactivate user accounts"
	// @action name=suspend desc="Suspend user accounts"
	mux.HandleFunc(RouteUsers+"/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, RouteUsers), "/")
		segments := strings.Split(path, "/")

		if len(segments) >= 2 {
			if handleNestedRoutes(w, r, handler, segments) {
				return
			}
		}

		handleStandardUserRoutes(w, r, handler)
	})
}

// handleNestedRoutes handles nested routes like /users/:id/roles or /users/:id/activate.
// Returns true if a route was handled, false otherwise.
func handleNestedRoutes(w http.ResponseWriter, r *http.Request, handler *UserHandler, segments []string) bool {
	if segments[1] == PathSegmentRoles {
		return handleRoleRoutes(w, r, handler, segments)
	}

	return handleUserActionRoutes(w, r, handler, segments)
}

// handleRoleRoutes handles role-related routes.
func handleRoleRoutes(w http.ResponseWriter, r *http.Request, handler *UserHandler, segments []string) bool {
	switch {
	case len(segments) == 2 && r.Method == http.MethodPost:
		handler.AssignRole(w, r)
		return true
	case len(segments) == 2 && r.Method == http.MethodGet:
		handler.ListUserRoles(w, r)
		return true
	case len(segments) == 3 && r.Method == http.MethodDelete:
		handler.RevokeRole(w, r)
		return true
	default:
		httputil.MethodNotAllowed(w)
		return true
	}
}

// handleUserActionRoutes handles user action routes like activate, deactivate, suspend.
func handleUserActionRoutes(w http.ResponseWriter, r *http.Request, handler *UserHandler, segments []string) bool {
	if len(segments) < 2 {
		return false
	}

	action := segments[1]
	switch action {
	case PathSegmentActivate:
		handler.ActivateUser(w, r)
		return true
	case PathSegmentDeactivate:
		handler.DeactivateUser(w, r)
		return true
	case PathSegmentSuspend:
		handler.SuspendUser(w, r)
		return true
	default:
		return false
	}
}

// handleStandardUserRoutes handles standard CRUD operations on users.
func handleStandardUserRoutes(w http.ResponseWriter, r *http.Request, handler *UserHandler) {
	switch r.Method {
	case http.MethodGet:
		handler.GetUser(w, r)
	case http.MethodPut:
		handler.UpdateUser(w, r)
	case http.MethodDelete:
		handler.DeleteUser(w, r)
	default:
		httputil.MethodNotAllowed(w)
	}
}
