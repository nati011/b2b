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
}

func Test_create_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	id, err := container.service.Create(ctx, CreateRequest{
		OrderId:        0,
		PartnerId:      1,
		TransactionRef: "test",
		Amount:         100,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := container.service.GetByID(ctx, id)
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
	id, err := container.service.Create(ctx, CreateRequest{
		OrderId:        0,
		PartnerId:      1,
		TransactionRef: "test",
		Amount:         100,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := container.service.GetByID(ctx, id)
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
		resp, err := container.service.GetByID(ctx, 99)
		if err != nil {
			t.Fatalf("Failed to get by id err %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v, Got id: %v", resp.Id, id)
		}
	})
}

func Test_getAll_happyPath(t *testing.T) {

}

func Test_getAll_unhappyPath(t *testing.T) {

}

func Test_getByTransactionRef_happyPath(t *testing.T) {

}

func Test_getByTransactionRef_unhappyPath(t *testing.T) {

}

func Test_getByOrderId_happyPath(t *testing.T) {

}

func Test_getByOrderId_unhappyPath(t *testing.T) {

}
