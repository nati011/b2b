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
			UnauthorizedErrorResponse(w, r, ErrUnAuthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			UnauthorizedErrorResponse(w, r, ErrUnAuthorized)
			return
		}

		token := parts[1]
		if token == "" {
			UnauthorizedErrorResponse(w, r, ErrUnAuthorized)
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
			UnauthorizedErrorResponse(w, r, ErrUnAuthorized)
			return
		}

		if !*result.Active {
			UnauthorizedErrorResponse(w, r, ErrUnAuthorized)
			return
		}

		decodedToken, mapClaims, err := am.client.DecodeAccessToken(
			r.Context(),
			token,
			am.Realm,
		)
		print(decodedToken)
		print(mapClaims)
		if err != nil {
			UnauthorizedErrorResponse(w, r, ErrUnAuthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "auth_info", result)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
