package routes

import (
	"net/http"

	"b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"

	config "b2b.nati011.github.com/config"
	authHandler "b2b.nati011.github.com/internal/adapter/primary/handlers/auth"
	service "b2b.nati011.github.com/internal/core/application/auth"
)

func RegisterRoutes(router *http.ServeMux) {
	appConfig := config.AppConfig
	keyCloakProvider := provider.NewKeycloakProvider(appConfig.KeycloakInstanceURL, appConfig.KeycloakUsername, appConfig.KeycloakPassword, appConfig.KeycloakRealm, appConfig.KeycloakApplicationRealm, appConfig.KeycloakClientId)
	container := service.NewContainer(keyCloakProvider)
	handler := authHandler.NewAuthHandler(container)

	router.HandleFunc("/api/auth/login", handler.Login)
	router.HandleFunc("/api/auth/refresh", handler.RefreshToken)
}
