package payment_verification

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
			BaseURL: "https://api.chapa.co",
			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		})
	if err != nil {
		panic("failed to create payment partner")
	}
}

func Test_Verify(t *testing.T) {
	t.Run("paymentPartnerIdMissing", func(t *testing.T) {
		ctx := context.Background()
		_, err := testContainer.PaymentVerificationService.Verify(ctx, 999, "90909090990909")
		wantErr := ErrPaymentPartnerNotSupported
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("transactionRefNotSupplied", func(t *testing.T) {
		ctx := context.Background()
		_, err := testContainer.PaymentVerificationService.Verify(ctx, PaymentPartnerId, "")
		wantErr := ErrTransactionReferenceNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}
