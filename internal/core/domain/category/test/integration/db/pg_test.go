package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	category "b2b.nati011.github.com/internal/core/domain/category"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var service category.Provider
var pgContainer *postgres.PostgresContainer
var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	var err error
	ctx := context.Background()

	pgContainer, err = RunContainer(ctx)
	if err != nil {
		panic(err)
	}

	connectionString, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	db, err = sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	service = category.NewCategory(
		db_adapter.NewPostgres(
			db,
		),
	)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// ddl
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db.sql")
	if err != nil {
		log.Fatalf("Error running migration: %v", err)
	}

	// functions
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db_functions.sql")
	if err != nil {
		log.Fatalf("Error running migration: %v", err)
	}

}

func teardown() {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("could not begin transaction: %v", err)
	}

	// Get all table names
	var tables []string
	rows, err := tx.Query("SELECT tablename FROM pg_tables WHERE schemaname = 'public';")
	if err != nil {
		tx.Rollback()
		log.Fatalf("could not fetch table names: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			tx.Rollback()
			log.Fatalf("could not scan table name: %v", err)
		}
		tables = append(tables, table)
	}

	// Prepare the TRUNCATE statement
	if len(tables) > 0 {
		truncateQuery := "TRUNCATE TABLE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE;"
		_, err = tx.Exec(truncateQuery)
		if err != nil {
			tx.Rollback()
			log.Fatalf("could not truncate tables: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("could not commit transaction: %v", err)
	}
}

func runMigration(db *sql.DB, filename string) error {
	// Read the SQL file
	sqlBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	// Execute the SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("could not execute SQL: %w", err)
	}

	return nil
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

func Test_Reader(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
		ctx := context.Background()
		in := &category.CreateRequest{
			Name: "test",
			Desc: "test",
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create category err: %v", err)
		}

		// get
		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get category err: %v", err)
		}

		if resp.Id != id {
			t.Errorf("Expected id: %v Got id: %v", resp.Id, id)
		}
	})

	t.Run("get_all", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
		ctx := context.Background()
		in := &category.CreateRequest{
			Name: "test",
			Desc: "test",
		}
		_, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create category err: %v", err)
		}

		// get
		resp, err := service.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to get category err: %v", err)
		}

		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp.List))
		}
	})
}

func Test_Writer(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &category.CreateRequest{
			Name: "test",
			Desc: "test",
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create category err: %v", err)
		}

		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get category err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got id: %v", id, resp.Id)
		}
	})

	t.Run("remove", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		in := &category.CreateRequest{
			Name: "test",
			Desc: "test",
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create category err: %v", err)
		}

		//remove
		err = service.Remove(ctx, id)
		if err != nil {
			t.Fatalf("Failed to remove category err: %v", err)
		}

		//check
		_, err = service.Get(ctx, id)
		wantErr := category.ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}
