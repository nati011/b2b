package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"b2b.nati011.github.com/config"

	util "b2b.nati011.github.com/internal/adapter/primary/rest/handler/util"
	application_core "b2b.nati011.github.com/internal/core/application"
	domain_core "b2b.nati011.github.com/internal/core/domain"
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
	flag.StringVar(&cfg.KeycloakClientSecret, "keycloak_client_secret", "", "Environment (development|staging|production)")

	//email
	flag.StringVar(&cfg.Email, "email", "", "Environment (development|staging|production)")
	flag.StringVar(&cfg.SMTP, "smtp", "", "Environment (development|staging|production)")

	//db
	flag.StringVar(&cfg.FileLocation, "migration_file_dir", "", "Environment (development|staging|production)")
	flag.StringVar(&cfg.CoreDBConnectionString, "db", "", "Environment (development|staging|production)")
	flag.Parse()
	validateFlags(cfg)

	db_pool := InitDB(cfg.CoreDBConnectionString, cfg.FileLocation)
	InitEmail(cfg.Email, cfg.SMTP)
	// InitAuth(cfg.Port, cfg.Env, cfg.KeycloakInstanceURL, cfg.KeycloakUsername, cfg.KeycloakPassword, cfg.KeycloakRealm, cfg.KeycloakApplicationRealm, cfg.KeycloakClientId, cfg.KeycloakClientSecret)
	InitSMS(cfg.Email, cfg.SMTP)

	application_constainer := application_core.NewContainer(db_pool, cfg.KeycloakInstanceURL, cfg.KeycloakUsername, cfg.KeycloakPassword, cfg.KeycloakRealm, cfg.KeycloakApplicationRealm, cfg.KeycloakClientId, cfg.Email, cfg.SMTP, cfg.KeycloakClientSecret)
	domain_container := domain_core.NewContainer(*application_constainer, db_pool)

	mux := http.NewServeMux()
	InitREST(mux, db_pool, application_constainer, domain_container)

	loggingingMiddleware := util.NewLoggingMiddleware()
	handler := loggingingMiddleware.Log(mux)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("starting %s server on %s", cfg.Env, srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
