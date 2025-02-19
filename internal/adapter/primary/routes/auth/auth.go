package auth

import (
	"net/http"

	config "b2b.nati011.github.com/config"
	authHandler "b2b.nati011.github.com/internal/adapter/primary/handlers/auth"
	"b2b.nati011.github.com/internal/core/domain/auth/provider"
	service "b2b.nati011.github.com/internal/core/domain/auth/service"
)

func RegisterRoutes(router *http.ServeMux) {
	appConfig := config.AppConfig
	container := service.NewContainer(provider.NewKeycloakProvider(appConfig.KeycloakClientId, appConfig.KeycloakInstanceURL, appConfig.KeycloakRealm, appConfig.KeycloakUsername, appConfig.KeycloakPassword, appConfig.KeycloakApplicationRealm))
	handler := authHandler.NewAuthHandler(container)

	router.HandleFunc("/api/auth/login", handler.Login)
	// router.HandleFunc("/api/auth/logout", handler.Logout)
	router.HandleFunc("/api/auth/refresh", handler.RefreshToken)
	router.HandleFunc("/api/auth/register/retailer", handler.RegisterRetailer)
}
