package main

import (
	"flag"
	"log"
	"net/http"

	authRouter "b2b.nati011.github.com/internal/core/application/routes/auth"
	authService "b2b.nati011.github.com/internal/core/domain/auth/service"
)

type config struct {
	port                     string
	env                      string
	keycloakInstanceURL      string
	keycloakUsername         string
	keycloakPassword         string
	keycloakRealm            string
	keycloakApplicationRealm string
	keycloakClientId         string
}

func main() {
	var cfg config

	flag.StringVar(&cfg.port, "port", ":8080", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakInstanceURL, "keycloak_base_url", "keycloak Instance Base URL", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakUsername, "keycloak_user_name", "keycloak Instance Base URL", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakPassword, "keycloak_password", "keycloak password", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakRealm, "keycloak_realm", "keycloak realm", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakApplicationRealm, "keycloak_application_realm", "keycloak application realm", "Environment (development|staging|production)")
	flag.StringVar(&cfg.keycloakClientId, "keycloak_client_id", "keycloak ClientId", "Environment (development|staging|production)")

	flag.Parse()

	router := http.NewServeMux()

	authService := &authService.AuthService{}
	authRouter.RegisterRoutes(router, authService)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting server on %s", cfg.port)
	if err := http.ListenAndServe(cfg.port, router); err != nil {
		log.Fatal(err)
	}
	//create test containers
	// _ = MasterTestContainer.NewMasterTestContainer()
}
