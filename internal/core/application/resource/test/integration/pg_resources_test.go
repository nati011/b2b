package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"

	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	resource "b2b.nati011.github.com/internal/core/application/resource"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var service resource.Provider
var pgContainer *postgres.PostgresContainer
var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	service = resource.NewResource(
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

func Test_create(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := resource.CreateRequest{
		Action: "test",
		Name:   "test",
	}

	got, err := service.Create(ctx, &in)
	if err != nil {
		t.Errorf("Failed to create resource err: %v", err)
	}
	if got == 0 {
		t.Errorf("Expected id != from %v", got)
	}
}

func Test_delete(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	in := resource.CreateRequest{
		Action: "test",
		Name:   "test",
	}
	id, _ := service.Create(ctx, &in)

	//delete resource
	err := service.Delete(ctx, id)
	if err != nil {
		t.Errorf("Failed to delete resource err: %v", err)
	}

	//verify deletion
	resp, _ := service.Get(ctx, &resource.GetRequest{
		Id:   id,
		Name: "",
	})
	if resp.Id == id {
		t.Error("Failed to delete resource")
	}
}

func Test_update(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	id, _ := service.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "test",
	})

	//update
	in := resource.UpdateRequest{
		Id:     id,
		Action: "test1",
		Name:   "test1",
	}
	resp, err := service.Update(ctx, &in)
	if err != nil {
		t.Errorf("Failed to update err %v", err)
	}
	if id != resp {
		t.Errorf("Failed to update resource")
	}
}

func Test_get(t *testing.T) {
	t.Run("getById", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "test",
		})

		// Get by Id
		in := resource.GetRequest{
			Id: id,
		}
		got, err := service.Get(ctx, &in)
		if err != nil {
			t.Errorf("Failed to get resource by Id err %v", err)
		}
		if got.Id != id {
			t.Errorf("Failed to get resource by id")
		}
	})

	t.Run("getByName", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "test",
		})

		// Get by Id
		in := resource.GetRequest{
			Name: "test",
		}
		got, err := service.Get(ctx, &in)
		if err != nil {
			t.Errorf("Failed to get resource by Id err %v", err)
		}
		if got.Id != id {
			t.Errorf("Failed to get resource by id")
		}
	})
}

func Test_getAll(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	_, err := service.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "test",
	})

	if err != nil {
		t.Errorf("Failed to create resource err %v", err)
	}

	_, err = service.Create(ctx, &resource.CreateRequest{
		Action: "test1",
		Name:   "test2",
	})

	if err != nil {
		t.Errorf("Failed to create resource err %v", err)
	}

	_, err = service.Create(ctx, &resource.CreateRequest{
		Action: "test4",
		Name:   "test4",
	})

	if err != nil {
		t.Errorf("Failed to create resource err %v", err)
	}

	pagination := &resource.Pagination{
		Limit:  2,
		Offset: 0,
	}

	// Get All
	got, err := service.GetAll(ctx, pagination)
	if err != nil {
		t.Errorf("Failed to get resource err %v", err)
	}
	if len(got.List) == pagination.Limit && len(got.List) > pagination.Limit {
		t.Errorf("Want len %v Got len %v", pagination.Limit, len(got.List))
	}
}
