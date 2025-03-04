package partner

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment/partner"
	"b2b.nati011.github.com/internal/core/application/payment/transaction"
	"b2b.nati011.github.com/internal/core/application/payment/transaction/test/integration"
)

var testContainer integration.TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = integration.NewPackageIntegrationTestContainer()
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
		User_Id:    1,
		Amount:     100,
		Partner_Id: id,
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
}

func Test_Validate_PartnerId_Upon_Fetch_By_PartnerId_Request_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//create transaction
	_, err := testContainer.TransactionService.Create(ctx, &transaction.CreateRequest{
		User_Id:    1,
		Amount:     100,
		Partner_Id: 99,
	})
	wantErr := transaction.ErrPartnerNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v Got : %v", wantErr, err)
	}
}
