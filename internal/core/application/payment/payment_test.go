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
	USERID = 0
	testContainer = NewPackageIntegrationTestContainer()
}

func Test_Checkout_happyPath(t *testing.T) {
	t.Run("init_checkout", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			User_Id:           USERID,
			Amount:            1,
			PaymentPartner_Id: PaymentPartnerId,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		if err != nil {
			t.Fatalf("Failed to checkout err: %v", err)
		}
	})
}

func Test_Checkout_unhappyPath(t *testing.T) {
	t.Run("userNotFound", func(t *testing.T) {
		ctx := context.Background()
		in := &CheckoutRequest{
			User_Id:           99,
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
			Amount:  -1,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		wantErr := ErrPaymentPartnerNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Verify_Payment(t *testing.T) {
}
