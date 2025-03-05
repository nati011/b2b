package config

import "flag"

type config struct {
	Port                     string
	Env                      string
	KeycloakInstanceURL      string
	KeycloakUsername         string
	KeycloakPassword         string
	KeycloakRealm            string
	KeycloakApplicationRealm string
	KeycloakClientId         string
	DB_URL                   string
	DB_Username              string
	DB_password              string
}

var AppConfig config

func LoadConfig() *config {
	flag.StringVar(&AppConfig.Port, "port", ":8080", "API server port")
	flag.StringVar(&AppConfig.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&AppConfig.KeycloakInstanceURL, "keycloak_base_url", "keycloak Instance Base URL", "Environment (development|staging|production)")
	flag.StringVar(&AppConfig.KeycloakUsername, "keycloak_user_name", "keycloak Instance Base URL", "Environment (development|staging|production)")
	flag.StringVar(&AppConfig.KeycloakPassword, "keycloak_password", "keycloak password", "Environment (development|staging|production)")
	flag.StringVar(&AppConfig.KeycloakRealm, "keycloak_realm", "keycloak realm", "Environment (development|staging|production)")
	flag.StringVar(&AppConfig.KeycloakApplicationRealm, "keycloak_application_realm", "keycloak application realm", "Environment (development|staging|production)")
	flag.StringVar(&AppConfig.KeycloakClientId, "keycloak_client_id", "keycloak ClientId", "Environment (development|staging|production)")
	flag.StringVar(&AppConfig.DB_URL, "db_url", "Database URL", "Environment (development|staging|production)")
	// flag.StringVar(&AppConfig.DB_Username, "db_username", "Database Username", "Environment(development|staging|production)")
	// flag.StringVar(&AppConfig.DB_password, "db_password", "Database Password", "Environment(development|staging|production)")

	flag.Parse()
	return &AppConfig
}
