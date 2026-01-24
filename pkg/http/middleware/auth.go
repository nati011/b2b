package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	authDomain "marketplace/internal/infra/auth"
	userDomain "marketplace/internal/infra/user/domain"
	httputil "marketplace/pkg/http"
	"marketplace/pkg/logger"
)

// UserLoader defines the interface for loading users by ID, email, or phone number.
type UserLoader interface {
	Get(ctx context.Context, id string) (*userDomain.User, error)
	FindByEmail(ctx context.Context, email string) (*userDomain.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*userDomain.User, error)
}

// BasicAuthService defines the interface for basic auth credential validation
type BasicAuthService interface {
	ValidateCredentials(ctx context.Context, username, password string) (*BasicCredential, error)
}

// BasicCredential represents a validated basic auth credential
type BasicCredential struct {
	Username string
	UserID   string
}

const errInvalidCredentials = "invalid credentials"

// AuthMiddleware handles authentication for the configured auth mode.
type AuthMiddleware struct {
	userLoader       UserLoader
	basicAuthService BasicAuthService
	mode             authDomain.AuthMode
	realm            string
	publicRoutes     []string
}

// NewAuthMiddleware creates a new authentication middleware.
func NewAuthMiddleware(mode authDomain.AuthMode, realm string, basicAuthService BasicAuthService, userLoader UserLoader, publicRoutes []string) *AuthMiddleware {
	if realm == "" {
		realm = "Restricted"
	}
	if publicRoutes == nil {
		publicRoutes = []string{}
	}

	return &AuthMiddleware{
		userLoader:       userLoader,
		basicAuthService: basicAuthService,
		mode:             mode,
		realm:            realm,
		publicRoutes:     publicRoutes,
	}
}

// Handle dispatches to the appropriate auth handler based on mode.
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	if m.mode == authDomain.AuthModeBasic {
		return m.handleBasic(next)
	}
	return m.handleUserLookup(next)
}

func (m *AuthMiddleware) handleBasic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers to all responses
		m.setCORSHeaders(w, r)

		// Handle CORS preflight requests (OPTIONS) - skip authentication
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Check if route is public (no authentication required)
		if m.isPublicRoute(r.Method, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if m.basicAuthService == nil {
			m.writeBasicUnauthorized(w, "basic auth is not configured")
			return
		}

		cred, ok := m.validateBasicCredentials(r, w)
		if !ok {
			return
		}

		ctx := m.authenticateUser(r.Context(), cred, w)
		if ctx == nil {
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateBasicCredentials validates the basic auth credentials from the request.
// The username can be:
// 1. A username stored in basic_auth_credentials table
// 2. An email address (if stored as username in basic_auth_credentials)
// 3. A phone number (if stored as username in basic_auth_credentials)
// UserID is mandatory in the credential and must link to an existing user.
// Returns the credential and true if valid, or writes error and returns false.
func (m *AuthMiddleware) validateBasicCredentials(r *http.Request, w http.ResponseWriter) (BasicCredential, bool) {
	username, password, ok := parseBasicAuth(r.Header.Get("Authorization"))
	if !ok {
		m.writeBasicUnauthorized(w, "missing or invalid basic auth header")
		return BasicCredential{}, false
	}

	// Validate credentials using basic auth service (checks database)
	cred, err := m.basicAuthService.ValidateCredentials(r.Context(), username, password)
	if err != nil {
		m.writeBasicUnauthorized(w, errInvalidCredentials)
		return BasicCredential{}, false
	}

	// UserID is mandatory in credential, so we can directly use it
	return BasicCredential{
		Username: cred.Username,
		UserID:   cred.UserID,
	}, true
}

// authenticateUser loads and validates the user if UserID is present.
// Returns the context with user if successful, or nil if authentication failed.
func (m *AuthMiddleware) authenticateUser(ctx context.Context, cred BasicCredential, w http.ResponseWriter) context.Context {
	if cred.UserID == "" || m.userLoader == nil {
		return ctx
	}

	user, err := m.userLoader.Get(ctx, cred.UserID)
	if err != nil {
		m.writeBasicUnauthorized(w, "authentication failed")
		return nil
	}

	if !user.CanLogin() {
		m.writeBasicUnauthorized(w, "user account is not active")
		return nil
	}

	ctx = httputil.WithUser(ctx, user)
	// Add user ID to logger context for request correlation
	ctx = logger.WithUserID(ctx, user.ID)
	return ctx
}

func (m *AuthMiddleware) handleUserLookup(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers to all responses
		m.setCORSHeaders(w, r)

		// Handle CORS preflight requests (OPTIONS) - skip authentication
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Check if route is public (no authentication required)
		if m.isPublicRoute(r.Method, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		userID := extractUserID(r)
		if userID == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		user, err := m.userLoader.Get(ctx, userID)
		if err != nil {
			httputil.JSON(w, http.StatusUnauthorized, map[string]string{
				"error": "authentication failed",
			})
			return
		}

		if !user.CanLogin() {
			httputil.JSON(w, http.StatusUnauthorized, map[string]string{
				"error": "user account is not active",
			})
			return
		}

		ctx = httputil.WithUser(ctx, user)
		// Add user ID to logger context for request correlation
		ctx = logger.WithUserID(ctx, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) writeBasicUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, m.realm))
	httputil.JSON(w, http.StatusUnauthorized, map[string]string{
		"error": message,
	})
}

// setCORSHeaders sets CORS headers to allow cross-origin requests from the frontend
func (m *AuthMiddleware) setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := "*"
	if r != nil {
		if reqOrigin := r.Header.Get("Origin"); reqOrigin != "" {
			origin = reqOrigin
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-User-ID")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "3600")
}

// extractUserID extracts the user ID from headers for non-basic modes.
func extractUserID(r *http.Request) string {
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return strings.TrimSpace(userID)
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return ""
	}

	return token
}

// isPublicRoute checks if the given method and path match any public route.
// Public routes can be specified as:
// - Exact path: "/users" (matches only GET /users)
// - Method-specific: "POST /users" (matches only POST /users)
// - Path prefix: "/public" (matches /public and /public/*)
func (m *AuthMiddleware) isPublicRoute(method, path string) bool {
	for _, publicRoute := range m.publicRoutes {
		// Check for method-specific routes (e.g., "POST /users")
		if strings.HasPrefix(publicRoute, method+" ") {
			routePath := strings.TrimPrefix(publicRoute, method+" ")
			if path == routePath || strings.HasPrefix(path, routePath+"/") {
				return true
			}
		}
		// Check for path-only routes (matches any method)
		if !strings.Contains(publicRoute, " ") {
			if path == publicRoute || strings.HasPrefix(path, publicRoute+"/") {
				return true
			}
		}
	}
	return false
}

func parseBasicAuth(header string) (username, password string, ok bool) {
	if header == "" {
		return "", "", false
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Basic") {
		return "", "", false
	}

	raw, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	creds := strings.SplitN(string(raw), ":", 2)
	if len(creds) != 2 {
		return "", "", false
	}

	return creds[0], creds[1], true
}
