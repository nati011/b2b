package checkout

import (
	"context"
	"log"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
)

var testContainer TestContainer
var PaymentPartnerId int
var OrderId int

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
			Name:          "chapa",
			Icon:          "etst",
			BaseURL:       "https://api.chapa.co",
			Secret:        "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
			PaymentMethod: payment_partner.PAYMENT_METHOD_DIGITAL,
		})
	if err != nil {
		panic("failed to create payment partner")
	}
}

func Test_Checkout(t *testing.T) {
	t.Run("paymentPartnerNotSupplied", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		setup()
		ctx := context.Background()

		_, err := testContainer.CheckoutService.Checkout(ctx, &CheckoutRequest{
			OrderId: 1,
			Amount:  100,
		})
		wantErr := ErrPaymentPartnerNotSupported
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("checkoutHappyPath", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(testContainer.TearDown)
		setup()
		in := &CheckoutRequest{
			OrderId:          1,
			Amount:           400,
			PaymentPartnerId: PaymentPartnerId,
		}

		_, err := testContainer.CheckoutService.Checkout(ctx, in)
		if err != nil {
			t.Errorf("Failed to checkout %v", err)
		}

	})

	t.Run("amountNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(testContainer.TearDown)
		setup()
		in := &CheckoutRequest{
			PaymentPartnerId: PaymentPartnerId,
		}
		_, err := testContainer.CheckoutService.Checkout(ctx, in)
		log.Printf("Got Error: %v", err)
		wantErr := ErrAmountNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_SubscriptionPayment(t *testing.T) {
	t.Run("paymentPartnerNotSupplied", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		setup()
		ctx := context.Background()

		_, err := testContainer.CheckoutService.Checkout(ctx, &CheckoutRequest{
			OrderId: 1,
			Amount:  100,
		})
		wantErr := ErrPaymentPartnerNotSupported
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("checkoutHappyPath", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(testContainer.TearDown)
		setup()
		in := &CheckoutRequest{
			OrderId:          1,
			Amount:           400,
			PaymentPartnerId: PaymentPartnerId,
		}

		_, err := testContainer.CheckoutService.Checkout(ctx, in)
		if err != nil {
			t.Errorf("Failed to checkout %v", err)
		}

	})

	t.Run("amountNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(testContainer.TearDown)
		setup()
		in := &CheckoutRequest{
			PaymentPartnerId: PaymentPartnerId,
		}
		_, err := testContainer.CheckoutService.Checkout(ctx, in)
		log.Printf("Got Error: %v", err)
		wantErr := ErrAmountNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}
