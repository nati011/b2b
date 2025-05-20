package main

import (
	"database/sql"
	"errors"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

var (
	ErrFailedToBuildRoutes = errors.New("oopsy, failed to build routes")
)

func InitREST(authMiddleWare *util.AuthMiddleware, mux *http.ServeMux, db *sql.DB, applicationServices *application_core.Container, domainServices *domain_core.Container) *http.ServeMux {
	err := rest.BuildRouter(authMiddleWare, mux, applicationServices, domainServices)
	if err != nil {
		panic(ErrFailedToBuildRoutes)
	}
	return mux
}
