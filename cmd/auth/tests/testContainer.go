package test

import (
	"context"
	"log"

	auth "b2b.nati011.github.com/cmd/auth/service"
	keycloak "github.com/stillya/testcontainers-keycloak"
)

// var authTestContainer = &AuthTestContainer{
// 	authService: &auth.AuthService{},
// 	keycloak:    *keycloak.KeycloakContainer,
// }

type AuthTestContainer struct {
	authService *auth.AuthService
	keycloak    *keycloak.KeycloakContainer
}

func (AuthTestContainer) NewTestContainer() (*AuthTestContainer, error) {
	tc := AuthTestContainer{}
	initTestContainer(&tc)
	initAuthService(&tc)
	return &tc, nil
}

func initTestContainer(tc *AuthTestContainer) {
	ctx := context.Background()
	keycloak, err := CreateKeycloakContainer(ctx)
	if err != nil {
		log.Fatalf("failed to create keycloak conatiner err: %v", err)
	}
	tc.keycloak = keycloak
}

func initAuthService(tc *AuthTestContainer) {
	authService, err := auth.NewAuthService()
	if err != nil {
		log.Fatalf("failed to create auth service err: %v", err)
	}
	tc.authService = authService
}

func CreateKeycloakContainer(ctx context.Context) (*keycloak.KeycloakContainer, error) {
	return keycloak.Run(ctx,
		"keycloak/keycloak:24.0",
		keycloak.WithContextPath("/auth"),
		keycloak.WithRealmImportFile("../testdata/realm-export.json"),
		keycloak.WithAdminUsername("admin"),
		keycloak.WithAdminPassword("admin"),
	)
}
