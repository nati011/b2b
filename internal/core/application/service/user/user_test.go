package user

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	role "b2b.nati011.github.com/internal/core/application/service/role"
)

var testContainer TestContainer
var service Provider

func TestMain(m *testing.M, testContainer TestContainer, service Provider) {
	setup(testContainer, service)
	code := m.Run()
	os.Exit(code)
}

func setup(testContainer TestContainer, service Provider) {
	testContainer = NewTestContainer()
	service = testContainer.UserService
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
			ExternalId: "123",
		}
		id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		_, err = testContainer.UserService.GetByParam(ctx, &GetByParam{
			ID: id,
		})
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
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
			ExternalId: "123",
		}
		id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		user, err := testContainer.UserService.GetByParam(ctx, &GetByParam{
			ID: id,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if user.List[0].IsActive != true {
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
			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
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
			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != ErrPhoneOrEmailMandatory {
			t.Errorf("Expected Err: %v Got: %v", ErrPhoneOrEmailMandatory, err)
		}
		//just phone
		in_only_phone := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Username:   "test",
			DOB:        parsedTime,
			Phone:      "+251949184879",
			ExternalId: "123",
		}
		_, err = testContainer.UserService.Create(ctx, &in_only_phone)
		if err != nil {
			t.Errorf("Expected Err: %v Got: %v", nil, err)
		}
		//just email
		in_only_email := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Username:   "test",
			DOB:        parsedTime,
			Phone:      "+251949184879",
			ExternalId: "123",
		}
		_, err = testContainer.UserService.Create(ctx, &in_only_email)
		if err != nil {
			t.Errorf("Expected Err: %v Got: %v", nil, err)
		}
	})

	t.Run("phone_validation", func(t *testing.T) {
		ctx := context.Background()
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Username:   "test",
			DOB:        parsedTime,
			Phone:      "011",
			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		wantErr := ErrPhoneNotValid
		if err != wantErr {
			t.Errorf("Expected Err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("email_validation", func(t *testing.T) {
		ctx := context.Background()
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Username:   "test",
			DOB:        parsedTime,
			Email:      "natnaeljemaneh001gmail.com",
			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		wantErr := ErrEmailNotValid
		if err != wantErr {
			t.Errorf("Expected Err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_getAll_happyPath(t *testing.T) {
	t.Run("non_empty_content", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//get
		resp, err := testContainer.UserService.GetAll(ctx)
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
		resp, err := testContainer.UserService.GetAll(ctx)
		if err != ErrEmptyGetContent {
			t.Errorf("Expected err: %v Got err: %v", ErrEmptyGetContent, err)
		}
		expecetdLen := 0
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
			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		inParam := &GetByParam{
			Email: email,
		}
		response, err := testContainer.UserService.GetByParam(ctx, inParam)
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
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    phone_number,
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		inParam := &GetByParam{
			Phone: phone_number,
		}
		response, err := testContainer.UserService.GetByParam(ctx, inParam)
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
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: username,
			DOB:      parsedTime,

			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		inParam := &GetByParam{
			Username: username,
		}
		response, err := testContainer.UserService.GetByParam(ctx, inParam)
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
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		inParam := &GetByParam{
			IsActive: status,
		}
		response, err := testContainer.UserService.GetByParam(ctx, inParam)
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
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName:   "natnael jemaneh asefa",
			Email:      "natnaeljemaneh001@gmail.com",
			Phone:      "+251949184879",
			Username:   "test",
			DOB:        parsedTime,
			ExternalId: "123",
		}
		_, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		parsedTime, _ = time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in_new := &CreateRequest{
			FullName:   "eyoel jemaneh asefa",
			Email:      "eyoeljemaneh011@gmail.com",
			Phone:      "+251933184880",
			Username:   "test",
			DOB:        parsedTime,
			ExternalId: "123",
		}
		_, err = testContainer.UserService.Create(ctx, in_new)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		inParam := &GetByParam{
			Username: "test",
		}
		response, err := testContainer.UserService.GetByParam(ctx, inParam)
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
		response, err := testContainer.UserService.GetByParam(ctx, inParam)
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
		ExternalId: "123",
	}
	id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = testContainer.UserService.Deactivate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}

	err = testContainer.UserService.Activate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_activate_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName: "natnael jemaneh asefa",
		Email:    "natnaeljemaneh001@gmail.com",
		Phone:    "+251949184879",
		Username: "test",
		DOB:      parsedTime,

		ExternalId: "123",
	}
	id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = testContainer.UserService.Deactivate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}

	err = testContainer.UserService.Activate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}
	expectedErr := ErrUserAlreadyActive
	err = testContainer.UserService.Activate(ctx, id)
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
		ExternalId: "123",
	}
	id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = testContainer.UserService.Deactivate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_deactivate_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName: "natnael jemaneh asefa",
		Email:    "natnaeljemaneh001@gmail.com",
		Phone:    "+251949184879",
		Username: "test",
		DOB:      parsedTime,

		ExternalId: "123",
	}
	id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = testContainer.UserService.Deactivate(ctx, id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}
	expectedErr := ErrUserAlreadyInactive
	err = testContainer.UserService.Deactivate(ctx, id)
	if err != expectedErr {
		t.Fatalf("Expected err: %v Got err: %v", expectedErr, err)
	}
}

func Test_isActive_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName: "natnael jemaneh asefa",
		Email:    "natnaeljemaneh001@gmail.com",
		Phone:    "+251949184879",
		Username: "test",
		DOB:      parsedTime,

		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	isActive, err := testContainer.UserService.IsActive(ctx, user_id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}
	expectedStatus := true
	if !isActive {
		t.Fatalf("Expected: %v Got: %v", expectedStatus, err)
	}

	//deactivate
	err = testContainer.UserService.Deactivate(ctx, user_id)
	if err != nil {
		t.Fatalf("Expected err: %v Got err: %v", nil, err)
	}

	expectedStatus = false
	if !isActive {
		t.Fatalf("Expected: %v Got: %v", expectedStatus, err)
	}
}

func Test_isActive_unhappyPath(t *testing.T) {
	ctx := context.Background()
	expectedErr := ErrIdNotFound
	_, err := testContainer.UserService.IsActive(ctx, rand.Int())
	if err != expectedErr {
		t.Fatalf("Expected err: %v Got err: %v", expectedErr, err)
	}
}

func Test_assignRole_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName: "natnael jemaneh asefa",
		Email:    "natnaeljemaneh001@gmail.com",
		Phone:    "+251949184879",
		Username: "test",
		DOB:      parsedTime,

		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_assignRole_unhappyPath(t *testing.T) {
	t.Run("roleDoesNotExist", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//assign
		err = testContainer.UserService.AssignRole(ctx, user_id, rand.Int())
		expectedErr := ErrRoleDoesNotExist
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

	t.Run("alreadyAssigned", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//create role
		role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//assign
		err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}

		//reassign
		err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
		expectedErr := ErrRoleAlreadyAssigned
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

}

func Test_removeAssignedRole_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName: "natnael jemaneh asefa",
		Email:    "natnaeljemaneh001@gmail.com",
		Phone:    "+251949184879",
		Username: "test",
		DOB:      parsedTime,

		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}

	//remove
	err = testContainer.UserService.RemoveAssignedRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_removeAssignedRole_unhappyPath(t *testing.T) {
	t.Run("roleNotAssigned", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//create role
		role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//remove
		expectedErr := ErrRoleNotAssigned
		err = testContainer.UserService.RemoveAssignedRole(ctx, user_id, role_id)
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

	t.Run("roleNotFound", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//remove
		expectedErr := ErrRoleDoesNotExist
		err = testContainer.UserService.RemoveAssignedRole(ctx, user_id, rand.Int())
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

}

func Test_getAllRole_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName: "natnael jemaneh asefa",
		Email:    "natnaeljemaneh001@gmail.com",
		Phone:    "+251949184879",
		Username: "test",
		DOB:      parsedTime,

		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	//create role
	role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
		Name: "test",
		Desc: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create role err: %v", err)
	}

	//assign
	err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}

	//get all assignment
	resp, err := testContainer.UserService.GetAllAssignedRoles(ctx, user_id)
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
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//create role
		_, err = testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//get all assignment
		wantErr := ErrNoRoleAssigned
		resp, err := testContainer.UserService.GetAllAssignedRoles(ctx, user_id)
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
		expecetdLen := 0
		if len(resp.List) != expecetdLen {
			t.Errorf("Expected len: %v Got len: %v", expecetdLen, err)
		}
	})
}

func Test_has_role_happyPath(t *testing.T) {
	t.Run("has", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//create role
		role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//assign
		err = testContainer.UserService.AssignRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}

		//check access to resource
		hasAccess, err := testContainer.UserService.HasRole(ctx, user_id, role_id)
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
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		//create role
		role_id, err := testContainer.RoleService.Create(ctx, &role.CreateRequest{
			Name: "test",
			Desc: "test",
		})
		if err != nil {
			t.Fatalf("Failed to create role err: %v", err)
		}

		//check access to role
		got, err := testContainer.UserService.HasRole(ctx, user_id, role_id)
		if err != nil {
			t.Errorf("Expected err: %v Got err: %v", nil, err)
		}
		want := false
		if got != want {
			t.Errorf("Expected has access status: %v Got err: %v", want, got)
		}
	})
}

func Test_has_role_unhappyPath(t *testing.T) {
	t.Run("roleDoesNotExist", func(t *testing.T) {
		ctx := context.Background()
		//setup
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
		in := CreateRequest{
			FullName: "natnael jemaneh asefa",
			Email:    "natnaeljemaneh001@gmail.com",
			Phone:    "+251949184879",
			Username: "test",
			DOB:      parsedTime,

			ExternalId: "123",
		}
		user_id, err := testContainer.UserService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check access to resource
		expectedErr := ErrRoleDoesNotExist
		got, err := testContainer.UserService.HasRole(ctx, user_id, rand.Int())
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
		want := false
		if got != want {
			t.Errorf("Expected has access status: %v Got: %v", want, got)
		}
	})
}

func Test_update_user_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	in_update := &UpdateRequest{
		Id:         user_id,
		FullName:   "test",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		ExternalId: "123",
	}
	got, err := testContainer.UserService.Update(ctx, in_update)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
	if got.FullName != in_update.FullName {
		t.Errorf("Expected : %v Got: %v", in_update.DOB, got.DOB)
	}
	if got.Email != in_update.Email {
		t.Errorf("Expected : %v Got: %v", in_update.Email, got.Email)
	}
	if got.Phone != in_update.Phone {
		t.Errorf("Expected : %v Got: %v", in_update.Phone, got.Phone)
	}
	if got.Username != in_update.Username {
		t.Errorf("Expected : %v Got: %v", in_update.Username, got.Username)
	}
	if got.DOB != in_update.DOB {
		t.Errorf("Expected : %v Got: %v", in_update.DOB, got.DOB)
	}
	if got.ExternalId != in_update.ExternalId {
		t.Errorf("Expected : %v Got: %v", in_update.ExternalId, got.ExternalId)
	}
}

func Test_update_user_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")

	in_update := &UpdateRequest{
		Id:         rand.Int(),
		FullName:   "test",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		ExternalId: "123",
	}
	_, err := testContainer.UserService.Update(ctx, in_update)
	expectedErr := ErrIdNotFound
	if err != expectedErr {
		t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
	}
}

func Test_remove_user_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := CreateRequest{
		FullName:   "natnael jemaneh asefa",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		ExternalId: "123",
	}
	user_id, err := testContainer.UserService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	err = testContainer.UserService.Remove(ctx, user_id)
	if err != nil {
		t.Errorf("Expected err: %v Got err: %v", nil, err)
	}
}

func Test_remove_user_unhappyPath(t *testing.T) {
	ctx := context.Background()
	err := testContainer.UserService.Remove(ctx, rand.Int())
	expectedErr := ErrIdNotFound
	if err != expectedErr {
		t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
	}
}
