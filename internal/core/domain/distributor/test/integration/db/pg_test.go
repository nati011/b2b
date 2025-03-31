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

	"b2b.nati011.github.com/internal/core/domain/distributor"
	test_container "b2b.nati011.github.com/internal/core/domain/distributor/test"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testContainer test_container.TestContainer
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
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// ddl
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db.sql")
	if err != nil {
		log.Fatalf("Error running ddl migration: %v", err)
	}

	// functions
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db_functions.sql")
	if err != nil {
		log.Fatalf("Error running stored func migration: %v", err)
	}
	testContainer = test_container.NewDBIntegrationTestContainer(
		db,
	)
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

func Test_Timeout(t *testing.T) {
}

func Test_Read(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := distributor.CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := testContainer.RetailerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got :%v", id, resp.Id)
		}
	})

	t.Run("get_all", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		// setup
		in := distributor.CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.RetailerService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(resp.List))
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})

	t.Run("get_by_name", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		// setup
		in := distributor.CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.RetailerService.GetByParam(ctx, &distributor.GetByParamRequest{
			Name: in.FirstName + in.LastName,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(resp.List))
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})

	t.Run("get_by_tin", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		// setup
		in := distributor.CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.RetailerService.GetByParam(ctx, &distributor.GetByParamRequest{
			Tin: in.Tin,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", len(resp.List), wantLen)
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})

	t.Run("get_all_user_agents", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		in := distributor.CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		_, err = testContainer.RetailerService.GetAllUsers(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get user agents %v", err)
		}

	})
}

func Test_Write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := distributor.CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",

			Email: "test@gmail.com",
		}
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := testContainer.RetailerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got :%v", id, resp.Id)
		}
	})

	t.Run("create_user_agent", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		in := distributor.CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",

			Email: "test@gmail.com",
		}
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check if user agent has been created
		_, err = testContainer.RetailerService.GetAllUsers(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get user agents %v", err)
		}
	})

	t.Run("updateName", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		in := distributor.CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		update_in := distributor.UpdateRequest{
			Id:   id,
			Name: "test",
		}
		_, err = testContainer.RetailerService.Update(ctx, &update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		//check
		resp, err := testContainer.RetailerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Name != update_in.Name {
			t.Errorf("Expected name: %v Got:%v", update_in.Name, resp.Name)
		}
	})

	t.Run("updateTin", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
		ctx := context.Background()
		in := distributor.CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		update_in := distributor.UpdateRequest{
			Id:  id,
			Tin: "1234567891",
		}
		_, err = testContainer.RetailerService.Update(ctx, &update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		// check
		resp, err := testContainer.RetailerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		if resp.Tin != update_in.Tin {
			t.Errorf("Expected Tin: %v Got:%v", update_in.Tin, resp.Tin)
		}
	})
}
