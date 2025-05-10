package product

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

var container order.TestContainer
var retailerId int
var productId int
var paymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func teardown() {
	container.Teardown()
}

func setup() {
	ctx := context.Background()
	container = order.NewPackageIntegrationTestContainer()
	retailerId, _ = container.RetailerService.Create(ctx, &retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "product_retailer",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@gmail.com",
	})
	productId, _ = container.ProductService.Create(ctx, &product.CreateRequest{
		Name:       "testProduct",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"test",
			"test",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
	})
	container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     productId,
		Amount: 20,
	})
	paymentPartnerId, _ = container.PartnerService.Create(ctx, &payment_partner.CreateRequest{
		Name:    "chapa",
		Icon:    "etst",
		BaseURL: "https://api.chapa.co",
		Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
	})
}

func Test_validate_Item_exists_upon_order_creation(t *testing.T) {
	t.Run("unhappy_Path", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &order.PlaceRequest{
			RetailerId: retailerId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  100},
			},
			PaymentPartnerId: paymentPartnerId,
		}
		_, err := container.OrderService.Place(ctx, in)
		wantErr := order.ErrItemMemberProductQuantityNotFound
		if err != wantErr {
			t.Errorf("Expected err : %v Got: %v", wantErr, err)
		}
	})
}

func Test_Reserve_Stock_Upon_order_creation(t *testing.T) {
	t.Run("happyPath", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		orderQty := 1
		in := &order.PlaceRequest{
			RetailerId: productId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  orderQty},
			},
			PaymentPartnerId: paymentPartnerId,
		}
		product, err := container.ProductService.Get(ctx, productId)
		if err != nil {
			t.Fatalf("Failed to fetch product err%v", err)
		}
		wantAvailableStock := product.AvailableStock - orderQty

		_, err = container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err%v", err)
		}

		product, err = container.ProductService.Get(ctx, productId)
		if err != nil {
			t.Fatalf("Failed to fetch product err%v", err)
		}
		if product.AvailableStock != wantAvailableStock {
			t.Errorf("Expected stock: %v Got stock %v", wantAvailableStock, product.Stock)
		}
	})
}

func Test_Free_Reserved_Stock_Upon_order_status_change(t *testing.T) {
	t.Run("delivered", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		orderQty := 1
		in := &order.PlaceRequest{
			RetailerId: productId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  orderQty},
			},
			PaymentPartnerId: paymentPartnerId,
		}
		product, err := container.ProductService.Get(ctx, productId)
		if err != nil {
			t.Fatalf("Failed to fetch product err%v", err)
		}
		wantAvailableStock := product.AvailableStock

		order_resp, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err%v", err)
		}

		container.OrderService.UpdateDeliveryStatus(ctx, &order.UpdateRequest{
			Id:             order_resp.Id,
			DeliveryStatus: order.DELIVERY_COMPLETED_STATUS,
		})

		product, err = container.ProductService.Get(ctx, productId)
		if err != nil {
			t.Fatalf("Failed to fetch product err%v", err)
		}
		if product.AvailableStock != wantAvailableStock {
			t.Errorf("Expected stock: %v Got stock %v", wantAvailableStock, product.AvailableStock)
		}
	})

	t.Run("canceled", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		orderQty := 1
		in := &order.PlaceRequest{
			RetailerId: productId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  orderQty},
			},
			PaymentPartnerId: paymentPartnerId,
		}
		product, err := container.ProductService.Get(ctx, productId)
		if err != nil {
			t.Fatalf("Failed to fetch product err%v", err)
		}
		wantAvailableStock := product.AvailableStock

		order_resp, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err%v", err)
		}

		container.OrderService.UpdateStatus(ctx, &order.UpdateRequest{
			Id:     order_resp.Id,
			Status: order.CANCELED_STATUS,
		})

		product, err = container.ProductService.Get(ctx, productId)
		if err != nil {
			t.Fatalf("Failed to fetch product err%v", err)
		}
		if product.AvailableStock != wantAvailableStock {
			t.Errorf("Expected stock: %v Got stock %v", wantAvailableStock, product.Stock)
		}
	})
}

func Test_Check_Stock_Availability_before_order_creation(t *testing.T) {
	t.Run("happyPath", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &order.PlaceRequest{
			RetailerId: productId,
			Items: []order.Item{
				{
					ProductId: productId,
					Quantity:  1000},
			},
			PaymentPartnerId: paymentPartnerId,
		}
		_, err := container.OrderService.Place(ctx, in)
		wantErr := order.ErrItemMemberProductQuantityNotFound
		if err != wantErr {
			t.Errorf("Expected err : %v Got: %v", wantErr, err)
		}
	})

	t.Run("unhappyPath", func(t *testing.T) {
		t.Cleanup(teardown)
	})
}
