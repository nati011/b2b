package paymentpartner

import (
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/payment_verification/test"
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

// func Test_DisallowVerificationIfPaymentProviderIsInactive(t *testing.T) {
// 	ctx := context.Background()
// 	var err error
// 	PaymentPartnerId, err := testContainer.PartnerService.Create(ctx,
// 		&payment_partner.CreateRequest{
// 			Name:    "chapa",
// 			Icon:    "test",
// 			BaseURL: "https://api.chapa.co",
// 			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
// 		})

// 	if err != nil {
// 		panic("failed to create order partner")
// 	}
// 	err = testContainer.PartnerService.Deactivate(ctx, PaymentPartnerId)
// 	if err != nil {
// 		panic("failed to deactivate payment partner")
// 	}
// 	_, err = testContainer.PaymementVerificationService.Verify(ctx, "txRef")
// 	wantErr := payment_verification.ErrPaymentPartnerNotSupported
// 	if err != wantErr {
// 		t.Errorf("Expected err: %v Got: %v", wantErr, err)
// 	}
// }

// func Test_AllowVerificationIfPaymentProviderIsActive(t *testing.T) {
// 	ctx := context.Background()
// 	PaymentPartnerId, err := testContainer.PartnerService.Create(ctx,
// 		&payment_partner.CreateRequest{
// 			Name:    "chapa",
// 			Icon:    "test",
// 			BaseURL: "https://api.chapa.co",
// 			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
// 		})
// 	if err != nil {
// 		t.Fatalf("Failed to create payment partner: %v", err)
// 	}
// 	_, err = testContainer.PaymementVerificationService.Verify(ctx,  "txRef")
// 	if err != nil {
// 		t.Fatalf("Failed to checkout: %v", err)
// 	}
// }
