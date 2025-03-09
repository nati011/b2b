package payment

import (
	"context"
	"os"
	"testing"
)

var testContainer TestContainer
var USERID int
var PaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	USERID = 1
	testContainer = NewPackageIntegrationTestContainer()
}

func Test_Checkout_happyPath(t *testing.T) {
	t.Run("init_checkout", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			User_Id:           USERID,
			Amount:            1,
			PaymentPartner_Id: 1,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		if err != nil {
			t.Fatalf("Failed to checkout err: %v", err)
		}
	})
}

func Test_Checkout_unhappyPath(t *testing.T) {
	t.Run("userNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			Amount:            1,
			PaymentPartner_Id: PaymentPartnerId,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		wantErr := ErrUserNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("amountNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			User_Id:           USERID,
			PaymentPartner_Id: PaymentPartnerId,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		wantErr := ErrAmountNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("amountGreaterThanZero", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			User_Id:           USERID,
			PaymentPartner_Id: PaymentPartnerId,
			Amount:            -1,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		wantErr := ErrAmountLessThanZero
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("PaymentPartnerNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			User_Id: USERID,
			Amount:  1,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		wantErr := ErrPaymentPartnerNotSupplied
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
		_, err := testContainer.PaymentService.Verify(ctx, in_tx_ref)
		wantErr := ErrTransactionReferenceNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}
