package routes

import (
	"database/sql"
	"net/http"

	config "b2b.nati011.github.com/config"
	distributorHandler "b2b.nati011.github.com/internal/adapter/primary/handlers/distributor"
	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/distributor/db"
	service "b2b.nati011.github.com/internal/core/domain/distributor/service"
)

func RegisterRoutes(router *http.ServeMux, db *sql.DB) {
	appConfig := config.AppConfig
	keyCloakProvider := provider.NewKeycloakProvider(
		appConfig.KeycloakInstanceURL,
		appConfig.KeycloakUsername,
		appConfig.KeycloakPassword,
		appConfig.KeycloakRealm,
		appConfig.KeycloakApplicationRealm,
		appConfig.KeycloakClientId,
	)
	distributorService := service.NewDistributorService(

		db_adapter.NewPostgres(
			db,
		),
		keyCloakProvider,
	)
	handler := distributorHandler.NewDistributorHandler(distributorService)
	// handler := authHandler.NewAuthHandler(keyCloakProvider)
	router.HandleFunc("POST /api/auth/register/distributor", handler.RegisterDistributor)
}
