package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"marketplace/internal/infra/auth/basic/domain"
	basicauthservice "marketplace/internal/infra/auth/basic/service"
	roleDomain "marketplace/internal/infra/authz/role/domain"
	userDomain "marketplace/internal/infra/user/domain"
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

// LoginRequest represents the request to login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response from login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// LoginResponseWrapper wraps the login response
type LoginResponseWrapper struct {
	Body LoginResponse `json:"body"`
}

// UserRoleService defines the interface for loading user roles
type UserRoleService interface {
	ListRolesForUser(ctx context.Context, userID string) ([]roleDomain.Role, error)
}

// Handler handles HTTP requests for basic auth credential management
type Handler struct {
	service        *basicauthservice.Service
	userLoader     basicauthservice.UserLoader
	userRoleService UserRoleService
	jwtSecret      []byte
}

// NewHandler creates a new basic auth credential handler
func NewHandler(service *basicauthservice.Service, userLoader basicauthservice.UserLoader, userRoleService UserRoleService) *Handler {
	// Use a default secret - in production this should come from config
	secret := []byte("change-me-in-production-secret-key-min-32-chars")
	return &Handler{
		service:        service,
		userLoader:     userLoader,
		userRoleService: userRoleService,
		jwtSecret:      secret,
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

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if !httputil.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	if req.Email == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("email is required"))
		return
	}

	if req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, errors.New("password is required"))
		return
	}

	ctx := r.Context()

	// Validate credentials
	cred, err := h.service.ValidateCredentials(ctx, req.Email, req.Password)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, errors.New("invalid credentials"))
		return
	}

	// Load user
	user, err := h.userLoader.Get(ctx, cred.UserID)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, errors.New("user not found"))
		return
	}

	if !user.CanLogin() {
		httputil.Error(w, http.StatusUnauthorized, errors.New("user account is not active"))
		return
	}

	// Generate JWT token
	accessToken, err := h.generateAccessToken(user)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, errors.New("failed to generate token"))
		return
	}

	// Generate refresh token (simple UUID for now)
	refreshToken := h.generateRefreshToken()

	response := LoginResponseWrapper{
		Body: LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	httputil.JSON(w, http.StatusAccepted, response)
}

// generateAccessToken creates a JWT token for the user
func (h *Handler) generateAccessToken(user *userDomain.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // Token expires in 24 hours

	// Get user roles - load actual role names from the database
	roles := []string{}
	if h.userRoleService != nil {
		ctx := context.Background()
		userRoles, err := h.userRoleService.ListRolesForUser(ctx, user.ID)
		if err == nil && len(userRoles) > 0 {
			// Extract role names (lowercase for consistency with frontend expectations)
			// Frontend checks for role names like 'admin', 'supplier', 'customer', etc.
			for _, role := range userRoles {
				roleName := strings.ToLower(role.Name)
				roles = append(roles, roleName)
			}
		}
	}
	
	// Fallback: if no roles found via service, use simple mapping
	if len(roles) == 0 {
		userTypeStr := strings.ToLower(string(user.UserType))
		if userTypeStr == "officer" {
			roles = append(roles, "officer")
		} else {
			roles = append(roles, "customer")
		}
	}

	claims := jwt.MapClaims{
		"sub":                user.ID,
		"name":               user.Name,
		"email":              user.Email.String(),
		"preferred_username": user.Email.String(),
		"realm_access": map[string]interface{}{
			"roles": roles,
		},
		"iat": now.Unix(),
		"exp": expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

// generateRefreshToken generates a simple refresh token (UUID-based)
func (h *Handler) generateRefreshToken() string {
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type": "refresh",
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
	}).SignedString(h.jwtSecret)
	return token
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
