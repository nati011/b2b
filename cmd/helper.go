package main

import (
	"errors"

	"b2b.nati011.github.com/config"
)

var (
	ErrPortMissing                               = errors.New("¯\\_(ツ)_/¯, port mandatory")
	ErrEnvMissing                                = errors.New("¯\\_(ツ)_/¯, env mandatory")
	ErrKeycloakInstanceUrlMissing                = errors.New("¯\\_(ツ)_/¯, keycloak instance url mandatory")
	ErrKeycloakUsernameMissing                   = errors.New("¯\\_(ツ)_/¯, keycloak username mandatory")
	ErrKeycloakPasswordMissing                   = errors.New("¯\\_(ツ)_/¯, keycloak password mandatory")
	ErrKeycloakRealmMissing                      = errors.New("¯\\_(ツ)_/¯, keycloak realm mandatory")
	ErrKeycloakApplicationRealmMissing           = errors.New("¯\\_(ツ)_/¯, keycloak application realm mandatory")
	ErrKeycloakClientSecretMissing               = errors.New("¯\\_(ツ)_/¯, keycloak client secret mandatory")
	ErrKeycloakClientIdMissing                   = errors.New("¯\\_(ツ)_/¯, keycloak clientId missing")
	ErrKeycloakClientSecretMandtory              = errors.New("¯\\_(ツ)_/¯, keycloak secret missing")
	ErrEmailMissing                              = errors.New("¯\\_(ツ)_/¯, email missing")
	ErrSMTPMissing                               = errors.New("¯\\_(ツ)_/¯, smtp missing")
	ErrCoreDBConnectionStringMissing             = errors.New("¯\\_(ツ)_/¯, core db conn string missing")
	ErrEmailDBConnectionStringMissing            = errors.New("¯\\_(ツ)_/¯, email db conn string missing")
	ErrMinMobileClientCompatibleVersionMandatory = errors.New("¯\\_(ツ)_/¯, min mobile client version string missing")
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
	if cfg.KeycloakClientSecret == "" {
		panic(ErrKeycloakClientSecretMandtory)
	}
	if cfg.MinMobileClientCompatibleVersion == "" {
		panic(ErrMinMobileClientCompatibleVersionMandatory)
	}
}

/*
	go run cmd/* --port=4001 --db="postgres://user:password@localhost:32769/test?" --keycloak_base_url="https://euc1.auth.ac/auth" --keycloak_user_name="admin" --keycloak_password=">Z3P5k8vHN?bch" --keycloak_realm="b2b" --keycloak_client_id="733bcd1a-dd25-4e59-b14f-3331872a3d4e" --keycloak_application_realm="b2b" --email="test@localhost.com" --smtp=1234
*/
