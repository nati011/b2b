package handler

import (
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
)

type Middleware struct {
	URL            string
	clientID       string
	clientSecret   string
	realm          string
	clientPassword string
}

func NewMiddleware(keycloak_base_url, keycloak_password, keycloak_realm, keycloak_client_id, keycloak_client_secret string) *Middleware {
	middleware := Middleware{}
	middleware.Init(
		keycloak_base_url, keycloak_password, keycloak_realm, keycloak_client_id, keycloak_client_secret,
	)
	return &middleware
}

func (a *Middleware) Init(keycloak_base_url, keycloak_password, keycloak_realm, keycloak_client_id, keycloak_client_secret string) {
	a.URL = keycloak_base_url
	a.clientID = keycloak_client_id
	a.clientSecret = keycloak_client_secret
	a.realm = keycloak_realm
	a.clientPassword = keycloak_password
}

func (a *Middleware) AuthenticationMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		fmt.Printf("Header %v", header)
		client := gocloak.NewClient(a.URL)

		if header.Get("Authorization") == "" {
			fmt.Printf("Authorization Header empty: %v", header)
		}

		result, err := client.RetrospectToken(r.Context(), header.Get("Authorization"), a.clientID, a.clientSecret, a.realm)
		if err != nil {
			print(header.Get("Authorization"))
			print(err.Error())
		}

		fmt.Printf("Introspection result %v", result)
		handler.ServeHTTP(w, r)
	})
}
