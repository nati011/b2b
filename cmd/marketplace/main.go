package main

import (
	"context"
	"marketplace/internal/app"
	"marketplace/internal/config"
	"marketplace/internal/infra/banner"
	"marketplace/pkg/logger"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := loadConfig()
	initializeLogger(cfg)
	displayBanner(cfg)

	application := initializeApplication(cfg)
	bootstrapApplication(application)

	runApplication(application)
	waitForShutdown()
	shutdownApplication(application)
}

// loadConfig parses command-line arguments and loads the configuration.
func loadConfig() *config.Config {
	configPath := flag.String("config", "config/config.yaml", "Path to the configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Default()
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	return cfg
}

// initializeLogger sets up logging based on configuration.
func initializeLogger(cfg *config.Config) {
	// Set defaults for new config options
	sourceLocation := cfg.Logging.SourceLocation
	samplingRate := cfg.Logging.SamplingRate
	if samplingRate == 0 {
		samplingRate = 1.0 // Default to log everything
	}
	async := cfg.Logging.Async
	sanitize := cfg.Logging.Sanitize

	logger.InitWithConfig(cfg.Logging.Level, cfg.Logging.Format, sourceLocation, samplingRate, async, sanitize, &cfg.App)
	logger.WithAppContext(cfg)
}

// displayBanner prints the application banner.
func displayBanner(cfg *config.Config) {
	banner.Print(cfg.App.Version, cfg.App.Env)
	logger.Info("Application starting", "env", cfg.App.Env)
}

// initializeApplication wires up and initializes the application dependencies.
func initializeApplication(cfg *config.Config) *app.Application {
	application, err := app.InitializeApp(cfg)
	if err != nil {
		logger.Error("Failed to initialize app", "error", err)
		os.Exit(1)
	}
	return application
}

// bootstrapApplication seeds resources from the manifest (idempotent operation).
func bootstrapApplication(application *app.Application) {
	bootstrapCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Bootstrap(bootstrapCtx); err != nil {
		logger.Error("Failed to bootstrap resources", "error", err)
		os.Exit(1)
	}
}

// runApplication starts the HTTP server in a goroutine.
func runApplication(application *app.Application) {
	go func() {
		if err := application.Start(); err != nil {
			logger.Error("Failed to start application", "error", err)
			os.Exit(1)
		}
	}()
}

// waitForShutdown blocks until a shutdown signal is received.
func waitForShutdown() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	<-ctx.Done()
	stop()
}

// shutdownApplication gracefully shuts down the application.
func shutdownApplication(application *app.Application) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Shutting down application...")
	if err := application.Shutdown(shutdownCtx); err != nil {
		logger.Error("Shutdown error", "error", err)
	}
	logger.Info("Application shutdown successfully")

	// Shutdown async logger if enabled
	logger.Shutdown()
}
