package currency

import (
	"net/http"

	httputil "core/pkg/http"
)

const (
	// RouteCurrencies is the base route for currency operations
	RouteCurrencies = "/organization/currencies"
	// ResourceCurrencies is the resource code for currencies
	ResourceCurrencies = "currencies"
)

// @resource code=currencies service=organization-management desc="Currency configuration"
// RegisterHTTPRoutes wires all currency HTTP routes into the provided mux.
func RegisterHTTPRoutes(mux *http.ServeMux, handler *CurrencyHandler) {
	// @action name=view desc="Get the configured currency"
	// @action name=create desc="Configure currency during registration (one-time only)"
	mux.HandleFunc(RouteCurrencies, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetCurrency(w, r)
		case http.MethodPost:
			handler.CreateCurrency(w, r)
		default:
			httputil.MethodNotAllowed(w)
		}
	})

	// @action name=view desc="Get currency dropdown options"
	mux.HandleFunc(RouteCurrencies+"/dropdown", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetCurrencyDropdown(w, r)
			return
		}
		httputil.MethodNotAllowed(w)
	})
}
