package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/resource"
	"b2b.nati011.github.com/internal/core/application/role"
	test_container "b2b.nati011.github.com/internal/core/application/role/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var container test_container.TestContainer
var db *sql.DB
var resource_id int

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	container = test_container.NewIntegrationTestContainer(db)
	ctx := context.Background()
	resource_id, _ = container.ResourceService.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "test",
	})
}

func teardown() {
	container.TeardownIntegrationTestContainer(db)
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_Read(t *testing.T) {

	t.Run("get_by_id", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// Get by Id
		in := role.GetRequest{
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

	t.Run("get_by_name", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// Get by Id
		in := role.GetRequest{
			Name: "test",
		}
		got, err := container.RoleService.Get(ctx, &in)
		if err != nil {
			t.Errorf("Failed to get role by Id err %v", err)
		}
		if got.Id != id {
			t.Errorf("Failed to get role by id")
		}
	})

	t.Run("get_all", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})
		got, err := container.RoleService.GetAll(ctx)
		if err != nil {
			t.Errorf("Failed to get role by Id err %v", err)
		}
		wantLen := 1
		if len(got.List) != wantLen {
			t.Errorf("Expected len: %v Got len %v", wantLen, len(got.List))
		}
	})
}

func Test_write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := role.CreateRequest{
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
	})

	t.Run("update_name", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})

		//update
		update_in := role.UpdateRequest{
			Id:   id,
			Name: "name",
		}
		_, err := container.RoleService.Update(ctx, &update_in)
		if err != nil {
			t.Errorf("Failed to update err %v", err)
		}
		resp, err := container.RoleService.Get(ctx, &role.GetRequest{
			// Id: id,
			Name: update_in.Name,
		})
		if err != nil {
			t.Fatalf("Failed to get role err: %v", err)
		}
		if resp.Name != update_in.Name {
			t.Errorf("Expected name:%v, Got: %v", update_in.Name, resp.Name)
		}
	})

	t.Run("update_desc", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})

		//update
		update_in := role.UpdateRequest{
			Id:   id,
			Desc: "desc",
		}
		_, err := container.RoleService.Update(ctx, &update_in)
		if err != nil {
			t.Errorf("Failed to update err %v", err)
		}
		resp, err := container.RoleService.Get(ctx, &role.GetRequest{
			Id: id,
		})
		if err != nil {
			t.Fatalf("Failed to get role err: %v", err)
		}
		if resp.Desc != update_in.Desc {
			t.Errorf("Expected desc:%v, Got: %v", update_in.Desc, resp.Desc)
		}
	})

	t.Run("delete", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		in := role.CreateRequest{
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
		resp, _ := container.RoleService.Get(ctx, &role.GetRequest{
			Id: id,
		})
		if resp.Id == id {
			t.Errorf("Failed to delete role err :%v", err)
		}
	})
}
