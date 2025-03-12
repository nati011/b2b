package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"b2b.nati011.github.com/internal/core/application/auth"
	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

type AuthHandler struct {
	service port.Provider
}

func NewAuthHandler(provider port.Provider) *AuthHandler {
	return &AuthHandler{
		service: auth.NewAuthService(provider),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req port.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	loginResponse, err := h.service.ClientLogin(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(loginResponse)
}

func (h *AuthHandler) RegisterDistributor(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req port.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	registerResponse, err := h.service.CreateNewClient(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(registerResponse)
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
