package resource

import (
	"context"
	"os"
	"testing"

	db_mock "b2b.nati011.github.com/internal/adapter/secondary/resource/db"
)

var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewResource(
		db_mock.NewMock(),
	)
}

func Test_create_happyPath(t *testing.T) {
	ctx := context.Background()
	in := CreateRequest{
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
		in := CreateRequest{
			Action: "test",
			Name:   "test",
		}
		_, err := service.Create(ctx, &in)
		if err != nil {
			t.Errorf("Failed to create resource err: %v", err)
		}

		//create duplicate
		wantErr := ErrDuplicateName
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
		in := CreateRequest{
			Action: "",
			Name:   "test",
		}

		wantErr := ErrEmptyAction
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
		in := CreateRequest{
			Action: "test",
			Name:   "",
		}

		wantErr := ErrEmptyName
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
	in := CreateRequest{
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
	resp, _ := service.Get(ctx, &GetRequest{
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
		wantErr := ErrIdNotFound
		if err != nil {
			switch err {
			case ErrIdNotFound:
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
	id, _ := service.Create(ctx, &CreateRequest{
		Action: "test",
		Name:   "test",
	})

	//update
	in := UpdateRequest{
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
		in := UpdateRequest{
			Id:     1,
			Action: "test",
			Name:   "test",
		}
		wantErr := ErrIdNotFound
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
		service.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "taken",
		})

		// create resource
		id, _ := service.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "test",
		})
		// update
		in := UpdateRequest{
			Id:     id,
			Action: "test",
			Name:   "takenName",
		}
		wantErr := ErrDuplicateName
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
		id, _ := service.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "test",
		})
		// update
		in := UpdateRequest{
			Id:     id,
			Action: "tets",
			Name:   "",
		}
		wantErr := ErrDuplicateName
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
		id, _ := service.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "test",
		})
		// update
		in := UpdateRequest{
			Id:     id,
			Action: "",
			Name:   "test",
		}
		wantErr := ErrEmptyAction
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
		id, _ := service.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "test",
		})

		// Get by Id
		in := GetRequest{
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
		id, _ := service.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "test",
		})

		// Get by Id
		in := GetRequest{
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
	id, _ := service.Create(ctx, &CreateRequest{
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
