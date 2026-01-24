package http

import (
	"net/http"

	httputil "marketplace/pkg/http"
	"marketplace/pkg/http/middleware"
)

// HTTP route path
const RouteCustomers = "/customer"

// Resource code
const ResourceCustomers = "customer"

// Action names used in customer context
const (
	ActionView   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// publicRoutesProvider registers public customer routes.
type publicRoutesProvider struct{}

func (p *publicRoutesProvider) PublicRoutes() []string {
	return []string{
		"POST " + RouteCustomers,
	}
}

func init() {
	middleware.RegisterPublicRoutesProvider(&publicRoutesProvider{})
}

// @resource code=customer service=customer-management desc="Customer profile records"
// RegisterHTTPRoutes wires all customer endpoints.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *CustomerHandler) {
	// @action name=view desc="List customer records"
	// @action name=create desc="Create customer records"
	mux.HandleFunc(RouteCustomers, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListCustomers(w, r)
		case http.MethodPost:
			handler.CreateCustomer(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="Read a specific customer record"
	// @action name=update desc="Update a customer record"
	// @action name=delete desc="Delete a customer record"
	mux.HandleFunc(RouteCustomers+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetCustomer(w, r)
		case http.MethodPut:
			handler.UpdateCustomer(w, r)
		case http.MethodDelete:
			handler.DeleteCustomer(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})
}
