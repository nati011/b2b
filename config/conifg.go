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

	//email
	Email string
	SMTP  string

	//db
	CoreDBConnectionString  string
	EmailDBConnectionString string
}
