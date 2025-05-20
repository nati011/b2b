package handler

import (
	"context"
	"net/http"
	"strings"

	"b2b.nati011.github.com/internal/core/application/auth"
)

type AuthMiddleware struct {
	auth auth.Provider
}

func NewAuthMiddleware(auth_service auth.Provider) *AuthMiddleware {
	return &AuthMiddleware{
		auth: auth_service,
	}
}

func (am *AuthMiddleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			UnauthorizedResponse(w)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			UnauthorizedResponse(w)
			return
		}

		token := parts[1]
		if strings.TrimSpace(token) == "" {
			UnauthorizedResponse(w)
			return
		}

		result, err := am.auth.RetrospectToken(r.Context(), token)
		if err != nil {
			UnauthorizedResponse(w)
			return
		}

		if !*result.Active {
			UnauthorizedResponse(w)
			return
		}

		decodedToken, err := am.auth.DecodeToken(
			r.Context(),
			token)

		claims := decodedToken.Claims

		if err != nil {
			UnauthorizedResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) RequireNoAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
