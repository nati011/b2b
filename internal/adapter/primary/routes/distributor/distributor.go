package routes

import (
	"net/http"

	distributorHandler "b2b.nati011.github.com/internal/adapter/primary/handlers/distributor"
	service "b2b.nati011.github.com/internal/core/domain/distributor/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *http.ServeMux, dbConnection *pgxpool.Pool) {
	container := service.NewContainer(
		dbConnection,
	)
	handler := distributorHandler.NewDistributorHandler(container)
	router.HandleFunc("/api/auth/register/distributor", handler.RegisterDistributor)
}
