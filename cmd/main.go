package main

import "flag"

type config struct {
	port                     int
	env                      string
	keycloakInstanceURL      string
	keycloakUsername         string
	keycloakPassword         string
	keycloakRealm            string
	keycloakApplicationRealm string
	keycloakClientId         string
	keycloakClientSecret     string
}

func main() {
	print("hello")
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakInstanceURL, "keycloak_base_url", "keycloak Instance Base URL", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakUsername, "keycloak_user_name", "keycloak Instance Base URL", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakPassword, "keycloak_password", "keycloak password", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakRealm, "keycloak_realm", "keycloak realm", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakApplicationRealm, "keycloak_application_realm", "keycloak application realm", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakClientId, "keycloak_client_id", "keycloak ClientId", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakClientSecret, "keycloak_client_secret", "keycloak ClientSecret", "Environment (development|staging|production)")

	flag.Parse()

	//create test containers
	_ = MasterTestContainer.NewMasterTestContainer()
}
