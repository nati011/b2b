package main

import (
	"errors"

	"b2b.nati011.github.com/config"
)

var (
	ErrPortMissing                               = errors.New(" port mandatory")
	ErrEnvMissing                                = errors.New(" env mandatory")
	ErrKeycloakInstanceUrlMissing                = errors.New(" keycloak instance url mandatory")
	ErrKeycloakUsernameMissing                   = errors.New(" keycloak username mandatory")
	ErrKeycloakPasswordMissing                   = errors.New(" keycloak password mandatory")
	ErrKeycloakRealmMissing                      = errors.New(" keycloak realm mandatory")
	ErrKeycloakApplicationRealmMissing           = errors.New(" keycloak application realm mandatory")
	ErrKeycloakClientSecretMissing               = errors.New(" keycloak client secret mandatory")
	ErrKeycloakClientIdMissing                   = errors.New(" keycloak clientId missing")
	ErrKeycloakClientSecretMandtory              = errors.New(" keycloak secret missing")
	ErrEmailMissing                              = errors.New(" email missing")
	ErrSMTPMissing                               = errors.New(" smtp missing")
	ErrCoreDBConnectionStringMissing             = errors.New(" core db conn string missing")
	ErrEmailDBConnectionStringMissing            = errors.New(" email db conn string missing")
	ErrMinMobileClientCompatibleVersionMandatory = errors.New(" min mobile client version string missing")
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
