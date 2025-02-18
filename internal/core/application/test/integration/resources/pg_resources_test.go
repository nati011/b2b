package integration

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/resource/db"
	resource "b2b.nati011.github.com/internal/core/application/service/resource"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var service resource.Provider
var pgContainer *postgres.PostgresContainer

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

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	service = resource.NewResource(
		db_adapter.NewPostgres(
			db,
		),
	)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// db.
}

func RunContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	return postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
	)
}

func Test_Timeout(t *testing.T) {

}

func Test_create_happyPath(t *testing.T) {
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

func Test_create_unhappyPath(t *testing.T) {
	t.Run("duplicateName", func(t *testing.T) {
		//setup
		ctx := context.Background()
		in := resource.CreateRequest{
			Action: "test",
			Name:   "test",
		}
		_, err := service.Create(ctx, &in)
		if err != nil {
			t.Errorf("Failed to create resource err: %v", err)
		}

		//create duplicate
		wantErr := resource.ErrDuplicateName
		got, err := service.Create(ctx, &in)
		if err != wantErr {
			switch err {
			case nil:
				t.Errorf("Expected err: %v Got err: %v", wantErr, err)
			default:
				t.Errorf("Failed to create resource err: %v", err)
			}
		}
		if got != 0 {
			t.Errorf("Expected: %v, Got: %v", 0, got)
		}
	})

	t.Run("emptyAction", func(t *testing.T) {
		ctx := context.Background()
		in := resource.CreateRequest{
			Action: "",
			Name:   "test",
		}

		wantErr := resource.ErrEmptyAction
		got, err := service.Create(ctx, &in)
		if err != wantErr {
			switch err {
			case nil:
				t.Errorf("Expected err: %v Got err: %v", wantErr, err)
			default:
				t.Errorf("Failed to create resource err: %v", err)
			}
		}
		if got != 0 {
			t.Errorf("Expected: %v, Got: %v", 0, got)
		}
	})

	t.Run("emptyName", func(t *testing.T) {
		ctx := context.Background()
		in := resource.CreateRequest{
			Action: "test",
			Name:   "",
		}

		wantErr := resource.ErrEmptyName
		got, err := service.Create(ctx, &in)
		if err != wantErr {
			switch err {
			case nil:
				t.Errorf("Expected err: %v Got err: %v", wantErr, err)
			default:
				t.Errorf("Failed to create resource err: %v", err)
			}
		}
		if got != 0 {
			t.Errorf("Expected: %v, Got: %v", 0, got)
		}
	})
}

func Test_delete_happyPath(t *testing.T) {
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

func Test_delete_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		err := service.Delete(ctx, 1010)
		wantErr := resource.ErrIdNotFound
		if err != nil {
			switch err {
			case resource.ErrIdNotFound:
				t.Errorf("Expected err: %q Got err: %q", wantErr, err)
			default:
				t.Error("Failed to delete resource")
			}
		}
	})
}

func Test_update_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	id, _ := service.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "test",
	})

	//update
	in := resource.UpdateRequest{
		Id:     id,
		Action: "",
		Name:   "",
	}
	resp, err := service.Update(ctx, &in)
	if err != nil {
		t.Errorf("Failed to update err %v", err)
	}
	if id != resp {
		t.Errorf("Failed to update resource")
	}
}

func Test_update_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		in := resource.UpdateRequest{
			Id:     1,
			Action: "test",
			Name:   "test",
		}
		wantErr := resource.ErrIdNotFound
		resp, err := service.Update(ctx, &in)
		if err != nil {
			switch err {
			case wantErr:
				t.Errorf("Expected %v Got %v", wantErr, err)
			default:
				t.Errorf("Failed to update err %v", err)
			}
		}
		if resp != 0 {
			t.Errorf("Failed to update resource")
		}
	})

	t.Run("duplicateName", func(t *testing.T) {
		// setup
		ctx := context.Background()

		//create resource with taken name
		service.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "taken",
		})

		// create resource
		id, _ := service.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "test",
		})
		// update
		in := resource.UpdateRequest{
			Id:     id,
			Action: "test",
			Name:   "takenName",
		}
		wantErr := resource.ErrDuplicateName
		resp, err := service.Update(ctx, &in)
		if err != nil {
			switch err {
			case wantErr:
				t.Errorf("Expected %v Got %v", wantErr, err)
			default:
				t.Errorf("Failed to update err %v", err)
			}
		}
		if id != resp {
			t.Errorf("Failed to update resource")
		}

	})

	t.Run("emptyName", func(t *testing.T) {
		// setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "test",
		})
		// update
		in := resource.UpdateRequest{
			Id:     id,
			Action: "tets",
			Name:   "",
		}
		wantErr := resource.ErrDuplicateName
		resp, err := service.Update(ctx, &in)
		if err != nil {
			switch err {
			case wantErr:
				t.Errorf("Expected %v Got %v", wantErr, err)
			default:
				t.Errorf("Failed to update err %v", err)
			}
		}
		if id != resp {
			t.Errorf("Failed to update resource")
		}
	})

	t.Run("emptyAction", func(t *testing.T) {
		// setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "test1",
		})
		// update
		in := resource.UpdateRequest{
			Id:     id,
			Action: "",
			Name:   "test2",
		}
		wantErr := resource.ErrEmptyAction
		resp, err := service.Update(ctx, &in)
		if err != nil {
			switch err {
			case wantErr:
				t.Errorf("Expected %v Got %v", wantErr, err)
			default:
				t.Errorf("Failed to update err %v", err)
			}
		}
		if id != resp {
			t.Errorf("Failed to update resource")
		}
	})
}

func Test_getResource_happyPath(t *testing.T) {
	t.Run("getById", func(t *testing.T) {
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

func Test_getResource_unhappyPath(t *testing.T) {
}

func Test_getAllResources_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	id, _ := service.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "test",
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

func Test_getAllResources_unhappyPath(t *testing.T) {
}
