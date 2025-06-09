package order_invoice

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

var container order.TestContainer
var retailer_id int
var product_id int
var digitalPaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = order.NewPackageIntegrationTestContainer()
	ctx := context.Background()
	var err error
	retailer_id, err = container.RetailerService.Create(ctx, &retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "order_retilertest",
		FirstName:   "test",
		LastName:    "test",
		Phone:       "+251949184879",
		Email:       "test@gmail.com",
	})
	if err != nil {
		panic(err)
	}
	product_id, err = container.ProductService.Create(ctx, &product.CreateRequest{
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
	})
	if err != nil {
		panic(err)
	}

	err = container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     product_id,
		Amount: 100,
	})
	if err != nil {
		panic(err)
	}

	digitalPaymentPartnerId, err = container.PartnerService.Create(ctx, &payment_partner.CreateRequest{
		Name:          "chapa",
		Icon:          "etst",
		BaseURL:       "https://api.chapa.co",
		Secret:        "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		PaymentMethod: payment_partner.PAYMENT_METHOD_DIGITAL,
	})
	if err != nil {
		panic(err)
	}
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
		PaymentPartnerId: digitalPaymentPartnerId,
	}
	order_resp, err := container.OrderService.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	_, err = container.InvoiceService.GetByParam(ctx, &invoice.GetByParamRequest{
		OrderId: order_resp.Id,
	})
	if err != nil {
		t.Fatalf("Failed to get invoice err: %v", err)
	}
}
