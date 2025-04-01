package distributor

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/user"
)

var testContainer TestContainer
var ctx context.Context

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	ctx = context.Background()
	testContainer = NewPackageIntegrationTestContainer()
}

func Test_Create_happyPath(t *testing.T) {
	t.Cleanup(testContainer.Cleanup)
	in := CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",

		FirstName: "test",
		LastName:  "test",

		Email: "test@gmail.com",
	}
	id, err := testContainer.DistributorService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//check
	resp, err := testContainer.DistributorService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	if resp.Id != id {
		t.Errorf("Expected id: %v Got :%v", id, resp.Id)
	}
}

func Test_Create_unhappyPath(t *testing.T) {
	t.Run("validate_invalid_Tin", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		// tin :has tobe 10 digits
		in := CreateRequest{
			Tin:         "111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		_, err := testContainer.DistributorService.Create(ctx, &in)
		wantErr := ErrInvalidTin
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("validate_duplicate_Tin", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		//setup
		_, err := testContainer.DistributorService.Create(ctx, &CreateRequest{
			Tin:         "1234567891",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		})
		if err != nil {
			t.Fatalf("Failed to create %v", err)
		}

		in := CreateRequest{
			Tin:         "1234567891",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		_, err = testContainer.DistributorService.Create(ctx, &in)
		wantErr := ErrDuplicateTin
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_Update_happyPath(t *testing.T) {
	t.Run("updateName", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		//setup
		in := CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.DistributorService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		update_in := UpdateRequest{
			Id:   id,
			Name: "test",
		}
		_, err = testContainer.DistributorService.Update(ctx, &update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		//check
		resp, err := testContainer.DistributorService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Name != update_in.Name {
			t.Errorf("Expected name: %v Got:%v", update_in.Name, resp.Name)
		}
	})

	t.Run("updateTin", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		// setup
		in := CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.DistributorService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		update_in := UpdateRequest{
			Id:  id,
			Tin: "1234567891",
		}
		_, err = testContainer.DistributorService.Update(ctx, &update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		// check
		resp, err := testContainer.DistributorService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		if resp.Tin != update_in.Tin {
			t.Errorf("Expected Tin: %v Got:%v", update_in.Tin, resp.Tin)
		}
	})
}

func Test_Update_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		update_in := UpdateRequest{
			Id:   99,
			Name: "test",
			Tin:  "1111111111",
		}
		_, err := testContainer.DistributorService.Update(ctx, &update_in)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("duplicate_Tin", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		//setup
		id, err := testContainer.DistributorService.Create(ctx, &CreateRequest{
			Tin:         "1234567891",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		})
		if err != nil {
			t.Fatalf("Failed to create %v", err)
		}

		update_in := UpdateRequest{
			Id:  id,
			Tin: "1234567891",
		}
		_, err = testContainer.DistributorService.Update(ctx, &update_in)
		wantErr := ErrDuplicateTin
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	t.Run("getById", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		//setup
		in := CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.DistributorService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.DistributorService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.Id)
		}
	})

	t.Run("getByName", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		// setup
		in := CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.DistributorService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.DistributorService.GetByParam(ctx, &GetByParamRequest{
			Name: in.FirstName + in.LastName,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(resp.List))
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})

	t.Run("getByTin", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		// setup
		in := CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.DistributorService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.DistributorService.GetByParam(ctx, &GetByParamRequest{
			Tin: in.Tin,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", len(resp.List), wantLen)
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})
	t.Run("getAll", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		// setup
		in := CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.DistributorService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.DistributorService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", len(resp.List), wantLen)
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.List[0].Id)
		}
	})
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		_, err := testContainer.DistributorService.GetByParam(ctx, &GetByParamRequest{
			Tin: "test",
		})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("empty_getAll", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		_, err := testContainer.DistributorService.GetAll(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_All_Users_happyPath(t *testing.T) {
	t.Run("getAllUsers", func(t *testing.T) {
		//setup
		t.Cleanup(testContainer.Cleanup)
		in := CreateRequest{
			Tin:         "1111111111",
			Latitude:    "9.0192° N",
			Longitude:   "38.7525° E",
			GeneralZone: "test",
			Region:      "test",
			Woreda:      "test",

			FirstName: "test",
			LastName:  "test",
			Email:     "test@gmail.com",
		}
		id, err := testContainer.DistributorService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		_, err = testContainer.DistributorService.GetAllUsers(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get user agents %v", err)
		}

	})
}

func Test_Get_All_Users_unhappyPath(t *testing.T) {

}

func Test_Create_Distributor_user_happyPath(t *testing.T) {
	t.Cleanup(testContainer.Cleanup)
	in := CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",

		FirstName: "test",
		LastName:  "test",
		Email:     "test@gmail.com",
	}
	id, err := testContainer.DistributorService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	user_id, err := testContainer.DistributorService.CreateUser(ctx, &CreateUserRequest{
		Distributor_Id: id,
		FirstName:      "test_user",
		LastName:       "test_user",
		Email:          "test@gmail.com",
	})
	_, err = testContainer.UserService.GetByParam(ctx, &user.GetByParam{
		ID: user_id,
	})
	if err != nil {
		t.Errorf("Failed to get user %v", err)
	}
}

func Test_Create_Distributor_user_unhappyPath(t *testing.T) {
	t.Run("distributorNotFound", func(t *testing.T) {
		t.Cleanup(testContainer.Cleanup)
		_, err := testContainer.DistributorService.CreateUser(ctx, &CreateUserRequest{
			Distributor_Id: 99,
			FirstName:      "test_user",
			LastName:       "test_user",
			Email:          "test@gmail.com",
		})
		WantErr := ErrIdNotFound
		if err != WantErr {
			t.Errorf("Expected err: %v Got: %v", WantErr, err)
		}
	})

}
