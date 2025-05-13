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

	// Base URL
	BaseUrl     string
	FrontendUrl string

	//email
	Email string
	SMTP  string

	//db
	FileLocation           string
	CoreDBConnectionString string

	// payment partner
	ChapaSecretKey string

	//mobile client version
	MinMobileClientCompatibleVersion string
}
