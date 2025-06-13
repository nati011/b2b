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
	ErrFailedToBuildRoutes = errors.New("¯\\_(o_o)_/¯, failed to build routes")
)

func InitREST(mux *http.ServeMux, db *sql.DB, applicationServices *application_core.Container, domainServices *domain_core.Container) *http.ServeMux {
	err := rest.BuildRouter(mux, applicationServices, domainServices)
	if err != nil {
		panic(ErrFailedToBuildRoutes)
	}
	return mux
}
