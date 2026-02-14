package http

import (
	"net/http"

	httputil "marketplace/pkg/http"
	"marketplace/pkg/http/middleware"
)

// HTTP route paths.
const (
	RouteProduct        = "/product"
	RouteProducts       = "/products"
	RouteCatalogue      = "/catalogue"
	RouteSupplierProducts = "/products/supplier"
)

// Resource code.
const ResourceProducts = "products"

// Action names used in product context.
const (
	ActionView   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

// publicRoutesProvider registers public catalogue routes.
type publicRoutesProvider struct{}

func (p *publicRoutesProvider) PublicRoutes() []string {
	return []string{
		"GET " + RouteCatalogue,  // Browse catalogue
		"GET " + RouteProducts,   // List products
		"GET " + RouteProduct,    // View product details
	}
}

// authenticatedRoutesProvider registers authenticated product routes.
type authenticatedRoutesProvider struct{}

func (p *authenticatedRoutesProvider) AuthenticatedRoutes() []string {
	return []string{
		"GET " + RouteSupplierProducts,        // Supplier products endpoint requires authentication
		"POST " + RouteProduct + "/grn",        // Create GRN requires authentication
		"PATCH " + RouteProduct + "/price",     // Update price requires authentication
	}
}

func init() {
	middleware.RegisterPublicRoutesProvider(&publicRoutesProvider{})
	middleware.RegisterAuthenticatedRoutesProvider(&authenticatedRoutesProvider{})
}

// @resource code=products service=product-management desc="Product catalog records"
// RegisterHTTPRoutes wires all product endpoints.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *ProductHandler) {
	// @action name=view desc="Read product details"
	// @action name=create desc="Create product records"
	// @action name=update desc="Update product records"
	// @action name=delete desc="Delete product records"
	mux.HandleFunc(RouteProduct, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetProduct(w, r)
		case http.MethodPost:
			handler.CreateProduct(w, r)
		case http.MethodPut:
			handler.UpdateProduct(w, r)
		case http.MethodDelete:
			handler.DeleteProduct(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="List products"
	mux.HandleFunc(RouteProducts, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListProducts(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="List catalogue products"
	mux.HandleFunc(RouteCatalogue, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListCatalogue(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="List supplier products (authenticated supplier-only endpoint)"
	mux.HandleFunc(RouteSupplierProducts, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListSupplierProducts(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=update desc="Create goods receiving note (GRN)"
	mux.HandleFunc(RouteProduct+"/grn", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateGRN(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=update desc="Update product price"
	mux.HandleFunc(RouteProduct+"/price", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			handler.UpdatePrice(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})
}
