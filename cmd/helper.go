package main

import (
	"errors"

	"b2b.nati011.github.com/config"
)

var (
	ErrPortMissing                     = errors.New("oopsy, port mandatory")
	ErrEnvMissing                      = errors.New("oopsy, env mandatory")
	ErrKeycloakInstanceUrlMissing      = errors.New("oopsy, keycloak instance url mandatory")
	ErrKeycloakUsernameMissing         = errors.New("oopsy, keycloak username mandatory")
	ErrKeycloakPasswordMissing         = errors.New("oopsy, keycloak password mandatory")
	ErrKeycloakRealmMissing            = errors.New("oopsy, keycloak realm mandatory")
	ErrKeycloakApplicationRealmMissing = errors.New("oopsy, keycloak application realm mandatory")
	ErrKeycloakClientIdMissing         = errors.New("oopsy, keycloak clientId missing")
	ErrEmailMissing                    = errors.New("oopsy, email missing")
	ErrSMTPMissing                     = errors.New("oopsy, smtp missing")
	ErrCoreDBConnectionStringMissing   = errors.New("oopsy, core db conn string missing")
	ErrEmailDBConnectionStringMissing  = errors.New("oopsy, email db conn string missing")
)

func validateFlags(cfg config.Config) {
	if cfg.Port == 0 {
		panic(ErrPortMissing)
	}
	if cfg.Env == "" {
		panic(ErrEnvMissing)
	}
	if cfg.KeycloakInstanceURL == "" {
		panic(ErrKeycloakInstanceUrlMissing)
	}
	if cfg.KeycloakUsername == "" {
		panic(ErrKeycloakUsernameMissing)
	}
	if cfg.KeycloakPassword == "" {
		panic(ErrKeycloakPasswordMissing)
	}
	if cfg.KeycloakRealm == "" {
		panic(ErrKeycloakRealmMissing)
	}
	if cfg.KeycloakApplicationRealm == "" {
		panic(ErrKeycloakApplicationRealmMissing)
	}
	if cfg.KeycloakClientId == "" {
		panic(ErrKeycloakClientIdMissing)
	}
	if cfg.Email == "" {
		panic(ErrEmailMissing)
	}
	if cfg.SMTP == "" {
		panic(ErrSMTPMissing)
	}
}

/*
	go run cmd/* --port=4001 --db="postgres://user:password@localhost:32769/test?" --keycloak_base_url="https://euc1.auth.ac/auth" --keycloak_user_name="admin" --keycloak_password=">Z3P5k8vHN?bch" --keycloak_realm="b2b" --keycloak_client_id="733bcd1a-dd25-4e59-b14f-3331872a3d4e" --keycloak_application_realm="b2b" --email="test@localhost.com" --smtp=1234
*/
