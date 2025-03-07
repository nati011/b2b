package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	config "b2b.nati011.github.com/config"
	authRouter "b2b.nati011.github.com/internal/adapter/primary/routes/auth"
	distributorRouter "b2b.nati011.github.com/internal/adapter/primary/routes/distributor"
	healthRouter "b2b.nati011.github.com/internal/adapter/primary/routes/health"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := *config.LoadConfig()

	DB, err := pgxpool.New(context.Background(), cfg.DB_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer DB.Close()
	router := http.NewServeMux()
	authRouter.RegisterRoutes(router)
	healthRouter.RegisterRoutes(router)
	distributorRouter.RegisterRoutes(router, DB)
	log.Printf("Starting server on %s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
