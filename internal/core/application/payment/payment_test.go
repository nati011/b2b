package payment

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
)

var testContainer TestContainer
var PaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = NewPackageIntegrationTestContainer()
	ctx := context.Background()
	var err error
	PaymentPartnerId, err = testContainer.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:    "chapa",
			Icon:    "etst",
			BaseUrl: "https://api.chapa.co",
			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		})
	if err != nil {
		panic("failed to create payment partner")
	}
}

func Test_Checkout(t *testing.T) {
	t.Run("paymentPartnerNotSupplied", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		in := &CheckoutRequest{
			TransactionRef: "1",
		}
		_, err := testContainer.PaymentService.Checkout(ctx, in)
		wantErr := ErrPaymentPartnerNotSupported
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("transactionRefNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			TransactionRef:   "1",
			PaymentPartnerId: 1,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		if err != nil {
			t.Fatalf("Failed to checkout err: %v", err)
		}
	})
}

func Test_Checkout_unhappyPath(t *testing.T) {

	t.Run("amountNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			PaymentPartnerId: PaymentPartnerId,
		}
		_, err := testContainer.PaymentService.Checkout(ctx, in)
		wantErr := ErrAmountNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("PaymentPartnerNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			TransactionRef: "1",
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		wantErr := ErrPaymentPartnerNotSupported
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Verify_Payment_happyPath(t *testing.T) {

}

func Test_Verify_Payment_unhappyPath(t *testing.T) {
	t.Run("TransactionRefNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		in_tx_ref := ""
		_, err := testContainer.PaymentService.Verify(ctx, 0, in_tx_ref)
		wantErr := ErrTransactionReferenceNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}
