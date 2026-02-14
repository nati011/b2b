package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	httputil "marketplace/pkg/http"
)

// RegistrationTokenValidator defines the interface for validating registration tokens
type RegistrationTokenValidator interface {
	ValidateAndConsume(ctx context.Context, token string) (userID string, err error)
}

// RegistrationTokenMiddleware validates registration tokens from request bodies
// and sets the userID in context for handlers to use.
type RegistrationTokenMiddleware struct {
	validator RegistrationTokenValidator
	paths     []string // Paths that should be checked for registration tokens
}

// NewRegistrationTokenMiddleware creates a new registration token middleware
func NewRegistrationTokenMiddleware(validator RegistrationTokenValidator, paths []string) *RegistrationTokenMiddleware {
	return &RegistrationTokenMiddleware{
		validator: validator,
		paths:     paths,
	}
}

// Handle processes requests and validates registration tokens when present
func (m *RegistrationTokenMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only process specified paths
		if !m.shouldProcess(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Only process POST requests (registration tokens are used for credential creation)
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		// Read request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			httputil.Error(w, http.StatusBadRequest, err)
			return
		}
		_ = r.Body.Close()

		// Restore body for handler
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		// Parse request to check for registration token
		var reqBody map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
			// If body is not JSON or malformed, let handler deal with it
			next.ServeHTTP(w, r)
			return
		}

		// Check if registration token is present
		token, ok := reqBody["registration_token"].(string)
		if !ok || token == "" {
			// No registration token, proceed normally
			next.ServeHTTP(w, r)
			return
		}

		// Validate registration token
		if m.validator == nil {
			httputil.Error(w, http.StatusInternalServerError, http.ErrBodyReadAfterClose)
			return
		}

		userID, err := m.validator.ValidateAndConsume(r.Context(), token)
		if err != nil {
			// Return error - handler or validator should handle specific error types
			httputil.Error(w, http.StatusBadRequest, errors.New("invalid or expired registration token"))
			return
		}

		// Set userID in context for handler to use
		ctx := httputil.WithRegistrationUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// shouldProcess checks if the path should be processed for registration tokens
func (m *RegistrationTokenMiddleware) shouldProcess(path string) bool {
	for _, p := range m.paths {
		if path == p || (len(path) > len(p) && path[:len(p)] == p) {
			return true
		}
	}
	return false
}
