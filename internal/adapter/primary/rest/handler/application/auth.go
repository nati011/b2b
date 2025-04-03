package handler

import (
	"encoding/json"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/auth"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

type AuthHandler struct {
	service    auth.Provider
	middleware util.AuthMiddleware
}

func InitAuth() {
	handler.Register(new(AuthHandler))
}

func (a *AuthHandler) Init(services *application_core.Container, domainService *domain_core.Container) error {
	a.service = services.AuthService
	a.middleware = *services.AuthMiddleware
	return nil
}

func (a *AuthHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/login", a.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", a.RefreshToken)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RequestErrorResponse(w, err)
		return
	}

	loginResponse, err := h.service.ClientLogin(r.Context(), req)

	if err != nil {
		util.UnauthorizedResponse(w)
		return
	}

	json.NewEncoder(w).Encode(loginResponse)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req auth.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	refreshResponse, err := h.service.RefreshToken(r.Context(), req)
	if err != nil {
		util.UnauthorizedErrorResponse(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(refreshResponse)
}
