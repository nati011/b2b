package http

import (
	"net/http"

	httputil "marketplace/pkg/http"
)

// HTTP route path
const RoutePermissions = "/permissions"

// Resource code
const ResourcePermissions = "permissions"

// Action names used in permission context
const (
	ActionView   = "view"
	ActionManage = "manage"
	ActionDelete = "delete"
)

// @resource code=permissions service=user-management desc="Permission catalog definitions"
// RegisterHTTPRoutes wires all permission endpoints.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *PermissionHandler) {
	// @action name=view desc="Read permissions"
	// @action name=manage desc="Create or update permissions"
	mux.HandleFunc(RoutePermissions, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListPermissions(w, r)
		case http.MethodPost:
			handler.CreatePermission(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="Read a specific permission"
	// @action name=manage desc="Update a specific permission"
	// @action name=delete desc="Delete a permission"
	mux.HandleFunc(RoutePermissions+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetPermission(w, r)
		case http.MethodPut:
			handler.UpdatePermission(w, r)
		case http.MethodDelete:
			handler.DeletePermission(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})
}
