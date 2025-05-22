package authorization

import (
	"context"
	"os"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/application/resource"
	role "b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/user"
)

var testContainer TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = NewTestContainer()
}

func teardown() {
	testContainer.teardown()
}

func Test_check_if_user_is_authorized(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//create user
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName:  "natnael asefa",
		LastName:   "jemaneh",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		ExternalId: "123",
		Password:   "test",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//create resource
	resource_in := resource.CreateRequest{
		Action: "test",
		Name:   "test",
	}
	resource_id, err := testContainer.ResourceService.Create(ctx, &resource_in)
	if err != nil {
		t.Errorf("Failed to create resource err: %v", err)
	}

	//create role
	role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	err = testContainer.RoleService.AddResource(ctx, &role.AddResourceRequest{
		ResourceId: resource_id,
		RoleId:     role_id})
	if err != nil {
		t.Errorf("Failed to add resource err: %v", err)
	}
	//give user role
	err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
	//check user authorization
	isAuthorized, err := testContainer.AuthorizationService.IsAuthorizedForResource(ctx, user_id, resource_id)
	if err != nil {
		t.Errorf("Failed to check err: %v", err)
	}
	wantResponse := true
	if isAuthorized != wantResponse {
		t.Errorf("Expected response: %v, Got: %v", wantResponse, isAuthorized)
	}
}

func Test_get_user_authorization(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//create user
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName:  "natnael asefa",
		LastName:   "jemaneh",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		ExternalId: "123",
		Password:   "test",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//create resource
	resource_in := resource.CreateRequest{
		Action: "test",
		Name:   "test",
	}
	resource_id, err := testContainer.ResourceService.Create(ctx, &resource_in)
	if err != nil {
		t.Errorf("Failed to create resource err: %v", err)
	}
	role_in := &role.CreateRequest{
		Name: "test",
		Desc: "test",
	}
	//create role
	role_id, err := testContainer.RoleService.Create(ctx, role_in)
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	err = testContainer.RoleService.AddResource(ctx, &role.AddResourceRequest{
		ResourceId: resource_id,
		RoleId:     role_id})
	if err != nil {
		t.Errorf("Failed to add resource err: %v", err)
	}
	//give user role
	err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}

	gotUserRole, err := testContainer.AuthorizationService.GetUserAuthorization(ctx, user_id)
	if err != nil {
		t.Fatalf("Failed to get user resources: %v", err)
	}
	wantUserRole := ResourceAccess{
		Roles: []string{
			role_in.Name},
	}
	if len(gotUserRole.Roles) == 0 {
		t.Fatalf("Expected resources lenght greater than 0")
	}

	if wantUserRole.Roles[0] != wantUserRole.Roles[0] {
		t.Errorf("Expected role:")
	}
}
