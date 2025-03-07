package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	config "b2b.nati011.github.com/config"
	authRouter "b2b.nati011.github.com/internal/adapter/primary/routes/auth"
	distributorRouter "b2b.nati011.github.com/internal/adapter/primary/routes/distributor"
	healthRouter "b2b.nati011.github.com/internal/adapter/primary/routes/health"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	cfg := *config.LoadConfig()

	var err error
	ctx := context.Background()

	db, err := sql.Open("pgx", cfg.DB_URL)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
	router := http.NewServeMux()
	authRouter.RegisterRoutes(router)
	healthRouter.RegisterRoutes(router)
	distributorRouter.RegisterRoutes(router, db)
	log.Printf("Starting server on %s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
