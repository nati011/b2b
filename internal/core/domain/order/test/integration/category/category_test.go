package category

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
)

var testContainer order.TestContainer
var productId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = order.NewPackageIntegrationTestContainer()
	ctx := context.Background()
	productId, _ = testContainer.ProductService.Create(ctx, &product.CreateRequest{
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
}

func Test_validate_category(t *testing.T) {
	ctx := context.Background()
	in := &order.PlaceRequest{
		RetailerId: 1,
		Items: []order.Item{
			{
				ProductId: productId,
				Quantity:  19},
		},
	}
	id, err := testContainer.OrderService.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	//check
	resp, err := testContainer.OrderService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to fetch order err: err %v", err)
	}
	if resp.Id != id {
		t.Errorf("Expected Id: %v Got Id: %v", id, resp.Id)
	}
}
