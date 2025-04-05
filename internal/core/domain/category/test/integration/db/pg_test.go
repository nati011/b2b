package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/category/db"
	category "b2b.nati011.github.com/internal/core/domain/category"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
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
	db = db_test_container.Setup()
	service = category.NewCategory(
		db_adapter.NewPostgres(
			db,
		),
	)

}

func teardown() {
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_Reader(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
		ctx := context.Background()
		in := &category.CreateRequest{
			Name: "test",
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
