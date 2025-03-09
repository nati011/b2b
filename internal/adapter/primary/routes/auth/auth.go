package routes

import (
	"net/http"

	config "b2b.nati011.github.com/config"
	authHandler "b2b.nati011.github.com/internal/adapter/primary/handlers/auth"
	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
)

func RegisterRoutes(router *http.ServeMux) {
	appConfig := config.AppConfig
	keyCloakProvider := provider.NewKeycloakProvider(
		appConfig.KeycloakInstanceURL,
		appConfig.KeycloakUsername,
		appConfig.KeycloakPassword,
		appConfig.KeycloakRealm,
		appConfig.KeycloakApplicationRealm,
		appConfig.KeycloakClientId,
	)
	handler := authHandler.NewAuthHandler(keyCloakProvider)

	router.HandleFunc("POST /api/v1/auth/login", handler.Login)
	router.HandleFunc("POST /api/v1/auth/refresh", handler.RefreshToken)
}
