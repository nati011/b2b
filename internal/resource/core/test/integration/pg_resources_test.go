package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/config"
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	resource "b2b.nati011.github.com/internal/core/application/resource"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var service resource.Provider
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
			config.DefaultPaginationBuilder().Build(),
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
		Action:   "test",
		Name:     "test121",
		Resource: "test",
		Scope:    "d",
	}

	_, err := service.Create(ctx, &in)
	if err != nil {
		t.Errorf("Failed to create resource err: %v", err)
	}
}

func Test_delete(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	in := resource.CreateRequest{
		Action:   "test",
		Name:     "test",
		Resource: "test",
	}
	id, err := service.Create(ctx, &in)
	if err != nil {
		t.Errorf("Failed to create resource err: %v", err)
	}

	//delete resource
	err = service.Delete(ctx, id)
	if err != nil {
		t.Errorf("Failed to delete resource err: %v", err)
	}

	//verify deletion
	resp, _ := service.Get(ctx, id)
	if resp.Id == id {
		t.Error("Failed to delete resource")
	}
}

func Test_update(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	id, _ := service.Create(ctx, &resource.CreateRequest{
		Action:   "test",
		Name:     "test",
		Resource: "test",
	})

	//update
	in := resource.UpdateRequest{
		Id:       id,
		Action:   "test1",
		Name:     "test1",
		Resource: "test2",
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
			Action:   "test",
			Name:     "test",
			Resource: "test",
		})

		// Get by Id
		got, err := service.Get(ctx, id)
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
			Action:   "test",
			Name:     "test",
			Resource: "test",
		})

		got, err := service.GetByName(ctx, "test")
		if err != nil {
			t.Errorf("Failed to get resource by Id err %v", err)
		}
		if got.Id != id {
			t.Errorf("Failed to get resource by id")
		}
	})

	t.Run("getByResource", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &resource.CreateRequest{
			Action:   "test",
			Name:     "test",
			Resource: "test",
		})

		got, err := service.GetByResource(ctx, "test")
		if err != nil {
			t.Errorf("Failed to get resource by resource err %v", err)
		}
		if got.Id != id {
			t.Errorf("Failed to get resource by resource")
		}
	})
}

func Test_getAll(t *testing.T) {
	t.Cleanup(teardown)
	//setup
	ctx := context.Background()
	id, _ := service.Create(ctx, &resource.CreateRequest{
		Action:   "test",
		Name:     "test",
		Resource: "test",
	})

	// Get All
	got, err := service.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get resource by Id err %v", err)
	}
	if len(got.List) == 0 || got.List[0].Id != id {
		t.Errorf("Failed to get resource by id")
	}
}
