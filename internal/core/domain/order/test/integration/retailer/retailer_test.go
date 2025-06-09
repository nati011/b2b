package retailer

import (
	"context"
	"os"
	"testing"

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

func Test_validate_retailer_upon_registration(t *testing.T) {
	ctx := context.Background()

	in := &order.PlaceRequest{
		RetailerId: 1,
		Items: []order.Item{
			{
				ProductId: 1,
				Quantity:  1},
		},
		PaymentMethod: order.PAYMENT_METHOD_DIGITAL,
	}
	_, err := testContainer.OrderService.Place(ctx, in)
	wantErr := order.ErrRetailerIdNotFound
	if err != wantErr {
		t.Errorf("Expected err : %v Got: %v", wantErr, err)
	}
}
