package transaction

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/checkout/test"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/transaction"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
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
			Icon:    "test",
			BaseURL: "https://api.chapa.co",
			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		})

	product_id, _ := testContainer.ProductService.Create(ctx, &product.CreateRequest{
		Name:       "testProduct",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"test",
			"test",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test":         "test",
			"another_test": "another_test",
		},
	})

	if err != nil {
		panic("failed to create product")
	}

	err = testContainer.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     product_id,
		Amount: 100,
	})

	if err != nil {
		panic("failed to add amount")
	}
	if err != nil {
		panic("failed to create payment partner")
	}
	OrderId, err = testContainer.OrderService.Place(ctx,
		&order.PlaceRequest{
			RetailerId: 1,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		},
	)

	if err != nil {
		panic("failed to create order partner")
	}
}

func Test_CreateTransactionUponPaymentInitAndSetStatusToPending(t *testing.T) {
	ctx := context.Background()

	//init transaction
	in := &checkout.CheckoutRequest{
		PaymentPartnerId: PaymentPartnerId,
		OrderId:          OrderId,
		Amount:           100,
	}

	checkout_response, err := testContainer.CheckoutService.Checkout(ctx, in)
	if err != nil {
		t.Errorf("Failed to checkout err: %v", err)
	}

	//check if transaction has been created
	resp, err := testContainer.TransactionService.GetByParam(ctx, &transaction.GetByParamRequest{
		TxRef: checkout_response.TransactionRef,
	})
	if err != nil {
		t.Fatalf("failed to get param err: %v", err)
	}
	wantStatus := transaction.PENDING_STATUS
	if resp.List[0].Status != wantStatus {
		t.Errorf("expected status: %v, got: %v", resp.List[0].Status, wantStatus)
	}
}
