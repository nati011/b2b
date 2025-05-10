package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	test_container "b2b.nati011.github.com/internal/core/application/payment_partner/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var container test_container.TestContainer
var db *sql.DB

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	container = test_container.NewIntegrationTestContainer(db)
}

func teardown() {
	container.TeardownIntegrationTestContainer(db)
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_Read(t *testing.T) {
	t.Run("get_by_id", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &payment_partner.CreateRequest{
			Name:    "test",
			Icon:    "test",
			Status:  "test",
			BaseURL: "https://google.com",
			Secret:  "randomSecret",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := container.PartnerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.Id)
		}
	})

	t.Run("get_all", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &payment_partner.CreateRequest{
			Name:    "test",
			Icon:    "test",
			Status:  "test",
			BaseURL: "test",
			Secret:  "randomSecret",
		}

		_, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		resp, err := container.PartnerService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to get all payment options err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp.List))
		}
	})

	t.Run("get_by_status", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &payment_partner.CreateRequest{
			Name:    "test4",
			Icon:    "test",
			Status:  "test",
			BaseURL: "test",
			Secret:  "randomSecret",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := container.PartnerService.GetActive(ctx)
		if err != nil {
			t.Fatalf("Failed to get all payment options err: %v", err)
		}
		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp.List))
		}

		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got id: %v", resp.List[0].Id, id)
		}
	})

	t.Run("get_client_secret", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &payment_partner.CreateRequest{
			Name:    "test1",
			Icon:    "test",
			Status:  "test",
			BaseURL: "test",
			Secret:  "randomSecret",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		_, err = container.PartnerService.GetPartnerSecret(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}

	})

	t.Run("get_by_name", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &payment_partner.CreateRequest{
			Name:    "test",
			Icon:    "test",
			Status:  "test",
			BaseURL: "test",
			Secret:  "randomSecret",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := container.PartnerService.GetByParam(ctx, &payment_partner.GetByParamRequest{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}
		if resp.List[0].Id != id {
			t.Errorf("Expected id: %v Got id: %v", id, resp.List[0].Id)
		}
	})
}

func Test_Write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &payment_partner.CreateRequest{
			Name:    "test",
			Icon:    "test",
			Status:  "test",
			BaseURL: "https://google.com",
			Secret:  "randomSecret",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := container.PartnerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.Id)
		}
	})

	t.Run("update_status", func(t *testing.T) {
		setup()
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &payment_partner.CreateRequest{
			Name:    "test",
			Icon:    "test",
			Status:  "test",
			BaseURL: "test",
			Secret:  "randomSecret",
		}

		id, err := container.PartnerService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//activate
		err = container.PartnerService.Activate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to activate payment option err: %v", err)
		}

		got, err := container.PartnerService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		if got.Status != payment_partner.ACTIVE_STATUS {
			t.Errorf("Expected status: %v Got: %v", got.Status, payment_partner.ACTIVE_STATUS)
		}
	})
}
