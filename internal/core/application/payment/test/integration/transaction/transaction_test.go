package transaction

import (
	"context"
	"os"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/application/payment"
	"b2b.nati011.github.com/internal/core/application/payment/test"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
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

func Test_CreateTransactionUponPaymentInitAndSetStatusToPending(t *testing.T) {
	ctx := context.Background()

	//init transaction
	currentTimestamp := time.Now()
	generatedTxRef := currentTimestamp.Format("2006_01_02_15_04_05")
	in := &payment.CheckoutRequest{
		Amount:           1,
		PaymentPartnerId: PaymentPartnerId,
		TransactionRef:   generatedTxRef,
	}

	_, err := testContainer.PaymentService.Checkout(ctx, in)
	if err != nil {
		t.Errorf("Failed to checkout err: %v", err)
	}

	//check if transaction has been created
	resp, err := testContainer.TransactionService.GetByParam(ctx, &transaction.GetByParamRequest{
		TxRef: in.TransactionRef,
	})
	if err != nil {
		t.Fatalf("failed to get param err: %v", err)
	}
	wantStatus := transaction.PENDING_STATUS
	if resp.List[0].Status != wantStatus {
		t.Errorf("expected status: %v, got: %v", resp.List[0].Status, wantStatus)
	}
}

func Test_CreateTransactionUponPaymentVerification(t *testing.T) {
	ctx := context.Background()
	//init transaction

	currentTimestamp := time.Now()
	generatedTxRef := currentTimestamp.Format("2006_01_02_15_04_05")
	in := &payment.CheckoutRequest{
		Amount:           1,
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
	wantIsValidStatus := true
	if resp != wantIsValidStatus {
		t.Fatalf("Expected status: %v Got: %v", wantIsValidStatus, resp)
	}

	//check if transaction has been created and status has been set to uploaded
	resp_get, err := testContainer.TransactionService.GetByParam(ctx, &transaction.GetByParamRequest{
		TxRef: in.TransactionRef,
	})
	if err != nil {
		t.Fatalf("failed to get param err: %v", err)
	}
	wantStatus := transaction.COMPLETED_STATUS
	if resp_get.List[0].Status != wantStatus {
		t.Errorf("expected status: %v, got: %v", resp_get.List[0].Status, wantStatus)
	}
}
