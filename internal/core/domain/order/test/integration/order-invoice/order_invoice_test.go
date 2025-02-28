package order_invoice

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/order"
)

var testContainer order.TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = order.NewPackageIntegrationTestContainer()
}

func Test_Create_Invoice_Upon_Order_Placement(t *testing.T) {
	ctx := context.Background()
	in := &order.PlaceRequest{
		RetailerId: 1,
		Items: []order.Item{
			{
				ProductId: 1,
				Quantity:  19},
		},
	}
	id, err := testContainer.OrderService.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	_, err = testContainer.InvoiceService.GetByParam(ctx, &invoice.GetByParamRequest{
		OrderId: id,
	})
	if err != nil {
		t.Fatalf("Failed to get invoice err: %v", err)
	}
}
