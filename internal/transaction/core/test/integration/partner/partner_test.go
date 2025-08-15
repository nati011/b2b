package partner

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
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
		Name:          "test",
		Icon:          "test",
		BaseURL:       "test",
		Secret:        "test",
		PaymentMethod: payment_partner.PAYMENT_METHOD_DIGITAL,
	})
	if err != nil {
		t.Fatalf("Failed to create err:%v", err)
	}

	//create transaction
	_, err = testContainer.TransactionService.Create(ctx, &transaction.CreateRequest{
		Amount:    100,
		PartnerId: id,
		TxRef:     "test",
	})
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
}
