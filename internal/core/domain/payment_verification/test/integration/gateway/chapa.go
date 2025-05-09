package gateway

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/payment_verification/test"
)

var testContainer test.TestContainer
var PaymentPartnerId int
var OrderId int

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
	OrderId, err = testContainer.OrderService.Place(ctx,
		&order.PlaceRequest{
			RetailerId: 1,
			Items: []order.Item{
				{
					ProductId: 1,
					Quantity:  1},
			},
		},
	)
	if err != nil {
		panic("failed to create payment partner")
	}
}

func Test_Verify_Payment(t *testing.T) {

	t.Run("validTransaction", func(t *testing.T) {
		ctx := context.Background()
		//init transaction
		in := &checkout.CheckoutRequest{
			PaymentPartnerId: PaymentPartnerId,
			OrderId:          OrderId,
		}

		checkout_resp, err := testContainer.CheckoutService.Checkout(ctx, in)
		if err != nil {
			t.Errorf("Failed to checkout err: %v", err)
		}

		resp, err := testContainer.PaymementVerificationService.Verify(ctx, PaymentPartnerId, checkout_resp.TransactionRef)
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

		resp, err := testContainer.PaymementVerificationService.Verify(ctx, PaymentPartnerId, "invalidTransaction")
		if err != nil {
			t.Fatalf("Failed to verify err: %v", err)
		}
		wantIsValidStatus := false
		if resp != wantIsValidStatus {
			t.Errorf("Expected status: %v Got: %v", wantIsValidStatus, resp)
		}
	})

}
