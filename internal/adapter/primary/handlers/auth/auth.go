package handlers

import (
	"encoding/json"
	"net/http"

	service "b2b.nati011.github.com/internal/core/application/auth"
	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

type AuthHandler struct {
	authContainer *service.Container
}

func NewAuthHandler(authContainer *service.Container) *AuthHandler {
	return &AuthHandler{
		authContainer: authContainer,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req port.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	loginResponse, err := h.authContainer.AuthService.LoginClient(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(loginResponse)
}

func (h *AuthHandler) RegisterDistributor(w http.ResponseWriter, r *http.Request) {
	print("In handler")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req port.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	registerResponse, err := h.authContainer.AuthService.CreateClient(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(registerResponse)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req port.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	refreshResponse, err := h.authContainer.AuthService.RefreshToken(req)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(refreshResponse)
}
