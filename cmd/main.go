package main

import (
	"flag"

	"b2b.nati011.github.com/config"
)

func main() {
	var cfg config.Config

	//keycloak
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.KeycloakInstanceURL, "keycloak_base_url", "", "Environment (development|staging|production)")
	flag.StringVar(&cfg.KeycloakUsername, "keycloak_user_name", "", "Environment (development|staging|production)")
	flag.StringVar(&cfg.KeycloakPassword, "keycloak_password", "", "Environment (development|staging|production)")
	flag.StringVar(&cfg.KeycloakRealm, "keycloak_realm", "", "Environment (development|staging|production)")
	flag.StringVar(&cfg.KeycloakApplicationRealm, "keycloak_application_realm", "", "Environment (development|staging|production)")
	flag.StringVar(&cfg.KeycloakClientId, "keycloak_client_id", "", "Environment (development|staging|production)")

	//email
	flag.StringVar(&cfg.Email, "email", "", "Environment (development|staging|production)")
	flag.StringVar(&cfg.SMTP, "smtp", "", "Environment (development|staging|production)")

	//db
	flag.StringVar(&cfg.CoreDBConnectionString, "core_db_connection_string", "", "Environment (development|staging|production)")
	flag.Parse()

	connection_pool := InitDB(cfg.CoreDBConnectionString)
	InitREST(connection_pool)
	InitEmail(cfg.Email, cfg.SMTP)
	InitAuth(cfg.Port, cfg.Env, cfg.KeycloakInstanceURL, cfg.KeycloakUsername, cfg.KeycloakPassword, cfg.KeycloakRealm, cfg.KeycloakApplicationRealm, cfg.KeycloakClientId)
	InitSMS()
	// _ := MasterContainer.NewMasterContainer()
}
