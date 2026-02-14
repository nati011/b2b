package payment

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment"
	"b2b.nati011.github.com/internal/core/application/payment/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var testContainer test.TestContainer
var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	testContainer = test.NewDBIntegrationTestContainer(db)
}

func teardown() {
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_Read(t *testing.T) {
	t.Run("getById", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		id, err := testContainer.Service.Create(ctx, &payment.CreateRequest{
			OrderId:        0,
			PartnerId:      1,
			TransactionRef: "test",
			Amount:         100,
		})
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.Service.GetByID(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get by id err %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v, Got id: %v", resp.Id, id)
		}
	})

	t.Run("getByOrderId", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in_orderId := 0
		_, err := testContainer.Service.Create(ctx, &payment.CreateRequest{
			OrderId:        in_orderId,
			PartnerId:      1,
			TransactionRef: "test",
			Amount:         100,
		})
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.Service.GetByOrderId(ctx, in_orderId)
		if err != nil {
			t.Fatalf("Failed to get by id err %v", err)
		}
		wantLen := 1
		if wantLen != len(resp.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp.List))
		}
	})

	t.Run("getByTxRef", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in_txRef := "test"
		id, err := testContainer.Service.Create(ctx, &payment.CreateRequest{
			OrderId:        0,
			PartnerId:      1,
			TransactionRef: in_txRef,
			Amount:         100,
		})
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.Service.GetByTransactionRef(ctx, in_txRef)
		if err != nil {
			t.Fatalf("Failed to get by id err %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v, Got id: %v", resp.Id, id)
		}
	})
}

func Test_Write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		id, err := testContainer.Service.Create(ctx, &payment.CreateRequest{
			OrderId:        0,
			PartnerId:      1,
			TransactionRef: "test",
			Amount:         100,
		})
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.Service.GetByID(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get by id err %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v, Got id: %v", resp.Id, id)
		}
	})
}
