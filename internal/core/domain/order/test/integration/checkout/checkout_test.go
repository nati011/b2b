package checkout

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/order"
	test_container "b2b.nati011.github.com/internal/core/domain/order/test"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var container test_container.TestContainer
var db *sql.DB
var retailerId int
var productId int
var manualPaymentPartnerId int
var digitalPaymentPartnerId int
var distributorId int

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	db = db_test_container.Setup()
	container = test_container.NewDBIntegrationTestContainer(
		db,
	)
	var err error
	retailerId, _ = container.RetailerService.Create(ctx, &retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "ordeR_retailer",
		FirstName:   "test",
		LastName:    "test",
		Phone:       "+251949184879",
		Email:       "test@gmail.com",
	})
	distributorId, err = container.DistributorService.Create(ctx, &distributor.CreateRequest{
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
		panic("failed to create distributor")
	}
	err = container.DistributorService.Activate(ctx, distributorId)
	if err != nil {
		panic("failed to activate distributor")
	}
	productId, _ = container.ProductService.Create(ctx, &product.CreateRequest{
		Name:       "testProduct",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
		DistributorId: distributorId,
	})
	container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     productId,
		Amount: 100,
	})

	manualPaymentPartnerId, _ = container.PartnerService.Create(ctx, &payment_partner.CreateRequest{
		Name:          "chapa",
		Icon:          "etst",
		BaseURL:       "https://api.chapa.co",
		Secret:        "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		PaymentMethod: payment_partner.PAYMENT_METHOD_MANUAL,
	})

	digitalPaymentPartnerId, _ = container.PartnerService.Create(ctx, &payment_partner.CreateRequest{
		Name:          "chapa",
		Icon:          "etst",
		BaseURL:       "https://api.chapa.co",
		Secret:        "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		PaymentMethod: payment_partner.PAYMENT_METHOD_MANUAL,
	})
}

func teardown() {
	container.Teardown(db)
	db_test_container.Teardown(db)
}

func Test_Checkout_happyPath(t *testing.T) {
	t.Run("generate_checkout_url_upon_order_placement_digital", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := &order.PlaceRequest{
			PaymentPartnerId: digitalPaymentPartnerId,
			RetailerId:       retailerId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  1},
			},
		}
		order_resp, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		if strings.Split(order_resp.CheckoutUrl, " ") == nil {
			t.Errorf("Expected checkoutUrl different from nil got: %v", order_resp.CheckoutUrl)
		}
	})

	t.Run("checkout_upon_order_placement_manual", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := &order.PlaceRequest{
			PaymentPartnerId: manualPaymentPartnerId,
			RetailerId:       retailerId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  1},
			},
		}
		_, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
	})

	t.Run("generate_checkout_url_upon_init_payment", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := &order.PlaceRequest{
			PaymentPartnerId: digitalPaymentPartnerId,
			RetailerId:       retailerId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  1},
			},
		}
		order_resp, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		init_resp, err := container.OrderService.InitPayment(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to init payment err: %v", err)
		}
		if strings.Split(init_resp.CheckoutUrl, " ") == nil {
			t.Errorf("Expected checkoutUrl different from nil got: %v", init_resp.CheckoutUrl)
		}
	})
}

func Test_Checkout_unhappyPath(t *testing.T) {
	t.Run("cannot_settle_payment_if_order_canceled", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := &order.PlaceRequest{
			PaymentPartnerId: digitalPaymentPartnerId,
			RetailerId:       retailerId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  1},
			},
		}
		order_resp, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		err = container.OrderService.Cancel(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to cancel order err: %v", err)
		}
		_, err = container.OrderService.InitPayment(ctx, order_resp.Id)
		wantErr := order.ErrCannotInitPaymentForCanceledOrder
		if err != wantErr {
			t.Errorf("Expeceted err: %v Want err: %v", err, wantErr)
		}

	})
}
