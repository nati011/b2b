package resource

import (
	"context"
	"os"
	"testing"
)

var container TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = NewTestContainer()
}

func Test_create_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	in := CreateRequest{
		Action: "test",
		Name:   "test",
	}

	got, err := container.ResourceService.Create(ctx, &in)
	if err != nil {
		t.Errorf("Failed to create resource err: %v", err)
	}
	if got == 0 {
		t.Errorf("Expected id != from %v", got)
	}
}

func Test_create_unhappyPath(t *testing.T) {
	t.Run("duplicateName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := CreateRequest{
			Action: "test",
			Name:   "test3",
		}
		_, err := container.ResourceService.Create(ctx, &in)
		if err != nil {
			t.Errorf("Failed to create resource err: %v", err)
		}

		//create duplicate
		wantErr := ErrDuplicateName
		got, err := container.ResourceService.Create(ctx, &in)
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
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := CreateRequest{
			Action: "",
			Name:   "test2",
		}

		wantErr := ErrEmptyAction
		got, err := container.ResourceService.Create(ctx, &in)
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
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := CreateRequest{
			Action: "test8",
			Name:   "",
		}

		wantErr := ErrEmptyName
		got, err := container.ResourceService.Create(ctx, &in)
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
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	in := CreateRequest{
		Action: "test",
		Name:   "test",
	}
	id, _ := container.ResourceService.Create(ctx, &in)

	//delete resource
	err := container.ResourceService.Delete(ctx, id)
	if err != nil {
		t.Errorf("Failed to delete resource err: %v", err)
	}

	//verify deletion
	resp, _ := container.ResourceService.Get(ctx, id)
	if resp.Id == id {
		t.Error("Failed to delete resource")
	}
}

func Test_delete_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		err := container.ResourceService.Delete(ctx, 1010)
		wantErr := ErrIdNotFound
		if err != wantErr {
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
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	id, _ := container.ResourceService.Create(ctx, &CreateRequest{
		Action: "test",
		Name:   "test",
	})

	//update
	in := UpdateRequest{
		Id:     id,
		Action: "anotherTest",
		Name:   "anotherTest",
	}
	resp, err := container.ResourceService.Update(ctx, &in)
	if err != nil {
		t.Errorf("Failed to update err %v", err)
	}
	if id != resp {
		t.Errorf("Failed to update resource")
	}
}

func Test_update_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		id, _ := service.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "test",
		})

		in := UpdateRequest{
			Id:     id,
			Action: "test",
			Name:   "demo test",
		}
		wantErr := ErrIdNotFound
		_, err := container.ResourceService.Update(ctx, &in)
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("duplicateName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()

		//create resource with taken name
		container.ResourceService.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "taken",
		})

		// create resource
		id, _ := container.ResourceService.Create(ctx, &CreateRequest{
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
		resp, err := container.ResourceService.Update(ctx, &in)
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
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		id, _ := container.ResourceService.Create(ctx, &CreateRequest{
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
		resp, err := container.ResourceService.Update(ctx, &in)
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
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		id, _ := container.ResourceService.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "test1",
		})
		// update
		in := UpdateRequest{
			Id:     id,
			Action: "",
			Name:   "test2",
		}
		wantErr := ErrEmptyAction
		resp, err := container.ResourceService.Update(ctx, &in)
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
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		id, _ := container.ResourceService.Create(ctx, &CreateRequest{
			Action: "test",
			Name:   "test",
		})
		got, err := container.ResourceService.Get(ctx, id)
		if err != nil {
			t.Errorf("Failed to get resource by Id err %v", err)
		}
		if got.Id != id {
			t.Errorf("Expected id: %v Got: %v", id, got.Id)
		}
	})

	t.Run("getByName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := &CreateRequest{
			Action: "test",
			Name:   "test",
		}
		_, err := container.ResourceService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		got, err := container.ResourceService.GetByName(ctx, "test")
		if err != nil {
			t.Errorf("Failed to get resource by Id err %v", err)
		}
		if got.Name != in.Name {
			t.Errorf("Expected name: %v Got: %v", in.Name, got.Name)
		}
	})
}

func Test_getResource_unhappyPath(t *testing.T) {
	t.Run("getById", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		// Get by Id
		_, err := container.ResourceService.Get(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err %v Got: %v", wantErr, err)
		}

	})

	t.Run("getByName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		// Get by Id
		_, err := container.ResourceService.GetByName(ctx, "name")
		wantErr := ErrNameNotFound
		if err != wantErr {
			t.Errorf("Expected err %v Got: %v", wantErr, err)
		}

	})

}

func Test_getAllResources_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	id, _ := container.ResourceService.Create(ctx, &CreateRequest{
		Action: "test",
		Name:   "test",
	})

	// Get container.ResourceService
	got, err := container.ResourceService.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get resource by Id err %v", err)
	}
	if len(got.List) == 0 || got.List[0].Id != id {
		t.Errorf("Failed to get resource by id")
	}
}

func Test_getAllResources_unhappyPath(t *testing.T) {
	t.Run("getAll", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.ResourceService.GetAll(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err %v Got: %v", wantErr, err)
		}

	})
}
