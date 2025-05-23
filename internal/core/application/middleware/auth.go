package middleware

import (
	"context"
	"net/http"
	"strings"

	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	auth "b2b.nati011.github.com/internal/core/application/auth"
)

type Auth struct {
	auth auth.Provider
}

type Option func(*Auth)

func NewAuthMiddleware(auth_service auth.Provider) *Auth {
	return &Auth{
		auth: auth_service,
	}
}

func WithRole(roles []string) Option {
	return func(a *Auth) {

	}
}

func (am *Auth) RequireAuthentication(next http.Handler, options ...Option) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			util.UnauthorizedResponse(w)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.UnauthorizedResponse(w)
			return
		}

		token := parts[1]
		if strings.TrimSpace(token) == "" {
			util.UnauthorizedResponse(w)
			return
		}

		result, err := am.auth.RetrospectToken(r.Context(), token)
		if err != nil {
			util.UnauthorizedResponse(w)
			return
		}

		if !result.Active {
			util.UnauthorizedResponse(w)
			return
		}

		decodedToken, err := am.auth.DecodeToken(
			r.Context(),
			token)

		claims := decodedToken.Claims

		if err != nil {
			util.UnauthorizedResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *Auth) RequireNoAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
