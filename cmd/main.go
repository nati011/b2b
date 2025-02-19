package main

import (
	"log"
	"net/http"

	config "b2b.nati011.github.com/config"
	authRouter "b2b.nati011.github.com/internal/adapter/primary/routes/auth"
)

func main() {
	cfg := *config.LoadConfig()
	router := http.NewServeMux()
	authRouter.RegisterRoutes(router)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting server on %s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.Fatal(err)
	}
	//create test containers
	// _ = MasterTestContainer.NewMasterTestContainer()
}
