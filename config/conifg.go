package config

type Config struct {
	//auth
	Port                     int
	Env                      string
	KeycloakInstanceURL      string
	KeycloakUsername         string
	KeycloakPassword         string
	KeycloakRealm            string
	KeycloakApplicationRealm string
	KeycloakClientId         string
	KeycloakClientSecret     string

	//email
	Email string
	SMTP  string

	//db
	FileLocation           string
	CoreDBConnectionString string
}
