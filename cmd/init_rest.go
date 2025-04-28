package main

import (
	"database/sql"
	"errors"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var (
	ErrFailedToBuildRoutes = errors.New("oopsy, failed to build routes")
)

func InitREST(mux *http.ServeMux, db *sql.DB, applicationServices *application_core.Container, domainServices *domain_core.Container) *http.ServeMux {
	err := rest.BuildRouter(mux, applicationServices, domainServices)
	if err != nil {
		panic(ErrFailedToBuildRoutes)
	}

	corsMux := http.NewServeMux()
	corsMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		mux.ServeHTTP(w, r)
	})

	return corsMux
}
