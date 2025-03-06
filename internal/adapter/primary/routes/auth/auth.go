package routes

import (
	"net/http"

	config "b2b.nati011.github.com/config"
	authHandler "b2b.nati011.github.com/internal/adapter/primary/handlers/auth"
	"b2b.nati011.github.com/pkg/auth/provider"
	service "b2b.nati011.github.com/pkg/auth/service"
)

func RegisterRoutes(router *http.ServeMux) {
	appConfig := config.AppConfig
	container := service.NewContainer(provider.NewKeycloakProvider(appConfig.KeycloakInstanceURL, appConfig.KeycloakUsername, appConfig.KeycloakPassword, appConfig.KeycloakRealm, appConfig.KeycloakApplicationRealm, appConfig.KeycloakClientId))
	handler := authHandler.NewAuthHandler(container)

	router.HandleFunc("/api/auth/login", handler.Login)
	router.HandleFunc("/api/auth/refresh", handler.RefreshToken)
}
