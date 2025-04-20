package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/user"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testContainer user.TestContainer
var service user.Provider
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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// ddl
	err = runMigration(db, "/home/delilah/Documents/work/non-kifya/b2b/migration/core_db.sql")
	if err != nil {
		log.Fatalf("Error running ddl migration: %v", err)
	}

	// functions
	err = runMigration(db, "/home/delilah/Documents/work/non-kifya/b2b/migration/core_db_functions.sql")
	if err != nil {
		log.Fatalf("Error running stored func migration: %v", err)
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

func Test_assignRole_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName: "Natnael",
		LastName:  "Jemaneh",
		Email:     "natnaeljemaneh001@gmail.com",
		Phone:     "+251949184879",
		Username:  "test",
		DOB:       parsedTime,

		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_assignRole_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("roleDoesNotExist", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "Natnael",
			LastName:  "Jemaneh",
			Email:     "natnaeljemaneh007@gmail.com",
			Phone:     "+251949184879",
			Username:  "test07",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//assign
		err = testContainer.UserService.AssignRole(ctx, user_id, int(rand.Int31()))
		expectedErr := user.ErrRoleDoesNotExist
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

	t.Run("alreadyAssigned", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "Natnael",
			LastName:  "Jemaneh",
			Email:     "natnaeljemaneh0009@gmail.com",
			Phone:     "+251949184879",
			Username:  "testAssign",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//create role
		role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//assign
		err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}

		//reassign
		err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
		expectedErr := user.ErrRoleAlreadyAssigned
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

}

func Test_removeAssignedRole_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName: "Natnael",
		LastName:  "Jemaneh",
		Email:     "natnaeljemaneh0001@gmail.com",
		Phone:     "+251949184879",
		Username:  "demo",
		DOB:       parsedTime,

		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}

	//remove
	err = testContainer.UserService.RemoveAssignedRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_removeAssignedRole_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("roleNotAssigned", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "Natnael",
			LastName:  "Jemaneh",
			Email:     "natnaeljemaneh0011@gmail.com",
			Phone:     "+251949184879",
			Username:  "testRole",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//create role
		role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//remove
		expectedErr := user.ErrRoleNotAssigned
		err = testContainer.UserService.RemoveAssignedRole(ctx, user_id, role_id)
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

	t.Run("roleNotFound", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "Natnael",
			LastName:  "Jemaneh",
			Email:     "natnaeljemaneh5001@gmail.com",
			Phone:     "+251949184879",
			Username:  "demoAssigned",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//remove
		expectedErr := user.ErrRoleDoesNotExist
		err = testContainer.UserService.RemoveAssignedRole(ctx, user_id, int(rand.Int31()))
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

}

func Test_getAllRole_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName: "Natnael",
		LastName:  "Jemaneh",
		Email:     "natnaeljemaneh101@gmail.com",
		Phone:     "+251949184879",
		Username:  "testRoleHappyPath",
		DOB:       parsedTime,

		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}

	//get all assignment
	resp, err := testContainer.UserService.GetAllAssignedRoles(ctx, user_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
	expecetdLen := 1
	if len(resp.List) != expecetdLen {
		t.Errorf("Expected len: %v Got len: %v", expecetdLen, err)
	}
}

func Test_getAllRole_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("emptyRoleList", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "Natnael",
			LastName:  "Jemaneh",
			Email:     "natnaeljemaneh0901@gmail.com",
			Phone:     "+251949184879",
			Username:  "testAllUnhappy",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//create role
		_, err = testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//get all assignment
		wantErr := user.ErrNoRoleAssigned
		resp, err := testContainer.UserService.GetAllAssignedRoles(ctx, user_id)
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
		expecetdLen := 0
		if len(resp.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, err)
		}
	})
}

func Test_has_role_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("has", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "Natnael",
			LastName:  "Jemaneh",
			Email:     "natnaeljemaneh000@gmail.com",
			Phone:     "+251949184879",
			Username:  "test001",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//create role
		role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//assign
		err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}

		//check access to resource
		hasAccess, err := testContainer.UserService.HasRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}
		expectedStatus := true
		if !hasAccess {
			t.Errorf("Expected has access status: %v Got err: %v", expectedStatus, hasAccess)
		}
	})

	t.Run("has-not", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "Natnael",
			LastName:  "Jemaneh",
			Email:     "natnaeljemaneh002@gmail.com",
			Phone:     "+251949184879",
			Username:  "test2",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//create role
		role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//check access to role
		got, err := testContainer.UserService.HasRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}
		want := false
		if got != want {
			t.Errorf("Expected has access status: %v Got err: %v", want, got)
		}
	})
}

func Test_has_role_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	t.Run("roleDoesNotExist", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := user.CreateRequest{
			FirstName: "Natnael",
			LastName:  "Jemaneh",
			Email:     "natnaeljemaneh9001@gmail.com",
			Phone:     "+251949184879",
			Username:  "testRoleUnhappy",
			DOB:       parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check access to resource
		expectedErr := user.ErrRoleDoesNotExist
		got, err := testContainer.UserService.HasRole(ctx, user_id, rand.Int())
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
		want := false
		if got != want {
			t.Errorf("Expected has access status: %v Got: %v", want, got)
		}
	})
}
