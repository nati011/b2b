package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	migrationsDir      = "db/migrations"
	publicSchemaTables = "SELECT quote_ident(tablename) FROM pg_tables WHERE schemaname = 'public';"
)

// findProjectRoot finds the project root by looking for go.mod file
func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	dir := wd
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("project root not found (go.mod not found)")
}

// SetupTestDB starts a disposable Postgres instance via Testcontainers, runs the
// project migrations, and returns a connection plus a cleanup function suitable
// for use in tests: db, cleanup := test.SetupTestDB(t)
func SetupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	container, err := RunContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	db, err := connectDB(ctx, container)
	if err != nil {
		_ = container.Terminate(context.Background())
		t.Fatalf("failed to connect to postgres container: %v", err)
	}

	connectionString, err := container.ConnectionString(ctx)
	if err != nil {
		cleanupContainerResources(db, container)
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Disable SSL for test database
	if !strings.Contains(connectionString, "sslmode=") {
		if strings.Contains(connectionString, "?") {
			connectionString += "&sslmode=disable"
		} else {
			connectionString += "?sslmode=disable"
		}
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		cleanupContainerResources(db, container)
		t.Fatalf("failed to find project root: %v", err)
	}

	absMigrationsPath := filepath.Join(projectRoot, migrationsDir)
	if _, err := os.Stat(absMigrationsPath); err != nil {
		cleanupContainerResources(db, container)
		t.Fatalf("migrations directory not found at %s: %v", absMigrationsPath, err)
	}

	if err := RunMigrations(connectionString, absMigrationsPath); err != nil {
		cleanupContainerResources(db, container)
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		if err := resetDatabase(db); err != nil {
			t.Logf("failed to reset database: %v", err)
		}
		if err := db.Close(); err != nil {
			t.Logf("failed to close database connection: %v", err)
		}
		if err := container.Terminate(context.Background()); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	}

	return db, cleanup
}

func connectDB(ctx context.Context, container *postgres.PostgresContainer) (*sql.DB, error) {
	connectionString, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("get connection string: %w", err)
	}

	// Disable SSL for test database
	if !strings.Contains(connectionString, "sslmode=") {
		if strings.Contains(connectionString, "?") {
			connectionString += "&sslmode=disable"
		} else {
			connectionString += "?sslmode=disable"
		}
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func resetDatabase(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var tables []string
	rows, err := tx.Query(publicSchemaTables)
	if err != nil {
		return fmt.Errorf("fetch public tables: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return fmt.Errorf("scan table name: %w", err)
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate table rows: %w", err)
	}

	if len(tables) > 0 {
		truncateQuery := "TRUNCATE TABLE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE;"
		if _, err := tx.Exec(truncateQuery); err != nil {
			return fmt.Errorf("truncate tables: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func cleanupContainerResources(db *sql.DB, container *postgres.PostgresContainer) {
	_ = db.Close()
	_ = container.Terminate(context.Background())
}

func RunContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	return postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),

		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
}

// GetResourceManifestPath returns the path to the resource manifest YAML file.
// This is a helper for tests that need to bootstrap resources.
func GetResourceManifestPath(t *testing.T) string {
	t.Helper()

	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("failed to find project root: %v", err)
	}

	return filepath.Join(projectRoot, "config", "resources.yaml")
}
