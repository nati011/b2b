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
		"POST " + RouteSuppliers,        // Register new supplier
		"GET " + RouteSuppliers,         // List suppliers (public browsing)
		"GET " + RouteSuppliers + "/",   // View supplier details (public)
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

	// Bank account routes
	// @resource code=bank-account service=supplier-management desc="Supplier bank account records for payment processing"
	mux.HandleFunc(RouteSuppliers+"/bank-account", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// @action name=create desc="Create a bank account for a supplier"
			handler.CreateBankAccount(w, r)
		case http.MethodGet:
			// @action name=view desc="List bank accounts for a supplier"
			handler.ListBankAccounts(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	mux.HandleFunc(RouteSuppliers+"/bank-account/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// @action name=update desc="Update a bank account"
			handler.UpdateBankAccount(w, r)
		case http.MethodDelete:
			// @action name=delete desc="Delete a bank account"
			handler.DeleteBankAccount(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})
}

