package routes

import (
	"database/sql"
	"net/http"

	distributorHandler "b2b.nati011.github.com/internal/adapter/primary/handlers/distributor"
	service "b2b.nati011.github.com/internal/core/domain/distributor/service"
)

func RegisterRoutes(router *http.ServeMux, db *sql.DB) {
	container := service.NewContainer(
		db,
	)
	handler := distributorHandler.NewDistributorHandler(container)
	router.HandleFunc("/api/auth/register/distributor", handler.RegisterDistributor)
}
