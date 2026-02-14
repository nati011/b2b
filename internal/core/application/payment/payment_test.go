package payment

import (
	"context"
	"os"
	"testing"
)

var container TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = NewTestContainer()
}

func tearDown() {
	container.Teardown()
}

func Test_create_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	id, err := container.Service.Create(ctx, &CreateRequest{
		OrderId:        0,
		PartnerId:      1,
		TransactionRef: "test",
		Amount:         100,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := container.Service.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get by id err %v", err)
	}
	if resp.Id != id {
		t.Errorf("Expected id: %v, Got id: %v", resp.Id, id)
	}
}

func Test_create_unhappyPath(t *testing.T) {
	t.Cleanup(tearDown)
}

func Test_getById_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	id, err := container.Service.Create(ctx, &CreateRequest{
		OrderId:        0,
		PartnerId:      1,
		TransactionRef: "test",
		Amount:         100,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := container.Service.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get by id err %v", err)
	}
	if resp.Id != id {
		t.Errorf("Expected id: %v, Got id: %v", resp.Id, id)
	}
}

func Test_getById_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		_, err := container.Service.GetByID(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expeced err: %v, Got: %v", wantErr, err)
		}
	})
}

func Test_getAll_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	_, err := container.Service.Create(ctx, &CreateRequest{
		OrderId:        0,
		PartnerId:      1,
		TransactionRef: "test",
		Amount:         100,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := container.Service.GetAll(ctx)
	if err != nil {
		t.Errorf("Failed to get all err: %v", err)
	}
	wantLen := 1
	if len(resp.List) != wantLen {
		t.Errorf("Expected len: %v, Got %v", wantLen, len(resp.List))
	}
}

func Test_getAll_unhappyPath(t *testing.T) {
	t.Run("emptyGetContent", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		_, err := container.Service.GetAll(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v, Want: %v", wantErr, err)
		}
	})
}

func Test_getByTransactionRef_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	in_txRef := "test"
	id, err := container.Service.Create(ctx, &CreateRequest{
		OrderId:        0,
		PartnerId:      1,
		TransactionRef: in_txRef,
		Amount:         100,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := container.Service.GetByTransactionRef(ctx, in_txRef)
	if err != nil {
		t.Fatalf("Failed to get by id err %v", err)
	}
	if resp.Id != id {
		t.Errorf("Expected id: %v, Got id: %v", resp.Id, id)
	}
}

func Test_getByTransactionRef_unhappyPath(t *testing.T) {
	t.Run("txRefNotFound", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		_, err := container.Service.GetByTransactionRef(ctx, "test")
		wantErr := ErrTxRefNotFound
		if err != wantErr {
			t.Errorf("Expeced err: %v, Got: %v", wantErr, err)
		}
	})
}

func Test_getByOrderId_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	in_orderId := 0
	_, err := container.Service.Create(ctx, &CreateRequest{
		OrderId:        in_orderId,
		PartnerId:      1,
		TransactionRef: "test",
		Amount:         100,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := container.Service.GetByOrderId(ctx, in_orderId)
	if err != nil {
		t.Fatalf("Failed to get by id err %v", err)
	}
	wantLen := 1
	if wantLen != len(resp.List) {
		t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp.List))
	}
}

func Test_getByOrderId_unhappyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	_, err := container.Service.GetByOrderId(ctx, 999)
	wantErr := ErrOrderIdNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v, Want: %v", wantErr, err)
	}
}
