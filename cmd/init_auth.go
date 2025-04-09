package main

import (
	"context"

	"b2b.nati011.github.com/config"
	keycloak "github.com/stillya/testcontainers-keycloak"
)

const (
	KEYCLOAK_VERSION        = "keycloak/keycloak:24.0"
	KEYCLOAK_ADMIN_USERNAME = "admin"
	KEYCLOAK_ADMIN_PASSWORD = "admin"
	// KEYCLOAK_ADMIN_IMPORTFILE   = "../testdata/realm-export.json"
	KEYCLOAK_ADMIN_CONTEXT_PATH = "/auth"
)

func RunContainer(ctx context.Context) (*keycloak.KeycloakContainer, error) {
	return keycloak.Run(ctx,
		KEYCLOAK_VERSION,
		keycloak.WithContextPath(KEYCLOAK_ADMIN_CONTEXT_PATH),
		// keycloak.WithRealmImportFile(KEYCLOAK_ADMIN_IMPORTFILE),
		keycloak.WithAdminUsername(KEYCLOAK_ADMIN_USERNAME),
		keycloak.WithAdminPassword(KEYCLOAK_ADMIN_PASSWORD),
	)
}

func InitAuth(cfg *config.Config) {
	/*
		DEVELOPMENT:
			Use TestContainers
		STAGING:
			Use External service, verify if connection can be established and operations can be performed
		PRODUCTION:
			Use External service, verify if connection can be established and operations can be performed
	*/

	// switch cfg.Env {
	// case "development":
	// 	InitAuthDevelopment(cfg)
	// case "staging":
	// 	InitAuthStaging(cfg)
	// case "production":
	// 	InitAuthProduction(cfg)
	// }
}

func InitAuthDevelopment(cfg *config.Config) {
	var err error
	var keycloakContainer *keycloak.KeycloakContainer
	ctx := context.Background()
	keycloakContainer, err = RunContainer(ctx)
	if err != nil {
		panic(err)
	}

	keycloakInstanceUrl, err := keycloakContainer.GetAuthServerURL(ctx)
	if err != nil {
		panic(err)
	}

	keycloakAdminClient, err := keycloakContainer.GetAdminClient(ctx)
	if err != nil {
		panic(err)
	}
	cfg.KeycloakUsername = keycloakAdminClient.Username
	cfg.KeycloakPassword = keycloakAdminClient.Password
	cfg.KeycloakRealm = keycloakAdminClient.Realm
	cfg.KeycloakApplicationRealm = keycloakAdminClient.Realm
	cfg.KeycloakClientId = keycloakAdminClient.ClientID
	cfg.KeycloakInstanceURL = keycloakInstanceUrl
}

func InitAuthStaging(cfg *config.Config) {
}

func InitAuthProduction(cfg *config.Config) {
}
