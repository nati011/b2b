package route

import (
	"net/http"

	handler "b2b.nati011.github.com/internal/adapter/primary/rest/handler/application"
)

func CreateResourceRoutes(handler handler.ResourceHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/resource", handler.GetResourceHandler)
	// mux.HandleFunc("POST /api/v1/resource", handler.GetResourceHandler)
	return mux
}
