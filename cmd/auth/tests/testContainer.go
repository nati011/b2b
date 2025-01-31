package test

import (
	"context"

	auth "b2b.nati011.github.com/cmd/auth/service"
	keycloak "github.com/stillya/testcontainers-keycloak"
)

var authTestContainer = &AuthTestContainer{
	authService: &auth.AuthService{},
}

type AuthTestContainer struct {
	authService *auth.AuthService
}

func (AuthTestContainer) NewTestContainer() (*AuthTestContainer, error) {
	tc := AuthTestContainer{}
	err := initTestContainer(&tc)
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func initTestContainer(tc *AuthTestContainer) error {
	ctx := context.Background()
	CreateKeycloakContainer(ctx)
	return nil
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
