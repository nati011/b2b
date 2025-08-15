package role

import (
	"context"
	"math/rand"
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
		Name: "test",
		Desc: "test",
	}

	_, err := container.RoleService.Create(ctx, &in)
	if err != nil {
		t.Errorf("Failed to create role err: %v", err)
	}

	getResp, _ := container.RoleService.GetAll(ctx)
	if len(getResp.List) == 0 {
		t.Errorf("No resources were created for role")
	}
}

func Test_create_unhappyPath(t *testing.T) {
	t.Run("duplicateName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := CreateRequest{
			Desc: "test",
			Name: "test",
		}
		_, err := container.RoleService.Create(ctx, &in)
		if err != nil {
			t.Errorf("Failed to create role err: %v", err)
		}

		//create duplicate
		wantErr := ErrDuplicateName
		got, err := container.RoleService.Create(ctx, &in)
		if err != wantErr {
			switch err {
			case nil:
				t.Errorf("Expected err: %v Got err: %v", wantErr, err)
			default:
				t.Errorf("Failed to create role err: %v", err)
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
			Desc: "test",
			Name: "",
		}

		wantErr := ErrEmptyName
		got, err := container.RoleService.Create(ctx, &in)
		if err != wantErr {
			switch err {
			case nil:
				t.Errorf("Expected err: %v Got err: %v", wantErr, err)
			default:
				t.Errorf("Failed to create role err: %v", err)
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
		Desc: "test",
		Name: "test",
	}
	id, _ := container.RoleService.Create(ctx, &in)

	//delete resource
	err := container.RoleService.Delete(ctx, id)
	if err != nil {
		t.Errorf("Failed to delete role err: %v", err)
	}

	//verify deletion
	resp, _ := container.RoleService.Get(ctx, &GetRequest{
		Id:   id,
		Name: "",
	})
	if resp.Id == id {
		t.Error("Failed to delete role")
	}
}

func Test_delete_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		err := container.RoleService.Delete(ctx, 1010)
		wantErr := ErrIdNotFound
		if err != wantErr {
			switch err {
			default:
				t.Errorf("Expected err: %q Got err: %q", wantErr, err)
			}
		}
	})
}

func Test_update_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	ctx := context.Background()
	id, _ := container.RoleService.Create(ctx, &CreateRequest{
		Desc: "test",
		Name: "test",
	})

	//update
	in := UpdateRequest{
		Id:   id,
		Desc: "test",
		Name: "tests",
	}
	resp, err := container.RoleService.Update(ctx, &in)
	if err != nil {
		t.Errorf("Failed to update err %v", err)
	}
	if id != resp {
		t.Errorf("Failed to update role")
	}
}

func Test_update_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		in := UpdateRequest{
			Id:   1,
			Desc: "test",
			Name: "test",
		}
		wantErr := ErrIdNotFound
		resp, err := container.RoleService.Update(ctx, &in)
		if err != wantErr {
			switch err {
			default:
				t.Errorf("Expected %v Got %v", wantErr, err)
			}
		}
		if resp != 0 {
			t.Errorf("Failed to update role")
		}
	})

	t.Run("duplicateName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()

		//create resource with taken name
		container.RoleService.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "taken",
		})

		// create resource
		id, _ := container.RoleService.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "test",
		})
		// update
		in := UpdateRequest{
			Id:   id,
			Desc: "test",
			Name: "taken",
		}
		wantErr := ErrDuplicateName
		got, err := container.RoleService.Update(ctx, &in)
		if err != wantErr {
			switch err {
			default:
				t.Errorf("Expected %v Got %v", wantErr, err)
			}
		}
		if got != 0 {
			t.Errorf("Failed to update role")
		}

	})

	t.Run("emptyName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		id, _ := container.RoleService.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "test",
		})
		// update
		in := UpdateRequest{
			Id:   id,
			Desc: "tets",
			Name: "",
		}
		wantErr := ErrDuplicateName
		resp, err := container.RoleService.Update(ctx, &in)
		if err != nil {
			switch err {
			case wantErr:
				t.Errorf("Expected %v Got %v", wantErr, err)
			default:
				t.Errorf("Failed to update err %v", err)
			}
		}
		if id != resp {
			t.Errorf("Failed to update role")
		}
	})

}

func Test_getRole_happyPath(t *testing.T) {
	t.Run("getById", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		id, _ := container.RoleService.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// Get by Id
		in := GetRequest{
			Id: id,
		}
		got, err := container.RoleService.Get(ctx, &in)
		if err != nil {
			t.Errorf("Failed to get role by Id err %v", err)
		}
		if got.Id != id {
			t.Errorf("Failed to get role by id")
		}
	})

	t.Run("getByName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		id, _ := container.RoleService.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// Get by Name
		in := GetRequest{
			Name: "test",
		}
		got, err := container.RoleService.Get(ctx, &in)
		if err != nil {
			switch err {
			case ErrIdNotFound:
			default:
				t.Errorf("Failed to get role by Id err %v", err)
			}
		}
		if got.Id != id {
			t.Errorf("Failed to get role by id")
		}
	})
}

func Test_getRole_unhappyPath(t *testing.T) {
	t.Run("getById_notfound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		_, err := container.RoleService.Get(ctx, &GetRequest{
			Id: rand.Int(),
		})
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}

	})

	t.Run("getByName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		// Get by Name
		in := GetRequest{
			Name: "test",
		}
		wantErr := ErrEmptyGetContent
		got, err := container.RoleService.Get(ctx, &in)
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
		if got.Id != 0 {
			t.Errorf("Failed, non existing resource found")
		}
	})
}
