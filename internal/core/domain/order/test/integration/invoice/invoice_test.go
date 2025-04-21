package order_invoice

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

var container order.TestContainer
var retailer_id int
var product_id int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = order.NewPackageIntegrationTestContainer()
	ctx := context.Background()
	retailer_id, _ = container.RetailerService.Create(ctx, &retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "order_retilertest",
		FirstName:   "test",
		LastName:    "test",

		Email: "test@gmail.com",
	})
	product_id, _ = container.ProductService.Create(ctx, &product.CreateRequest{
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
		Id:     product_id,
		Amount: 100,
	})
}

func Test_Create_Invoice_Upon_Order_Placement(t *testing.T) {
	ctx := context.Background()
	in := &order.PlaceRequest{
		RetailerId: retailer_id,
		Items: []order.Item{
			{
				ProductId: product_id,
				Quantity:  19},
		},
	}
	id, err := container.OrderService.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	_, err = container.InvoiceService.GetByParam(ctx, &invoice.GetByParamRequest{
		OrderId: id,
	})
	if err != nil {
		t.Fatalf("Failed to get invoice err: %v", err)
	}
}
