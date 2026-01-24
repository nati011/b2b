package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"marketplace/internal/infra/auth/basic/domain"
	basicauthservice "marketplace/internal/infra/auth/basic/service"
	httputil "marketplace/pkg/http"
)

const (
	errUsernameRequired = "username is required"
	errPasswordRequired = "password is required"
)

// CreateCredentialRequest represents the request to create a basic auth credential
type CreateCredentialRequest struct {
	Username          string `json:"username"`                     // Can be email, phone number, or custom username
	Password          string `json:"password"`                     // Plain text password (will be hashed)
	UserID            string `json:"user_id,omitempty"`            // Optional: for authenticated officer users creating credentials for others
	RegistrationToken string `json:"registration_token,omitempty"` // Required for self-service users: token from user creation
}

// UpdatePasswordRequest represents the request to update a password
type UpdatePasswordRequest struct {
	Password string `json:"password"` // New plain text password (will be hashed)
}

// CredentialResponse represents a credential in API responses
type CredentialResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	UserID    string    `json:"user_id"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToCredentialResponse converts a domain Credential to a CredentialResponse DTO
func ToCredentialResponse(cred *domain.Credential) CredentialResponse {
	return CredentialResponse{
		ID:        cred.ID,
		Username:  cred.Username,
		UserID:    cred.UserID,
		Active:    cred.Active,
		CreatedAt: cred.CreatedAt,
		UpdatedAt: cred.UpdatedAt,
	}
}

// Handler handles HTTP requests for basic auth credential management
type Handler struct {
	service *basicauthservice.Service
}

// NewHandler creates a new basic auth credential handler
func NewHandler(service *basicauthservice.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateCredential handles POST /auth/basic/credentials
func (h *Handler) CreateCredential(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	// Validate required fields
	if req.Username == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New(errUsernameRequired))
		return
	}

	if req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New(errPasswordRequired))
		return
	}

	ctx := r.Context()
	serviceReq := basicauthservice.CreateCredentialParams{
		Username: req.Username,
		Password: req.Password,
		UserID:   req.UserID,
	}

	cred, err := h.service.CreateCredentialWithAuthorization(ctx, serviceReq)
	if err != nil {
		h.writeError(w, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, ToCredentialResponse(cred))
}

// UpdatePassword handles PUT /auth/basic/credentials/:username/password
func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPut) {
		return
	}

	username := h.extractUsername(r.URL.Path)
	if username == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New(errUsernameRequired))
		return
	}

	var req UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	if req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New(errPasswordRequired))
		return
	}

	ctx := r.Context()
	if err := h.service.UpdatePasswordWithAuthorization(ctx, username, req.Password); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeactivateCredential handles POST /auth/basic/credentials/:username/deactivate
func (h *Handler) DeactivateCredential(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	username := h.extractUsername(r.URL.Path)
	if username == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New(errUsernameRequired))
		return
	}

	ctx := r.Context()
	if err := h.service.DeactivateCredentialWithAuthorization(ctx, username); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	status := h.mapError(err)
	httputil.Error(w, status, err)
}

func (h *Handler) mapError(err error) int {
	if errors.Is(err, domain.ErrCredentialNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, basicauthservice.ErrTargetUserNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, basicauthservice.ErrRegistrationTokenRequired) ||
		errors.Is(err, basicauthservice.ErrUserIDRequired) {
		return http.StatusBadRequest
	}
	if errors.Is(err, basicauthservice.ErrSelfServiceOwnCredentialsOnly) ||
		errors.Is(err, basicauthservice.ErrInsufficientPermissions) ||
		errors.Is(err, basicauthservice.ErrRegistrationTokenForSelfServiceOnly) {
		return http.StatusForbidden
	}
	if errors.Is(err, basicauthservice.ErrAuthenticationRequired) {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}

// extractUsername extracts the username from the URL path
// Expected format: /auth/basic/credentials/:username/...
func (h *Handler) extractUsername(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Path format: /auth/basic/credentials/:username/...
	if len(parts) >= 4 && parts[0] == PathSegmentAuth && parts[1] == PathSegmentBasic && parts[2] == PathSegmentCredentials {
		return parts[3]
	}
	return ""
}
