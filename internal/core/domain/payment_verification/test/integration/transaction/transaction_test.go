package transaction

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/payment_verification"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var testContainer payment_verification.TestContainer
var db *sql.DB
var PaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	testContainer = payment_verification.NewPackageIntegrationTestContainer()
	ctx := context.Background()
	var err error
	PaymentPartnerId, err = testContainer.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:    "chapa",
			Icon:    "etst",
			BaseURL: "https://api.chapa.co",
			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		},
	)
	if err != nil {
		panic("failed to create payment partner")
	}
}

func Test_VerifyTransaction(t *testing.T) {
	ctx := context.Background()
	checkout_response, err := testContainer.CheckoutService.Checkout(ctx, &checkout.CheckoutRequest{
		PaymentPartnerId: PaymentPartnerId,
		OrderId:          1,
		Amount:           100,
	})
	if err != nil {
		t.Errorf("Failed to checkout err: %v", err)
	}

	resp, err := testContainer.PaymentVerificationService.Verify(ctx, PaymentPartnerId, checkout_response.TransactionRef)
	if err != nil {
		t.Fatalf("Failed to verify err: %v", err)
	}
	wantIsValidStatus := false
	if resp != wantIsValidStatus {
		t.Fatalf("Expected status: %v Got: %v", wantIsValidStatus, resp)
	}
}

func Test_StatusChangeUponSuccessfulVerification(t *testing.T) {
	ctx := context.Background()
	checkout_response, err := testContainer.CheckoutService.Checkout(ctx, &checkout.CheckoutRequest{
		PaymentPartnerId: PaymentPartnerId,
		OrderId:          1,
		Amount:           100,
	})
	if err != nil {
		t.Errorf("Failed to checkout err: %v", err)
	}

	resp, err := testContainer.PaymentVerificationService.Verify(ctx, PaymentPartnerId, checkout_response.TransactionRef)
	if err != nil {
		t.Fatalf("Failed to verify err: %v", err)
	}
	wantIsValidStatus := true
	if resp != wantIsValidStatus {
		t.Fatalf("Expected status: %v Got: %v", wantIsValidStatus, resp)
	}

	//check if transaction has been created and status has been set to COMPLETED
	resp_get, err := testContainer.TransactionService.GetByParam(ctx, &transaction.GetByParamRequest{
		TxRef: checkout_response.TransactionRef,
	})
	if err != nil {
		t.Fatalf("failed to get param err: %v", err)
	}
	wantStatus := transaction.COMPLETED_STATUS
	if resp_get.List[0].Status != wantStatus {
		t.Errorf("expected status: %v, got: %v", wantStatus, resp_get.List[0].Status)
	}
}
