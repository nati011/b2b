package auth

import (
	"encoding/json"
	"net/http"

	authDTO "b2b.nati011.github.com/internal/core/domain/auth/model/dto"
	service "b2b.nati011.github.com/internal/core/domain/auth/service"
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

	var req authDTO.LoginUserRequest
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

func (h *AuthHandler) RegisterRetailer(w http.ResponseWriter, r *http.Request) {
	print("In handler")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authDTO.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	registerResponse, err := h.authContainer.AuthService.CreateClient(req)
	if err != nil {
		http.Error(w, "Registration failed", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(registerResponse)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req authDTO.RefreshTokenRequest

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
