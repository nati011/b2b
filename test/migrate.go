package test

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies database migrations from the specified directory using golang-migrate.
// It uses the migrate library to handle version tracking and migration execution.
func RunMigrations(connectionString, migrationsPath string) error {
	// Convert absolute path to file:// URL format
	migrationsURL := "file://" + filepath.ToSlash(migrationsPath)

	// Create migrate instance with file source and postgres database
	m, err := migrate.New(migrationsURL, connectionString)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
