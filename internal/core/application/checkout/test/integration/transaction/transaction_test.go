package transaction

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/checkout/test"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
)

var testContainer test.TestContainer
var PaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = test.NewPackageIntegrationTestContainer()
	ctx := context.Background()
	var err error
	PaymentPartnerId, err = testContainer.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:    "chapa",
			Icon:    "test",
			BaseURL: "https://api.chapa.co",
			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		})

	if err != nil {
		panic("failed to create order partner")
	}
}

func Test_CreateTransactionUponPaymentInitAndSetStatusToPending(t *testing.T) {
	ctx := context.Background()

	//init transaction
	in := &checkout.CheckoutRequest{
		PaymentPartnerId: PaymentPartnerId,
		OrderId:          1,
		Amount:           100,
	}

	checkout_response, err := testContainer.CheckoutService.Checkout(ctx, in)
	if err != nil {
		t.Errorf("Failed to checkout err: %v", err)
	}

	//check if transaction has been created
	resp, err := testContainer.TransactionService.GetByParam(ctx, &transaction.GetByParamRequest{
		TxRef: checkout_response.TransactionRef,
	})
	if err != nil {
		t.Fatalf("failed to get param err: %v", err)
	}
	wantStatus := transaction.PENDING_STATUS
	if resp.List[0].Status != wantStatus {
		t.Errorf("expected status: %v, got: %v", resp.List[0].Status, wantStatus)
	}
}
