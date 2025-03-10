package distributor

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/application/distributor/business/db"
	port "b2b.nati011.github.com/internal/port/distributor/business"
)

var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {

	service = NewDistributorBusinessService(
		db.NewMock(),
	)
}

func Test_Get_All_unhappyPath(t *testing.T) {
	t.Run("no_distributor_found", func(t *testing.T) {
		ctx := context.Background()
		wantErr := ErrEmptyGetContent
		_, err := service.GetAll(ctx)
		if err != wantErr {
			t.Errorf("Expected err:%v Got err: %v", wantErr, err)
		}
	})
}
func Test_Create_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &port.CreateBusinessInformation{
			Name: "Test",
			Tin:  "124576",

			GeneralZone:   "Test Zone",
			Region:        "Test Region",
			Woreda:        "Test Woreda",
			DistributorId: rand.Int(),
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
	in := &port.CreateBusinessInformation{
		Name: "Test",
		Tin:  "124576",

		GeneralZone:   "Test Zone",
		Region:        "Test Region",
		Woreda:        "Test Woreda",
		DistributorId: rand.Int(),
	}
	resp, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create business %v", err)
	}

	_, err = service.GetById(ctx, resp.BusinessId)

	if err != nil {
		t.Fatalf("Failed to fetch business %v", err)
	}
}
