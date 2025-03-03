package distributor

import (
	"context"
	"log"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/distributor/db"
)

var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewDistributorService(
		db.NewMock(),
	)
}

func Test_Create_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &RegisterDistributorRequest{
			FullName:        "Test User",
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
	got, err := service.GetAll(ctx)
	if err != nil {
		t.Errorf("Expected err:%v Got err: %v", nil, err)
	}
	wantNum := 1
	if len(got.List) != wantNum {
		t.Errorf("Expected len: %v, Got len: %v", wantNum, len(got.List))
	}
}

func Test_Get_happyPath(t *testing.T) {
	ctx := context.Background()
	in := &RegisterDistributorRequest{
		FullName:        "Test User",
		Email:           "test11@email.com",
		Password:        "test@123",
		ConfirmPassword: "test@123",
		Username:        "username11",
	}

	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create distributor %v", err)
	}

	_, err = service.Get(ctx, id)

	if err != nil {
		t.Fatalf("Failed to fetch distributor %v", err)
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("no_distributor_found", func(t *testing.T) {
		ctx := context.Background()
		//get
		wantErr := ErrEmptyGetContent
		_, err := service.Get(ctx, 99)
		if err != wantErr {
			t.Errorf("Expected err:%v Got err: %v", wantErr, err)
		}
	})
}
func Test_Add_Business_Information_happyPath(t *testing.T) {
	ctx := context.Background()
	createRequest := &RegisterDistributorRequest{
		FullName:        "Test User",
		Email:           "test11@email.com",
		Password:        "test@123",
		ConfirmPassword: "test@123",
		Username:        "username11",
	}

	id, err := service.Create(ctx, createRequest)
	if err != nil {
		t.Fatalf("Failed to create distributor %v", err)
	}

	region := &BusinessLocation{
		GeneralZone: "Test Zone",
		Region:      "Test Region",
		Woreda:      "Test Woreda",
	}

	in := &UpdateBusinessRequest{
		DistributorId: id,
		Name:          "Test location",
		Tin:           127897024567,
		Region:        *region,
	}

	_, err = service.CreateBusinessInformation(ctx, in)

	if err != nil {
		t.Fatalf("Failed to add business info. Error: %v", err)
	}

}
