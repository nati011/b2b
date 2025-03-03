package routes

import (
	"net/http"

	distributorHandler "b2b.nati011.github.com/internal/adapter/primary/handlers/distributor"
	db "b2b.nati011.github.com/internal/adapter/secondary/distributor/db"
	service "b2b.nati011.github.com/internal/core/domain/distributor/service"
)

func RegisterRoutes(router *http.ServeMux) {
	service := service.NewDistributorService(db.NewPostgres())
	handler := distributorHandler.NewDistributorHandler(service)

	router.HandleFunc("/api/auth/login", handler.Login)
	// router.HandleFunc("/api/auth/logout", handler.Logout)
	router.HandleFunc("/api/auth/refresh", handler.RefreshToken)
	router.HandleFunc("/api/auth/register/distributor", handler.RegisterDistributor)
}
