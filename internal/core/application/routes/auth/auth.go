package auth

import (
	"net/http"

	authHandler "b2b.nati011.github.com/internal/core/application/handlers/auth"
	service "b2b.nati011.github.com/internal/core/domain/auth/service"
)

func RegisterRoutes(router *http.ServeMux, authService *service.AuthService) {
	handler := authHandler.NewAuthHandler(authService)

	router.HandleFunc("/api/auth/login", handler.Login)
	// router.HandleFunc("/api/auth/logout", handler.Logout)
	// router.HandleFunc("/api/auth/refresh", handler.RefreshToken)
	// router.HandleFunc("/api/auth/validate", handler.ValidateToken)
}
