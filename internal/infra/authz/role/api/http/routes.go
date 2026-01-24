package role

import (
	"net/http"

	httputil "marketplace/pkg/http"
)

// HTTP route path
const RouteRoles = "/roles"

// Resource code
const ResourceRoles = "roles"

// Action names used in role context
const (
	ActionView   = "view"
	ActionManage = "manage"
	ActionDelete = "delete"
)

// @resource code=roles service=user-management desc="RBAC role definitions"
// RegisterHTTPRoutes wires all role HTTP routes into the provided mux.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *RoleHandler) {
	// @action name=view desc="Read roles"
	// @action name=manage desc="Create or update roles"
	mux.HandleFunc(RouteRoles, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListRoles(w, r)
		case http.MethodPost:
			handler.CreateRole(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="Read a specific role"
	// @action name=manage desc="Update a specific role"
	// @action name=delete desc="Delete a role"
	mux.HandleFunc(RouteRoles+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetRole(w, r)
		case http.MethodPut:
			handler.UpdateRole(w, r)
		case http.MethodDelete:
			handler.DeleteRole(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})
}
