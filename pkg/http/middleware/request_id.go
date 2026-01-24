package middleware

import (
	"marketplace/pkg/logger"
	"net/http"

	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// RequestIDMiddleware adds a request ID to each request context
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get request ID from header first
		requestID := r.Header.Get(requestIDHeader)

		// Generate new request ID if not present
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add to response header
		w.Header().Set(requestIDHeader, requestID)

		// Add to context
		ctx := logger.WithRequestID(r.Context(), requestID)

		// Continue with enriched context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
