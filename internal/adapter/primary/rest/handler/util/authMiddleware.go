package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
)

type AuthMiddleware struct {
	client       *gocloak.GoCloak
	BaseURL      string
	ClientID     string
	ClientSecret string
	Realm        string
	Password     string
}

var (
	ErrUnAuthorized = errors.New("oopsy, unauthorized user")
)

func NewAuthMiddleware(
	BaseURL string,
	ClientID string,
	ClientSecret string,
	Realm string,
	Password string,
) *AuthMiddleware {
	return &AuthMiddleware{
		client:       gocloak.NewClient(BaseURL),
		BaseURL:      BaseURL,
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Realm:        Realm,
		Password:     Password,
	}
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			UnauthorizedResponse(w)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			UnauthorizedResponse(w)
			return
		}

		token := parts[1]
		if token == "" {
			UnauthorizedResponse(w)
			return
		}

		result, err := am.client.RetrospectToken(
			r.Context(),
			token,
			am.ClientID,
			am.ClientSecret,
			am.Realm,
		)
		if err != nil {
			UnauthorizedResponse(w)
			return
		}

		if !*result.Active {
			UnauthorizedResponse(w)
			return
		}

		ctx := context.WithValue(r.Context(), "auth_info", result)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
