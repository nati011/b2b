package distributor

import (
	"context"
	"log"
	"os"
	"testing"

	authProvider "b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	db "b2b.nati011.github.com/internal/adapter/secondary/application/distributor/db"
	"b2b.nati011.github.com/internal/core/application/auth"
	distributor "b2b.nati011.github.com/internal/core/domain/distributor/service"
	port "b2b.nati011.github.com/internal/port/distributor"
)

var service distributor.Provider
var authService auth.Provider

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
	)

	authService = auth.NewAuthService(KeycloakProvider)
	service = distributor.NewDistributorService(
		db.NewMock(),

		KeycloakProvider,
	)
}

func Test_Create_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &port.RegisterDistributorRequest{
			FirstName:       "Test User",
			Email:           "test@email.com",
			Password:        "test@123",
			ConfirmPassword: "test@123",
			Username:        "username",
		}

		_, err := service.Create(ctx, in)

		if err != nil {
			log.Fatalf("Create distributor test failed %v", err)
		}
	})
}

func Test_Get_All_happyPath(t *testing.T) {
	ctx := context.Background()

	in := &port.RegisterDistributorRequest{
		FirstName:       "Test User",
		Email:           "test116546@email.com",
		Password:        "test@123",
		ConfirmPassword: "test@123",
		Username:        "username11",
	}

	_, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create distributor %v", err)
	}
	got, err := service.GetAll(ctx)
	if err != nil {
		t.Errorf("Expected err:%v Got err: %v", nil, err)
	}
	wantNum := 1
	if len(got.List) != wantNum {
		t.Errorf("Expected len: %v, Got len: %v", wantNum, len(got.List))
	}
}

func Test_Get_All_unhappyPath(t *testing.T) {
	t.Run("no_distributor_found", func(t *testing.T) {
		ctx := context.Background()
		wantErr := distributor.ErrEmptyGetContent
		_, err := service.GetAll(ctx)
		if err != wantErr {
			t.Errorf("Expected err:%v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	ctx := context.Background()
	in := &port.RegisterDistributorRequest{
		FirstName:       "Test User",
		Email:           "tesfhjt2_11@gmail.com",
		Password:        "test@123",
		ConfirmPassword: "test@123",
		Username:        "username11",
	}

	resp, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create distributor %v", err)
	}
	params := port.GetByParamRequest{
		Id: resp.DistributorId,
	}
	_, err = service.GetByParam(ctx, &params)

	if err != nil {
		t.Fatalf("Failed to fetch distributor %v", err)
	}
}
