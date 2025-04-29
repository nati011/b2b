package resource

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
var err error

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	container = test_container.NewIntegrationTestContainer(db)
	ctx := context.Background()
	resource_id, err = container.ResourceService.Create(ctx, &resource.CreateRequest{
		Action: "resourceTest",
		Name:   "resourceTest",
	})
	if err != nil {
		panic(err.Error())
	}
}

func teardown() {
	container.TeardownIntegrationTestContainer(db)
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_Read(t *testing.T) {
	t.Run("has_resource", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		role_id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "testResource0012",
			Name: "testResource0012",
		})
		err := container.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to add resouce err %v", err)
		}
		has_resource_status, err := container.RoleService.HasResource(ctx, &role.HasResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to get resource: %v", err)
		}
		wantStatus := true
		if has_resource_status != wantStatus {
			t.Errorf("Expected status: %v Got: %v", wantStatus, has_resource_status)
		}
	})
	t.Run("get_resource", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test2",
		})
		resId, err := container.ResourceService.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "taken_2",
		})
		if err != nil {
			t.Errorf("Failed")
		}
		err = container.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: resId,
			RoleId:     id,
		})
		if err != nil {
			t.Errorf("Failed to add resource")
		}
		// Get All
		got, err := container.RoleService.GetAllResources(ctx, id)
		if err != nil {
			t.Errorf("Failed to get role by Id err %v", err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len:%v Want:%v", wantLen, len(got.List))
		}
	})

}

func Test_Write(t *testing.T) {
	t.Run("add_resource_happyPath", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		role_id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})
		err := container.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to add resouce err %v", err)
		}
		has_resource_status, err := container.RoleService.HasResource(ctx, &role.HasResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to get resource: %v", err)
		}
		wantStatus := true
		if has_resource_status != wantStatus {
			t.Errorf("Expected status: %v Got: %v", wantStatus, has_resource_status)
		}
	})

	t.Run("add_resource_unhappyPath", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		role_id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})
		err := container.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to add resouce err %v", err)
		}

		err = container.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		wantErr := role.ErrResourceAlreadyExistsInRole
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("remove_resource_happyPath", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		role_id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})
		err := container.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to add resouce err %v", err)
		}

		err = container.RoleService.RemoveResource(ctx, &role.RemoveResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to remove resouce err %v", err)
		}
		has_resource_status, err := container.RoleService.HasResource(ctx, &role.HasResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to get resource: %v", err)
		}
		wantStatus := false
		if has_resource_status != wantStatus {
			t.Errorf("Expected status: %v Got: %v", wantStatus, has_resource_status)
		}
	})

	t.Run("remove_resource_unhappyPath", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		//setup
		ctx := context.Background()
		role_id, _ := container.RoleService.Create(ctx, &role.CreateRequest{
			Desc: "test",
			Name: "test",
		})

		err := container.RoleService.RemoveResource(ctx, &role.RemoveResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		var wantErr = role.ErrResourceNotFoundInRole
		if wantErr != err {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}
