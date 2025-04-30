package gateway

import (
	"context"
	"os"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/application/payment"
	"b2b.nati011.github.com/internal/core/application/payment/test"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
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
			Icon:    "etst",
			BaseURL: "https://api.chapa.co",
			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		})
	if err != nil {
		panic("failed to create payment partner")
	}
}

func Test_Checkout(t *testing.T) {

	t.Run("init_checkout", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		currentTimestamp := time.Now()
		generatedTxRef := currentTimestamp.Format("2006_01_02_15_04_05")
		in := &payment.CheckoutRequest{
			PaymentPartnerId: PaymentPartnerId,
			TransactionRef:   generatedTxRef,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		if err != nil {
			t.Errorf("Failed to checkout err: %v", err)
		}
	})
}

func Test_Verify_Payment(t *testing.T) {

	t.Run("validTransaction", func(t *testing.T) {
		ctx := context.Background()
		//init transaction

		currentTimestamp := time.Now()
		generatedTxRef := currentTimestamp.Format("2006_01_02_15_04_05")
		in := &payment.CheckoutRequest{
			PaymentPartnerId: PaymentPartnerId,
			TransactionRef:   generatedTxRef,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		if err != nil {
			t.Errorf("Failed to checkout err: %v", err)
		}

		resp, err := testContainer.PaymentService.Verify(ctx, PaymentPartnerId, generatedTxRef)
		if err != nil {
			t.Errorf("Failed to verify err: %v", err)
		}
		wantIsValidStatus := true
		if resp != wantIsValidStatus {
			t.Errorf("Expected status: %v Got: %v", wantIsValidStatus, resp)
		}
	})

	t.Run("invalidTransaction", func(t *testing.T) {
		ctx := context.Background()
		//init transaction
		currentTimestamp := time.Now()
		generatedTxRef := currentTimestamp.Format("2006_01_02_15_04_05")
		in := &payment.CheckoutRequest{
			PaymentPartnerId: PaymentPartnerId,
			TransactionRef:   generatedTxRef,
		}

		_, err := testContainer.PaymentService.Checkout(ctx, in)
		if err != nil {
			t.Errorf("Failed to checkout err: %v", err)
		}

		resp, err := testContainer.PaymentService.Verify(ctx, PaymentPartnerId, generatedTxRef)
		if err != nil {
			t.Fatalf("Failed to verify err: %v", err)
		}
		wantIsValidStatus := false
		if resp != wantIsValidStatus {
			t.Errorf("Expected status: %v Got: %v", wantIsValidStatus, resp)
		}
	})

}
