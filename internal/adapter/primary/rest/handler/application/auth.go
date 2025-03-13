package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"b2b.nati011.github.com/internal/adapter/primary/rest/handler"
	"b2b.nati011.github.com/internal/core"
	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

type AuthHandler struct {
	service port.Provider
}

func InitAuth() {
	handler.Register(new(AuthHandler))
}

func (a *AuthHandler) Init(services *core.MasterContainer) error {
	print("Registred Services")
	print(services)
	a.service = services.AuthService
	return nil
}

func (a *AuthHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/login", a.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", a.RefreshToken)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	print("Ezi.....")
	ctx := context.Background()
	var req port.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	print(h.service)
	loginResponse, err := h.service.ClientLogin(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(loginResponse)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req port.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	refreshResponse, err := h.service.RefreshToken(ctx, req)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(refreshResponse)
}
