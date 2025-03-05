package routes

import (
	"net/http"
)

func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}
