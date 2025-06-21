package order_expiry

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/config"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	"github.com/go-co-op/gocron/v2"
)

var container TestContainer
var db *sql.DB
var retailer_id int
var distributor_id int
var product_id int
var DigitalPaymentPartnerId int
var ManualPaymentPartnerId int

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	db = db_test_container.Setup()
	container = NewIntegrationTestContainer(db)
	var err error
	retailer_id, err = container.RetailerService.Create(ctx, &retailer.CreateRequest{
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

	distributor_id, err = container.DistributorService.Create(ctx, &distributor.CreateRequest{
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
	product_id, err = container.ProductService.Create(ctx, &product.CreateRequest{
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
		DistributorId: distributor_id,
	})
	if err != nil {
		panic("failed to create product")
	}
	container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     product_id,
		Amount: 10000,
	})
	DigitalPaymentPartnerId, err = container.PartnerService.Create(ctx,
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

	ManualPaymentPartnerId, err = container.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:          "chapa",
			Icon:          "etst",
			BaseURL:       "https://api.chapa.co",
			Secret:        "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
			PaymentMethod: payment_partner.PAYMENT_METHOD_MANUAL,
		})
	if err != nil {
		panic("failed to create payment partner")
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	InitOrderExpiry(
		s,
		&container.OrderService,
		&container.ConfigService,
		&container.PaymentService,
		&container.PartnerService)

	s.Start()
}

func teardown() {
	db_test_container.Teardown(db)
}

func Test_CancelExpiredOrders(t *testing.T) {
	t.Run("cancel_expired_order_digitalPaymentMethod", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		// create an order
		in := &order.PlaceRequest{
			PaymentPartnerId: DigitalPaymentPartnerId,
			RetailerId:       retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		log.Printf("Partner id %v", DigitalPaymentPartnerId)
		order_resp, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		resp, err := container.OrderService.Get(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		if resp.Id != order_resp.Id {
			t.Errorf("Expected Id: %v Got Id: %v", order_resp.Id, resp.Id)
		}

		// set order expiry config to 1 minute
		err = container.ConfigService.SetOrderExpiryConfig(ctx, &config.SetOrderExpiryRequest{
			ExpiryDurationInMinues: 1,
		})
		if err != nil {
			t.Fatalf("Failed to set order expiry date: err %v", err)
		}
		// wait two minutes
		time.Sleep(2 * time.Minute)

		// confirm order status is canceled
		order_after, err := container.OrderService.Get(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to get order expiry date: err %v", err)
		}
		wantOrderStatus := order.CANCELED_STATUS
		if order_after.Status != wantOrderStatus {
			t.Errorf("Expected order status: %v Got: %v", wantOrderStatus, order_after.Status)
		}
	})

	t.Run("do_not_cancel_not_expired_order_digitalPaymentMethod", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		// create an order
		in := &order.PlaceRequest{
			PaymentPartnerId: DigitalPaymentPartnerId,
			RetailerId:       retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		log.Printf("Partner id %v", DigitalPaymentPartnerId)
		order_resp, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		resp, err := container.OrderService.Get(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		if resp.Id != order_resp.Id {
			t.Errorf("Expected Id: %v Got Id: %v", order_resp.Id, resp.Id)
		}

		// set order expiry config to 1 minute
		err = container.ConfigService.SetOrderExpiryConfig(ctx, &config.SetOrderExpiryRequest{
			ExpiryDurationInMinues: 1,
		})
		if err != nil {
			t.Fatalf("Failed to set order expiry date: err %v", err)
		}

		// confirm order status is canceled
		order_after, err := container.OrderService.Get(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to get order expiry date: err %v", err)
		}
		wantOrderStatus := order.CANCELED_STATUS
		if order_after.Status != wantOrderStatus {
			t.Errorf("Expected order status: %v Got: %v", wantOrderStatus, order_after.Status)
		}
	})

	t.Run("cancel_exipired_order_manualPaymentMethod", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		// create an order
		in := &order.PlaceRequest{
			PaymentPartnerId: DigitalPaymentPartnerId,
			RetailerId:       retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		log.Printf("Partner id %v", ManualPaymentPartnerId)
		order_resp, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		resp, err := container.OrderService.Get(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		if resp.Id != order_resp.Id {
			t.Errorf("Expected Id: %v Got Id: %v", order_resp.Id, resp.Id)
		}

		// set order expiry config to 1 minute
		err = container.ConfigService.SetOrderExpiryConfig(ctx, &config.SetOrderExpiryRequest{
			ExpiryDurationInMinues: 1,
		})
		if err != nil {
			t.Fatalf("Failed to set order expiry date: err %v", err)
		}
		// wait two minutes
		time.Sleep(2 * time.Minute)

		// confirm order status is canceled
		order_after, err := container.OrderService.Get(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to get order expiry date: err %v", err)
		}
		wantOrderStatus := order.CANCELED_STATUS
		if order_after.Status != wantOrderStatus {
			t.Errorf("Expected order status: %v Got: %v", wantOrderStatus, order.CANCELED_STATUS)
		}
	})

	t.Run("donot_cancel_not_exipired_order_manualPaymentMethod", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		// create an order
		in := &order.PlaceRequest{
			PaymentPartnerId: DigitalPaymentPartnerId,
			RetailerId:       retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		log.Printf("Partner id %v", ManualPaymentPartnerId)
		order_resp, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		resp, err := container.OrderService.Get(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		if resp.Id != order_resp.Id {
			t.Errorf("Expected Id: %v Got Id: %v", order_resp.Id, resp.Id)
		}

		// set order expiry config to 1 minute
		err = container.ConfigService.SetOrderExpiryConfig(ctx, &config.SetOrderExpiryRequest{
			ExpiryDurationInMinues: 1,
		})
		if err != nil {
			t.Fatalf("Failed to set order expiry date: err %v", err)
		}

		// confirm order status is canceled
		order_after, err := container.OrderService.Get(ctx, order_resp.Id)
		if err != nil {
			t.Fatalf("Failed to get order expiry date: err %v", err)
		}
		wantOrderStatus := order.CANCELED_STATUS
		if order_after.Status != wantOrderStatus {
			t.Errorf("Expected order status: %v Got: %v", wantOrderStatus, order.CANCELED_STATUS)
		}
	})
}
