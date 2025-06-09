package paymentpartner

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/checkout/test"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
)

var testContainer test.TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = test.NewPackageIntegrationTestContainer()
}

func Test_DisallowCheckOutIfPaymentProviderIsInactive(t *testing.T) {
	ctx := context.Background()
	var err error
	PaymentPartnerId, err := testContainer.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:    "chapa",
			Icon:    "test",
			BaseURL: "https://api.chapa.co",
			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		})

	if err != nil {
		panic("failed to create order partner")
	}

	//activate
	err = testContainer.PartnerService.Deactivate(ctx, PaymentPartnerId)
	if err != nil {
		t.Fatalf("Failed to activate payment partner: %v", err)
	}
	_, err = testContainer.CheckoutService.Checkout(ctx, &checkout.CheckoutRequest{
		PaymentPartnerId: PaymentPartnerId,
		OrderId:          1,
		Amount:           100,
		PaymentMethod:    checkout.PAYMENT_METHOD_DIGITAL,
	})
	wantErr := checkout.ErrPaymentPartnerNotSupported
	if err != wantErr {
		t.Errorf("Expected err: %v Got: %v", wantErr, err)
	}
}

func Test_AllowCheckOutIfPaymentProviderIsActive(t *testing.T) {
	ctx := context.Background()
	PaymentPartnerId, err := testContainer.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:    "chapa",
			Icon:    "test",
			BaseURL: "https://api.chapa.co",
			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		})
	if err != nil {
		t.Fatalf("Failed to create payment partner: %v", err)
	}

	_, err = testContainer.CheckoutService.Checkout(ctx, &checkout.CheckoutRequest{
		PaymentPartnerId: PaymentPartnerId,
		OrderId:          1,
		Amount:           100,
		PaymentMethod:    checkout.PAYMENT_METHOD_DIGITAL,
	})
	if err != nil {
		t.Fatalf("Failed to checkout: %v", err)
	}
}
