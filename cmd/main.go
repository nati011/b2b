package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"b2b.nati011.github.com/config"
	"github.com/go-co-op/gocron/v2"

	application_core "b2b.nati011.github.com/internal/core/application"
	"b2b.nati011.github.com/internal/core/application/middleware"
	domain_core "b2b.nati011.github.com/internal/core/domain"
)

func main() {
	var cfg config.Config

	flag.IntVar(&cfg.Port, "port", 4000, "api server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.KeycloakInstanceURL, "keycloak_base_url", "", "authServer baseurl")
	flag.StringVar(&cfg.KeycloakUsername, "keycloak_user_name", "", "authServer username")
	flag.StringVar(&cfg.KeycloakPassword, "keycloak_password", "", "authServer password")
	flag.StringVar(&cfg.KeycloakRealm, "keycloak_realm", "", "authServer realm")
	flag.StringVar(&cfg.KeycloakApplicationRealm, "keycloak_application_realm", "", "authServer application realm")
	flag.StringVar(&cfg.KeycloakClientId, "keycloak_client_id", "", "authServer clientId")
	flag.StringVar(&cfg.KeycloakClientSecret, "keycloak_client_secret", "", "authServer client secret")
	flag.StringVar(&cfg.Email, "email", "", "email server address")
	flag.StringVar(&cfg.SMTP, "smtp", "", "email server address smtp")
	flag.StringVar(&cfg.EmailPassword, "email_password", "", "email server password")
	flag.StringVar(&cfg.FileLocation, "migration_file_dir", "", "migration file location dir")
	flag.StringVar(&cfg.CoreDBConnectionString, "db", "", "coreDB connection string (eg. postgres://postgres:1234@localhost:5432/b2b_1136)")
	flag.StringVar(&cfg.MinMobileClientCompatibleVersion, "min_compatible_client_version", "1.0.0", "Min Mobile client version")
	flag.StringVar(&cfg.BaseUrl, "base_url", "http://localhost:8080", "backend base url")
	flag.StringVar(&cfg.FrontendUrl, "frontend_base_url", "http://localhost:8080", "frontend base url")
	flag.StringVar(&cfg.JWTSecret, "jwt_secret", "", "jwt secret keys")
	flag.StringVar(&cfg.DefaultSuperAdminUserEmail, "default_superadmin_email", "", "")
	flag.Parse()
	validateFlags(cfg)

	db_pool := InitDB(&cfg)
	InitAuth(&cfg)
	// InitEmail(cfg.Email, cfg.SMTP, cfg.EmailPassword)
	// InitSMS(cfg.Email, cfg.SMTP)

	pagination := config.DefaultPaginationBuilder().Build()
	paginationMiddleware := middleware.NewPaginationMiddleware(pagination)
	application_container := application_core.NewContainer(
		db_pool,
		&cfg,
		pagination)

	domain_container := domain_core.NewContainer(*application_container, cfg.BaseUrl, cfg.FrontendUrl, db_pool)

	mux := http.NewServeMux()
	InitREST(mux, db_pool, application_container, domain_container)

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	InitCron(s, application_container, domain_container)
	s.Start()

	InitEventLister(application_container, domain_container)

	InitDefaultConfig(cfg, application_container)

	loggingingMiddleware := middleware.NewLoggingMiddleware()
	handler := paginationMiddleware.Paginate(loggingingMiddleware.Log(mux))
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("Ahoy! server running %s on %s ...", cfg.Env, srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
