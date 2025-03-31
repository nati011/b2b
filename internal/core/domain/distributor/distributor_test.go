package distributor

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"

	authProvider "b2b.nati011.github.com/internal/adapter/secondary/application/auth/provider"
	db "b2b.nati011.github.com/internal/adapter/secondary/application/distributor/db"
	"b2b.nati011.github.com/internal/core/application/auth"
	port "b2b.nati011.github.com/internal/port/application/distributor"
)

var service Provider
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
		"",
	)

	authService = auth.NewAuthService(KeycloakProvider)
	service = NewDistributorService(
		db.NewMock(),

		KeycloakProvider,
	)
}

func Test_Create_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &RegisterDistributorRequest{
			FirstName:       "Test User",
			Email:           "test789@email.com",
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

	in := &RegisterDistributorRequest{
		FirstName:       "Test User",
		Email:           "t6546@email.com",
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
		wantErr := ErrEmptyGetDistributorContent
		_, err := service.GetAll(ctx)
		if err != wantErr {
			t.Errorf("Expected err:%v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	ctx := context.Background()
	in := &RegisterDistributorRequest{
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
		Id: resp.Id,
	}
	_, err = service.GetByParam(ctx, &params)

	if err != nil {
		t.Fatalf("Failed to fetch distributor %v", err)
	}
}

func Test_Get_All_Businesses_unhappyPath(t *testing.T) {
	t.Run("no_businesses_found", func(t *testing.T) {
		ctx := context.Background()
		wantErr := ErrEmptyGetBusinessContent
		_, err := service.GetBusinessAll(ctx)
		if err != wantErr {
			t.Errorf("Expected err:%v Got err: %v", wantErr, err)
		}
	})
}
func Test_Create_Business_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &port.CreateBusinessInformation{
			Name: "Test",
			Tin:  124576,

			GeneralZone:   "Test Zone",
			Region:        "Test Region",
			Woreda:        "Test Woreda",
			DistributorId: rand.Int(),
		}

		_, err := service.AddBusinessInformattion(ctx, in)

		if err != nil {
			log.Fatalf("Create distributor test failed %v", err)
		}
	})
}

func Test_Get_All_Businesses_happyPath(t *testing.T) {
	ctx := context.Background()
	in := &RegisterDistributorRequest{
		FirstName:       "Test User",
		Email:           "tesfh231@gmail.com",
		Password:        "test@123",
		ConfirmPassword: "test@123",
		Username:        "username11",
	}
	_, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create distributor %v", err)
	}
	got, err := service.GetBusinessAll(ctx)
	if err != nil {
		t.Errorf("Expected err:%v Got err: %v", nil, err)
	}
	wantNum := 1
	if len(got.List) != wantNum {
		t.Errorf("Expected len: %v, Got len: %v", wantNum, len(got.List))
	}
}

func Test_Get_Business_happyPath(t *testing.T) {
	ctx := context.Background()

	in := &port.CreateBusinessInformation{
		Name: "Test",
		Tin:  124576,

		GeneralZone:   "Test Zone",
		Region:        "Test Region",
		Woreda:        "Test Woreda",
		DistributorId: rand.Int(),
	}
	resp, err := service.AddBusinessInformattion(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create business %v", err)
	}

	_, err = service.GetById(ctx, resp.BusinessId)

	if err != nil {
		t.Fatalf("Failed to fetch business %v", err)
	}
}

func Test_Update_Business_happyPath(t *testing.T) {
	ctx := context.Background()
	distIn := &RegisterDistributorRequest{
		FirstName:       "Test User",
		Email:           "test475@gmail.com",
		Password:        "test@123",
		ConfirmPassword: "test@123",
		Username:        "username11",
	}
	dist, err := service.Create(ctx, distIn)
	if err != nil {
		t.Fatalf("Failed to create distributor %v", err)
	}
	businessIn := &CreateBusinessInformation{
		Name: "Test",
		Tin:  124576,

		GeneralZone:   "Test Zone",
		Region:        "Test Region",
		Woreda:        "Test Woreda",
		DistributorId: dist.Id,
	}

	businessResp, err := service.AddBusinessInformattion(ctx, businessIn)

	if err != nil {
		t.Fatalf("Failed to create business %v", err)
	}
	in := &port.UpdateBusinessRequest{
		Id:   businessResp.BusinessId,
		Name: "Test",
		Tin:  124576,
	}
	resp, err := service.UpdateBusiness(ctx, in)
	if err != nil {
		t.Fatalf("Failed to update business %v", err)
	}

	want := port.RegisterDistributorResponse{
		Id:      businessResp.BusinessId,
		Message: SUCCESS_MESSAGE,
	}

	if resp != want {
		t.Errorf("Expected: %v, Got: %v", want, resp)
	}
}

func Test_Update_Business_unhappyPath(t *testing.T) {
	ctx := context.Background()
	id := rand.Int()
	in := &port.UpdateBusinessRequest{
		Id:            id,
		Name:          "Test",
		Tin:           124576,
		DistributorId: rand.Int(),
	}
	_, err := service.UpdateBusiness(ctx, in)
	wantErr := ErrEmptyGetDistributorContent
	if err != wantErr {
		t.Errorf("Expected: %v, Got: %v", wantErr, err)
	}
}
