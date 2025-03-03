package role

import (
	"context"
	"math/rand"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/service/resource"
)

var testContainer TestContainer
var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = NewTestContainer()
	service = testContainer.RoleService
}

func Test_create_happyPath(t *testing.T) {
	ctx := context.Background()
	in := CreateRequest{
		Name: "test",
		Desc: "test",
	}

	_, err := service.Create(ctx, &in)
	if err != nil {
		t.Errorf("Failed to create role err: %v", err)
	}

	//get resource
	getResp, _ := service.GetAll(ctx)
	if len(getResp.List) == 0 {
		t.Errorf("No resources were created for role")
	}
}

func Test_create_unhappyPath(t *testing.T) {
	t.Run("duplicateName", func(t *testing.T) {
		//setup
		ctx := context.Background()
		in := CreateRequest{
			Desc: "test",
			Name: "test",
		}
		_, err := service.Create(ctx, &in)
		if err != nil {
			t.Errorf("Failed to create role err: %v", err)
		}

		//create duplicate
		wantErr := ErrDuplicateName
		got, err := service.Create(ctx, &in)
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
		ctx := context.Background()
		in := CreateRequest{
			Desc: "test",
			Name: "",
		}

		wantErr := ErrEmptyName
		got, err := service.Create(ctx, &in)
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
	//setup
	ctx := context.Background()
	in := CreateRequest{
		Desc: "test",
		Name: "test",
	}
	id, _ := service.Create(ctx, &in)

	//delete resource
	err := service.Delete(ctx, id)
	if err != nil {
		t.Errorf("Failed to delete role err: %v", err)
	}

	//verify deletion
	resp, _ := service.Get(ctx, &GetRequest{
		Id:   id,
		Name: "",
	})
	if resp.Id == id {
		t.Error("Failed to delete role")
	}
}

func Test_delete_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		err := service.Delete(ctx, 1010)
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
	//setup
	ctx := context.Background()
	id, _ := service.Create(ctx, &CreateRequest{
		Desc: "test",
		Name: "test",
	})

	//update
	in := UpdateRequest{
		Id:   id,
		Desc: "test",
		Name: "tests",
	}
	resp, err := service.Update(ctx, &in)
	if err != nil {
		t.Errorf("Failed to update err %v", err)
	}
	if id != resp {
		t.Errorf("Failed to update role")
	}
}

func Test_update_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		in := UpdateRequest{
			Id:   1,
			Desc: "test",
			Name: "test",
		}
		wantErr := ErrIdNotFound
		resp, err := service.Update(ctx, &in)
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
		// setup
		ctx := context.Background()

		//create resource with taken name
		service.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "taken",
		})

		// create resource
		id, _ := service.Create(ctx, &CreateRequest{
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
		got, err := service.Update(ctx, &in)
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
		// setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &CreateRequest{
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
			t.Errorf("Failed to update role")
		}
	})

}

func Test_getRole_happyPath(t *testing.T) {
	t.Run("getById", func(t *testing.T) {
		//setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// Get by Id
		in := GetRequest{
			Id: id,
		}
		got, err := service.Get(ctx, &in)
		if err != nil {
			t.Errorf("Failed to get role by Id err %v", err)
		}
		if got.Id != id {
			t.Errorf("Failed to get role by id")
		}
	})

	t.Run("getByName", func(t *testing.T) {
		//setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// Get by Name
		in := GetRequest{
			Name: "test",
		}
		got, err := service.Get(ctx, &in)
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
		//setup
		ctx := context.Background()

		// Get by Id
		in := GetRequest{
			Id: rand.Int(),
		}
		wantErr := ErrEmptyGetContent
		got, err := service.Get(ctx, &in)
		if err != wantErr {
			switch err {
			case wantErr:
			default:
				t.Errorf("Expected err: %v Got err %v", wantErr, err)
			}

		}
		if got.Id != 0 {
			t.Errorf("Failed, non existing resource found")
		}
	})

	t.Run("getByName", func(t *testing.T) {
		//setup
		ctx := context.Background()
		// Get by Name
		in := GetRequest{
			Name: "test",
		}
		wantErr := ErrEmptyGetContent
		got, err := service.Get(ctx, &in)
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
		if got.Id != 0 {
			t.Errorf("Failed, non existing resource found")
		}
	})
}

func Test_getAllResources_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	id, _ := service.Create(ctx, &CreateRequest{
		Desc: "test",
		Name: "test",
	})

	// Get All
	got, err := service.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get role by Id err %v", err)
	}
	if len(got.List) == 0 || got.List[0].Id != id {
		t.Errorf("Failed to get role by id")
	}
}

func Test_hasResource_happyPath(t *testing.T) {
	t.Run("has", func(t *testing.T) {
		//setup
		ctx := context.Background()
		roleId, _ := service.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "test",
		})

		// create resource with taken name
		resId, err := testContainer.ResourceService.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "taken",
		})
		if err != nil {
			t.Errorf("Failed")
		}

		err = service.AddResource(ctx, &AddResourceRequest{
			ResourceId: resId,
			RoleId:     roleId,
		})
		if err != nil {
			t.Errorf("Failed to add resource")
		}

		in := HasResourceRequest{
			ResourceId: resId,
			RoleId:     roleId,
		}
		got, _ := service.HasResource(ctx, &in)
		if got != true {
			t.Errorf("Expected: %v Got: %v", true, got)
		}
	})

}

func Test_hasResource_unhappyPath(t *testing.T) {
	t.Run("has-not", func(t *testing.T) {
		//setup
		ctx := context.Background()
		id, _ := service.Create(ctx, &CreateRequest{
			Desc: "test",
			Name: "test",
		})

		in := HasResourceRequest{
			ResourceId: rand.Intn(200),
			RoleId:     id,
		}
		got, err := service.HasResource(ctx, &in)
		if err != nil {
			t.Errorf("Failed to check role resources %v", err)
		}
		if got != false {
			t.Errorf("Expected: %v Got: %v", false, got)
		}
	})
}

func Test_addResource_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	// create resource with taken name
	resId, err := testContainer.ResourceService.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "taken",
	})
	if err != nil {
		t.Errorf("Failed")
	}

	//create role
	in := CreateRequest{
		Name: "test2",
		Desc: "test",
	}
	roleId, err := service.Create(ctx, &in)
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}

	//add role to resource
	req := AddResourceRequest{
		ResourceId: resId,
		RoleId:     roleId,
	}
	err = service.AddResource(ctx, &req)
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}
	inRes := HasResourceRequest{
		ResourceId: req.ResourceId,
		RoleId:     req.RoleId,
	}
	hasResource, err := service.HasResource(ctx, &inRes)
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}
	if !hasResource {
		t.Errorf("Failed to add resource to role, resource not added in role")
	}
}

func Test_addResource_unhappyPath(t *testing.T) {
	t.Run("resourceNotFound", func(t *testing.T) {
		ctx := context.Background()
		//create role
		in := CreateRequest{
			Name: "test2",
			Desc: "test",
		}
		roleId, err := service.Create(ctx, &in)
		if err != nil {
			t.Errorf("Failed to add resource to role %v", err)
		}

		//add role to resource
		req := AddResourceRequest{
			ResourceId: rand.Int(),
			RoleId:     roleId,
		}
		err = service.AddResource(ctx, &req)
		wantErr := ErrResourceNotFound
		if err != ErrResourceNotFound {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_removeResource_happyPath(t *testing.T) {
	// setup
	ctx := context.Background()

	// create resource
	resId, err := testContainer.ResourceService.Create(ctx, &resource.CreateRequest{
		Action: "test",
		Name:   "test",
	})
	if err != nil {
		t.Errorf("Failed")
	}

	// create role
	roleId, err := service.Create(ctx, &CreateRequest{
		Name: "test2",
		Desc: "test",
	})
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}

	// add role to resource
	req := AddResourceRequest{
		ResourceId: resId,
		RoleId:     roleId,
	}
	err = service.AddResource(ctx, &req)
	if err != nil {
		t.Errorf("Failed to add resource to role %v", err)
	}

	// check if role has resource
	hasResource, err := service.HasResource(ctx, &HasResourceRequest{
		ResourceId: req.ResourceId,
		RoleId:     req.RoleId,
	})
	if err != nil {
		t.Errorf("Failed %v", err)
	}
	if !hasResource {
		t.Errorf("Failed")
	}

	//remove resource
	err = service.RemoveResource(ctx, &RemoveResourceRequest{
		ResourceId: req.ResourceId,
		RoleId:     req.RoleId,
	})
	if err != nil {
		t.Errorf("Failed to remove resource from role %v", err)
	}

	//recheck
	hasResource, err = service.HasResource(ctx, &HasResourceRequest{
		ResourceId: req.ResourceId,
		RoleId:     req.RoleId,
	})
	if err != nil {
		t.Errorf("Failed %v", err)
	}
	if hasResource {
		t.Errorf("Failed to remove resource from role")
	}
}

func Test_removeResource_unhappyPath(t *testing.T) {
	t.Run("resource_not_found", func(t *testing.T) {
		// setup
		ctx := context.Background()

		// create resource
		resId, err := testContainer.ResourceService.Create(ctx, &resource.CreateRequest{
			Action: "test",
			Name:   "test",
		})
		if err != nil {
			t.Errorf("Failed")
		}

		// create role
		roleId, err := service.Create(ctx, &CreateRequest{
			Name: "test2",
			Desc: "test",
		})
		if err != nil {
			t.Errorf("Failed to add resource to role %v", err)
		}

		// add role to resource
		req := AddResourceRequest{
			ResourceId: resId,
			RoleId:     roleId,
		}
		err = service.AddResource(ctx, &req)
		if err != nil {
			t.Errorf("Failed to add resource to role %v", err)
		}

		// check if role has resource
		hasResource, err := service.HasResource(ctx, &HasResourceRequest{
			ResourceId: req.ResourceId,
			RoleId:     req.RoleId,
		})
		if err != nil {
			t.Errorf("Failed %v", err)
		}
		if !hasResource {
			t.Errorf("Failed")
		}

		//remove resource
		err = service.RemoveResource(ctx, &RemoveResourceRequest{
			ResourceId: req.ResourceId,
			RoleId:     req.RoleId,
		})
		if err != nil {
			t.Errorf("Failed to remove resource from role %v", err)
		}

		//recheck
		hasResource, err = service.HasResource(ctx, &HasResourceRequest{
			ResourceId: req.ResourceId,
			RoleId:     req.RoleId,
		})
		if err != nil {
			t.Errorf("Failed %v", err)
		}
		if hasResource {
			t.Errorf("Failed to remove resource from role")
		}

		//remove again
		err = service.RemoveResource(ctx, &RemoveResourceRequest{
			ResourceId: req.ResourceId,
			RoleId:     req.RoleId,
		})
		if err != nil {
			switch err {
			case ErrResourceNotFound:
			default:
				t.Errorf("Failed to remove resource from role %v", err)
			}
		}
	})
}
