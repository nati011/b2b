package main

import (
	"database/sql"
	"errors"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest"
	"b2b.nati011.github.com/internal/core"
)

var (
	ErrFailedToBuildRoutes = errors.New("oopsy, failed to build routes")
)

func InitREST(mux *http.ServeMux, db *sql.DB, services *core.MasterContainer) *http.ServeMux {
	print("Rest Initialized....")
	err := rest.BuildRouter(mux, services)
	if err != nil {
		panic(ErrFailedToBuildRoutes)
	}
	return mux
}
