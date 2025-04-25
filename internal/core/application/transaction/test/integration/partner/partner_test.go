package partner

import (
	"context"
	"os"
	"testing"

	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
)

var testContainer transaction.TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = transaction.NewPackageIntegrationTestContainer()
}

func Test_Validate_PartnerId_Upon_Transaction_Create_happyPath(t *testing.T) {
	ctx := context.Background()
	//create partner
	id, err := testContainer.PartnerService.Create(ctx, &partner.CreateRequest{
		Name:             "test",
		Icon:             "test",
		Init_payment_url: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create err:%v", err)
	}

	//create transaction
	_, err = testContainer.TransactionService.Create(ctx, &transaction.CreateRequest{
		Amount:     100,
		Partner_Id: id,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
}

func Test_Validate_PartnerId_Upon_Transaction_Create_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	_, err := testContainer.TransactionService.Create(ctx, &transaction.CreateRequest{
		Amount:     100,
		Partner_Id: 99,
	})
	wantErr := transaction.ErrPartnerDoesNotExist
	if err != wantErr {
		t.Errorf("Expected err: %v Got : %v", wantErr, err)
	}
}
