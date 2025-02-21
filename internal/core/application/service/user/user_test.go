package user

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/application/service/resource"
	role "b2b.nati011.github.com/internal/core/application/service/role"
)

var test_container TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
}

func Test_create_happyPath(t *testing.T) {

	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		id, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		_, err = test_container.UserService.Get(ctx, id)
		if err != ErrIdNotFound {
			t.Errorf("Expected err: %v Got err: %v", ErrIdNotFound, err)
		}
	})

	t.Run("user_active_by_default", func(t *testing.T) {
		ctx := context.Background()
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		id, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		user, err := test_container.UserService.Get(ctx, id)
		if err != nil {
			switch err {
			case ErrIdNotFound:
			default:
				t.Fatalf("Failed to get err: %v", err)
			}
		}
		if user.IsActive != true {
			t.Errorf("Expected user active status: %v Got: %v", true, false)
		}
	})
}

func Test_create_unhappyPath(t *testing.T) {
	t.Run("fullName_mandatory", func(t *testing.T) {
		ctx := context.Background()
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		_, err := test_container.UserService.Create(ctx, &in)
		if err != ErrFullNameMandatory {
			t.Errorf("Expected Err: %v Got: %v", ErrFullNameMandatory, err)
		}
	})
	t.Run("phone_or_email_mandatory", func(t *testing.T) {
		//none
		ctx := context.Background()
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		_, err := test_container.UserService.Create(ctx, &in)
		if err != ErrPhoneOrEmailMandatory {
			t.Errorf("Expected Err: %v Got: %v", ErrPhoneOrEmailMandatory, err)
		}
		//just phone
		in_only_phone := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		_, err = test_container.UserService.Create(ctx, &in_only_phone)
		if err != nil {
			t.Errorf("Expected Err: %v Got: %v", nil, err)
		}
		//just email
		in_only_email := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		_, err = test_container.UserService.Create(ctx, &in_only_email)
		if err != nil {
			t.Errorf("Expected Err: %v Got: %v", nil, err)
		}
	})
}

func Test_getAll_happyPath(t *testing.T) {
	t.Run("non_empty_content", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		_, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//get
		resp, err := test_container.UserService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to getAll err: %v", err)
		}
		expecetdLen := 1
		if len(resp.List) == 0 {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(resp.List))
		}
	})
}

func Test_getAll_unhappyPath(t *testing.T) {
	t.Run("empty_content", func(t *testing.T) {
		ctx := context.Background()
		resp, err := test_container.UserService.GetAll(ctx)
		if err != ErrEmptyGetContent {
			t.Errorf("Expected err: %v Got err: %v", ErrEmptyGetContent, err)
		}
		expecetdLen := 1
		if len(resp.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(resp.List))
		}
	})
}

func Test_getByParam_happyPath(t *testing.T) {

	t.Run("email", func(t *testing.T) {
		ctx := context.Background()
		email := "natnaeljemaneh001@gmail.com"
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      email,
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		_, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		inParam := &GetByParam{
			Email: email,
		}
		response, err := test_container.UserService.GetByParam(ctx, inParam)
		if err != nil {
			t.Fatalf("Failed to get user by param err: %v", err)
		}
		expecetdLen := 1
		if len(response.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(response.List))
		}
	})

	t.Run("phone", func(t *testing.T) {
		ctx := context.Background()
		phone_number := "+251949184879"
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      phone_number,
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		_, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		inParam := &GetByParam{
			Phone: phone_number,
		}
		response, err := test_container.UserService.GetByParam(ctx, inParam)
		if err != nil {
			t.Fatalf("Failed to get user by param err: %v", err)
		}
		expecetdLen := 1
		if len(response.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(response.List))
		}
	})

	t.Run("username", func(t *testing.T) {
		ctx := context.Background()
		username := "test"
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   username,
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		_, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		inParam := &GetByParam{
			Username: username,
		}
		response, err := test_container.UserService.GetByParam(ctx, inParam)
		if err != nil {
			t.Fatalf("Failed to get user by param err: %v", err)
		}
		expecetdLen := 1
		if len(response.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(response.List))
		}
	})

	t.Run("active_status", func(t *testing.T) {
		ctx := context.Background()
		status := true
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   status,
			ExternalId: "123",
		}
		_, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		inParam := &GetByParam{
			IsActive: status,
		}
		response, err := test_container.UserService.GetByParam(ctx, inParam)
		if err != nil {
			t.Fatalf("Failed to get user by param err: %v", err)
		}
		expecetdLen := 1
		if len(response.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(response.List))
		}
	})

	t.Run("aggregate_fetch", func(t *testing.T) {
		ctx := context.Background()
		//setup
		status := true
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   status,
			ExternalId: "123",
		}
		_, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		parsedTime, _ = time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in = CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh011@gmail.com",
			Phone:      "+251949184880",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   status,
			ExternalId: "123",
		}
		_, err = test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		inParam := &GetByParam{
			IsActive: status,
		}
		response, err := test_container.UserService.GetByParam(ctx, inParam)
		if err != nil {
			t.Fatalf("Failed to get user by param err: %v", err)
		}
		expecetdLen := 2
		if len(response.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(response.List))
		}
	})
}

func Test_getByParam_unhappyPath(t *testing.T) {
	t.Run("empty_content", func(t *testing.T) {
		ctx := context.Background()
		inParam := &GetByParam{
			Email: "test",
		}
		response, err := test_container.UserService.GetByParam(ctx, inParam)
		if err != ErrEmptyGetContent {
			t.Fatalf("Expected err: %v Got err: %v", ErrEmptyGetContent, err)
		}
		expecetdLen := 0
		if len(response.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, len(response.List))
		}
	})
}

func Test_activate_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		IsActive:   false,
		ExternalId: "123",
	}
	id, err := test_container.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = test_container.UserService.Activate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_activate_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		IsActive:   false,
		ExternalId: "123",
	}
	id, err := test_container.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = test_container.UserService.Activate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}
	expectedErr := ErrUserAlreadyActive
	err = test_container.UserService.Activate(ctx, id)
	if err != expectedErr {
		t.Fatalf("Expected err: %v Got err: %v", expectedErr, err)
	}
}

func Test_deactivate_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		IsActive:   true,
		ExternalId: "123",
	}
	id, err := test_container.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = test_container.UserService.Deactivate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_deactivate_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		IsActive:   true,
		ExternalId: "123",
	}
	id, err := test_container.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = test_container.UserService.Deactivate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}
	expectedErr := ErrUserAlreadyActive
	err = test_container.UserService.Deactivate(ctx, id)
	if err != expectedErr {
		t.Fatalf("Expected err: %v Got err: %v", expectedErr, err)
	}
}

func Test_assignRole_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		IsActive:   true,
		ExternalId: "123",
	}
	user_id, err := test_container.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := test_container.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = test_container.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_assignRole_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		IsActive:   true,
		ExternalId: "123",
	}
	user_id, err := test_container.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//assign
	err = test_container.UserService.AssignRole(ctx, user_id, rand.Int())
	expectedErr := ErrRoleDoesNotExist
	if err != expectedErr {
		t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
	}
}

func Test_removeAssignedRole_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		IsActive:   true,
		ExternalId: "123",
	}
	user_id, err := test_container.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := test_container.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = test_container.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}

	//remove
	err = test_container.UserService.RemoveAssignedRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_removeAssignedRole_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		IsActive:   true,
		ExternalId: "123",
	}
	user_id, err := test_container.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := test_container.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = test_container.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}

	//remove
	expectedErr := ErrRoleNotAssigned
	err = test_container.UserService.RemoveAssignedRole(ctx, user_id, role_id)
	if err != expectedErr {
		t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
	}
}

func Test_getAllRole_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		IsActive:   true,
		ExternalId: "123",
	}
	user_id, err := test_container.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := test_container.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = test_container.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}

	//get all assignment
	resp, err := test_container.UserService.GetAllAssignedRoles(ctx, user_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
	expecetdLen := 1
	if len(resp.List) != expecetdLen {
		t.Errorf("Expected len: %v Got len: %v", expecetdLen, err)
	}
}

func Test_getAllRole_unhappyPath(t *testing.T) {
	t.Run("emptyRoleList", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		user_id, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//create role
		_, err = test_container.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//get all assignment
		wantErr := ErrNoRoleAssigned
		resp, err := test_container.UserService.GetAllAssignedRoles(ctx, user_id)
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
		expecetdLen := 0
		if len(resp.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, err)
		}
	})
}

func Test_has_access_to_resource_happyPath(t *testing.T) {
	t.Run("has", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		user_id, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//create role
		role_id, err := test_container.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//create resource
		resource_id, err := test_container.ResourceService.Create(ctx, &resource.CreateRequest{
			Name:   "test",
			Action: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//assign resource
		err = test_container.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to assign resource to role err: %v", err)
		}

		//assign
		err = test_container.UserService.AssignRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}

		//check access to resource
		hasAccess, err := test_container.UserService.HasAccessToResource(ctx, user_id, resource_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}
		expectedStatus := true
		if !hasAccess {
			t.Errorf("Expected has access status: %v Got err: %v", expectedStatus, hasAccess)
		}
	})

	t.Run("has-not", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			IsActive:   true,
			ExternalId: "123",
		}
		user_id, err := test_container.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//create role
		role_id, err := test_container.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//create resource
		resource_id, err := test_container.ResourceService.Create(ctx, &resource.CreateRequest{
			Name:   "test",
			Action: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//assign resource
		err = test_container.RoleService.AddResource(ctx, &role.AddResourceRequest{
			ResourceId: resource_id,
			RoleId:     role_id,
		})
		if err != nil {
			t.Fatalf("Failed to assign resource to role err: %v", err)
		}

		//assign
		err = test_container.UserService.AssignRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}

		//remove
		err = test_container.UserService.RemoveAssignedRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}

		//check access to resource
		hasAccess, err := test_container.UserService.HasAccessToResource(ctx, user_id, resource_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}
		expectedStatus := false
		if !hasAccess {
			t.Errorf("Expected has access status: %v Got err: %v", expectedStatus, hasAccess)
		}
	})
}

func Test_has_access_to_resource_unhappyPath(t *testing.T) {
	t.Run("resource_not_found", func(t *testing.T) {
	})
}
