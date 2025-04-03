package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	db_resource_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	db_role_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/user/db"
	"b2b.nati011.github.com/internal/core/application/auth"
	resource "b2b.nati011.github.com/internal/core/application/resource"
	role "b2b.nati011.github.com/internal/core/application/role"
	user "b2b.nati011.github.com/internal/core/application/user"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testContainer user.TestContainer
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

	testContainer = user.NewIntegrationTestContainer(db)

	testContainer.UserService = user.NewUser(
		db_adapter.NewPostgres(
			db,
		),
		role.NewRole(
			db_role_adapter.NewPostgres(
				db,
			), resource.NewResource(
				db_resource_adapter.NewPostgres(
					db,
				),
			),
		),
		auth.NewIntegrationAuthContainer(),
	)

	testContainer.RoleService = role.NewRole(
		db_role_adapter.NewPostgres(
			db,
		), resource.NewResource(
			db_resource_adapter.NewPostgres(
				db,
			),
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

func Test_Timeout(t *testing.T) {}

func Test_write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName:  "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			ExternalId: "123",
		}
		id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		_, err = testContainer.UserService.Get(ctx, id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}
	})
}

func Test_read(t *testing.T) {
	t.Run("get_all", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "natnael jemaneh asefa",
			Email:     "natnaeljemaneh001@gmail.com",
			Phone:     "+251949184879",
			Username:  "test",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//get
		resp, err := testContainer.UserService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to getAll err: %v", err)
		}
		expecetdLen := 1
		if len(resp.List) == 0 {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(resp.List))
		}
	})

	t.Run("get_by_param", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName:  "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		inParam := &user.GetByParam{
			Username: "test",
		}
		response, err := testContainer.UserService.GetByParam(ctx, inParam)
		if err != nil {
			t.Fatalf("Failed to get user by param err: %v", err)
		}
		expecetdLen := 1
		if len(response.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(response.List))
		}
	})

	t.Run("update", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02", "2024-09-20")
		in := user.CreateRequest{
			FirstName:  "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		in_updateParsedTime, _ := time.Parse("2006-01-02", "2024-09-19")
		in_update := &user.UpdateRequest{
			Id:         user_id,
			FirstName:  "test",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        in_updateParsedTime,
			ExternalId: "123",
		}
		got, err := testContainer.UserService.Update(ctx, in_update)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}
		if got.FirstName != in_update.FirstName {
			t.Errorf("Expected : %v Got: %v", in_update.DOB, got.DOB)
		}
		if got.Email != in_update.Email {
			t.Errorf("Expected : %v Got: %v", in_update.Email, got.Email)
		}
		if got.Phone != in_update.Phone {
			t.Errorf("Expected : %v Got: %v", in_update.Phone, got.Phone)
		}
		if got.Username != in_update.Username {
			t.Errorf("Expected : %v Got: %v", in_update.Username, got.Username)
		}
		if got.DOB != in_update.DOB {
			t.Errorf("Expected : %v Got: %v", in_update.DOB, got.DOB)
		}
		if got.ExternalId != in_update.ExternalId {
			t.Errorf("Expected : %v Got: %v", in_update.ExternalId, got.ExternalId)
		}
	})

}
