package http

import (
	"net/http"

	httputil "marketplace/pkg/http"
	"marketplace/pkg/http/middleware"
)

// HTTP route path
const RouteSuppliers = "/supplier"

// Resource code
const ResourceSuppliers = "supplier"

// Action names used in supplier context
const (
	ActionView   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// publicRoutesProvider registers public supplier routes.
type publicRoutesProvider struct{}

func (p *publicRoutesProvider) PublicRoutes() []string {
	return []string{
		"POST " + RouteSuppliers,
	}
}

func init() {
	middleware.RegisterPublicRoutesProvider(&publicRoutesProvider{})
}

// @resource code=supplier service=supplier-management desc="Supplier profile records"
// RegisterHTTPRoutes wires all supplier endpoints.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *SupplierHandler) {
	// @action name=view desc="List supplier records"
	// @action name=create desc="Create supplier records"
	mux.HandleFunc(RouteSuppliers, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListSuppliers(w, r)
		case http.MethodPost:
			handler.CreateSupplier(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="Read a specific supplier record"
	// @action name=update desc="Update a supplier record"
	// @action name=delete desc="Delete a supplier record"
	mux.HandleFunc(RouteSuppliers+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetSupplier(w, r)
		case http.MethodPut:
			handler.UpdateSupplier(w, r)
		case http.MethodDelete:
			handler.DeleteSupplier(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})
}

