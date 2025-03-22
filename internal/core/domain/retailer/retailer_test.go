package retailer

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/adapter/secondary/application/user/db"
	"b2b.nati011.github.com/internal/core/application/auth"
)

var testContainer TestContainer
var ctx context.Context

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {

	KeycloakUsername := "admin@do.not.edit"
	KeycloakPassword := "W>-553:F?XWXpmV"
	KeycloakRealm := "b2b"
	keycloakApplicationRealm := "b2b"
	keycloakClientId := "733bcd1a-dd25-4e59-b14f-3331872a3d4e"
	keycloakInstanceUrl := "https://euc1.auth.ac/auth"

	KeycloakProvider := authProvider.NewKeycloakProvider(
		keycloakInstanceUrl,
		KeycloakUsername,
		KeycloakPassword,
		KeycloakRealm,
		keycloakApplicationRealm,
		keycloakClientId,
		"",
	)

	authService = auth.NewAuthService(KeycloakProvider)
	service = NewRetailerService(
		db.NewMock(),

		KeycloakProvider,
	)
	ctx = context.Background()
	testContainer = NewPackageIntegrationTestContainer()
}

func Test_Create_happyPath(t *testing.T) {
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
	id, err := testContainer.RetailerService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//check
	resp, err := testContainer.RetailerService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	if resp.Id != id {
		t.Errorf("Expected id: %v Got :%v", id, resp.Id)
	}
}

func Test_Create_unhappyPath(t *testing.T) {
	t.Run("validate_invalid_Tin", func(t *testing.T) {
		// tin :has tobe 10 digits
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
		_, err := testContainer.RetailerService.Create(ctx, &in)
		wantErr := ErrInvalidTin
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("validate lat", func(t *testing.T) {
		//validate lat: has tobe validate coordinates
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
		_, err := testContainer.RetailerService.Create(ctx, &in)
		wantErr := ErrInvalidLatitude
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("validate long", func(t *testing.T) {
		//validate long: has tobe validate coordinates
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
		_, err := testContainer.RetailerService.Create(ctx, &in)
		wantErr := ErrInvalidLatitude
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_Update_happyPath(t *testing.T) {
	t.Run("updateName", func(t *testing.T) {
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
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		update_in := UpdateRequest{
			Id:   id,
			Name: "test",
			Tin:  "1111111111",
		}
		err = testContainer.RetailerService.Update(ctx, &update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		//check
		resp, err := testContainer.RetailerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

	})

	t.Run("updateTin", func(t *testing.T) {

	})
}

func Test_Update_unhappyPath(t *testing.T) {

}

func Test_Get_happyPath(t *testing.T) {
	t.Run("getById", func(t *testing.T) {
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
		id, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		_, err = testContainer.RetailerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
	})

	t.Run("getByName", func(t *testing.T) {
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
		_, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		_, err = testContainer.RetailerService.GetByParam(ctx, &GetByParamRequest{
			Name: in.FirstName + in.LastName,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
	})

	t.Run("getByTin", func(t *testing.T) {
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
		_, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		_, err = testContainer.RetailerService.GetByParam(ctx, &GetByParamRequest{
			Tin: in.Tin,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
	})

	t.Run("getAll", func(t *testing.T) {
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
		_, err := testContainer.RetailerService.Create(ctx, &in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		_, err = testContainer.RetailerService.GetByParam(ctx, &GetByParamRequest{
			Tin: in.Tin,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
	})
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		_, err := testContainer.RetailerService.GetByParam(ctx, &GetByParamRequest{
			Tin: "test",
		})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}
