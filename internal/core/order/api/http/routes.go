package http

import (
	"net/http"

	httputil "marketplace/pkg/http"
	"marketplace/pkg/http/middleware"
)

// HTTP route paths.
const (
	RouteOrder          = "/order"
	RouteOrdersCustomer = "/orders/customer"
	RouteOrdersSupplier = "/order/supplier"
)

// Resource code.
const ResourceOrder = "order"

// Action names used in order context.
const (
	ActionView   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
)

// authenticatedRoutesProvider registers order routes that should bypass permission checks.
type authenticatedRoutesProvider struct{}

func (p *authenticatedRoutesProvider) AuthenticatedRoutes() []string {
	return []string{
		"GET " + RouteOrder,
		"POST " + RouteOrder,
		"PATCH " + RouteOrder,
		"GET " + RouteOrdersCustomer,
		"GET " + RouteOrdersSupplier,
	}
}

func init() {
	middleware.RegisterAuthenticatedRoutesProvider(&authenticatedRoutesProvider{})
}

// @resource code=order service=order-management desc="Order records and status updates"
// RegisterHTTPRoutes wires all order endpoints.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *OrderHandler) {
	// @action name=view desc="Read order details"
	// @action name=create desc="Create orders"
	// @action name=update desc="Update order status"
	mux.HandleFunc(RouteOrder, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetOrder(w, r)
		case http.MethodPost:
			handler.CreateOrder(w, r)
		case http.MethodPatch:
			// Check if this is a payment status update
			if r.URL.Query().Get("payment_status") != "" {
				handler.UpdatePaymentStatus(w, r)
			} else {
				handler.UpdateOrderStatus(w, r)
			}
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="List customer orders (requires customer_id, optional status filter)"
	mux.HandleFunc(RouteOrdersCustomer, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListCustomerOrders(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="List orders for supplier (automatically filtered by authenticated supplier)"
	mux.HandleFunc(RouteOrdersSupplier, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListSupplierOrders(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})
}
