package payment_verification

import (
	"context"
	"log"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

var testContainer TestContainer
var DigitalPaymentPartnerId int
var ManualPaymentPartnerId int
var retailerId int
var distributorId int
var productId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = NewPackageIntegrationTestContainer()
	ctx := context.Background()
	var err error
	retailerId, err = testContainer.RetailerService.Create(ctx, &retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "order_test",
		FirstName:   "test",
		LastName:    "test",
		Phone:       "+251949184879",
		Email:       "test@gmail.com",
	})
	if err != nil {
		panic("failed to create product")
	}
	distributorId, err = testContainer.DistributorService.Create(ctx, &distributor.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "username",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@gmail.com",
	})
	if err != nil {
		panic("failed to create distributor err: ")
	}
	productId, err = testContainer.ProductService.Create(ctx, &product.CreateRequest{
		Name:       "testProduct",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://picsum.photos/200",
			"https://picsum.photos/200",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
		DistributorId: distributorId,
	})
	if err != nil {
		panic("failed to create product")
	}
	testContainer.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     productId,
		Amount: 10000,
	})

	DigitalPaymentPartnerId, err = testContainer.PartnerService.Create(ctx,
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

	ManualPaymentPartnerId, err = testContainer.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:          "pay_delivery",
			Icon:          "test",
			BaseURL:       "https://localhost:8080",
			Secret:        "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
			PaymentMethod: payment_partner.PAYMENT_METHOD_MANUAL,
		})
	if err != nil {
		panic("failed to create payment partner")
	}
}

func teardown() {
	testContainer.TearDown()
}

func Test_Verify_UnhappyPath(t *testing.T) {
	t.Run("transactionRefNotSupplied", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := testContainer.PaymentVerificationService.Verify(ctx, "")
		wantErr := ErrTransactionReferenceNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Verify_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := &order.PlaceRequest{
		PaymentPartnerId: DigitalPaymentPartnerId,
		RetailerId:       retailerId,
		Items: []order.Item{
			{
				ProductId: productId,
				Quantity:  1},
		},
	}
	log.Printf("Partner id %v", DigitalPaymentPartnerId)
	place_resp, err := testContainer.OrderService.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	verifyResp, err := testContainer.PaymentVerificationService.Verify(ctx, place_resp.TxRef)
	if err != nil {
		t.Fatalf("Failed to verify err: %v", err)
	}
	wantVerifyRespStatus := false
	if verifyResp != wantVerifyRespStatus {
		t.Errorf("Expected verify resp: %v Got: %v", wantVerifyRespStatus, verifyResp)
	}
}

func Test_Confirm_unhappyPath(t *testing.T) {
	t.Run("transactionRefNotSupplied", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		err := testContainer.PaymentVerificationService.Confirm(ctx, "")
		wantErr := ErrTransactionReferenceNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("cannotConfirmDigitalPayments", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &order.PlaceRequest{
			PaymentPartnerId: DigitalPaymentPartnerId,
			RetailerId:       retailerId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  1},
			},
		}
		log.Printf("Partner id %v", DigitalPaymentPartnerId)
		place_resp, err := testContainer.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		err = testContainer.PaymentVerificationService.Confirm(ctx, place_resp.TxRef)
		wantErr := ErrCannotConfirmDigitalPayment
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_Confirm_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &order.PlaceRequest{
		PaymentPartnerId: ManualPaymentPartnerId,
		RetailerId:       retailerId,
		Items: []order.Item{
			{
				ProductId: productId,
				Quantity:  1},
		},
	}
	log.Printf("Partner id %v", DigitalPaymentPartnerId)
	place_resp, err := testContainer.OrderService.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	err = testContainer.PaymentVerificationService.Confirm(ctx, place_resp.TxRef)
	if err != nil {
		t.Fatalf("Failed to confirm err: %v", err)
	}
	verifyResp, err := testContainer.TransactionService.GetByParam(ctx, &transaction.GetByParamRequest{
		TxRef: place_resp.TxRef,
	})
	if err != nil {
		t.Fatalf("Failed to verify err: %v", err)
	}
	wantVerifyRespStatus := transaction.COMPLETED_STATUS
	if verifyResp.List[len(verifyResp.List)-1].Status != transaction.COMPLETED_STATUS {
		t.Errorf("Expected transaction status resp: %v Got: %v", wantVerifyRespStatus, verifyResp.List[len(verifyResp.List)-1].Status)
	}
}

func Test_Callback_happyPath(t *testing.T) {
}

func Test_Callback_unhappyPath(t *testing.T) {
}
