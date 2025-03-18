package product

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

func Test_validate_Item_exists_upon_order_creation(t *testing.T) {
	t.Run("unhappy_Path", func(t *testing.T) {
		ctx := context.Background()
		in := &order.PlaceRequest{
			RetailerId: 1,
			Items: []order.Item{
				{
					ProductId: 1,
					Quantity:  19},
			},
		}
		_, err := testContainer.OrderService.Place(ctx, in)
		wantErr := order.ErrItemMemberProductNotFound
		if err != wantErr {
			t.Errorf("Expected err : %v Got: %v", wantErr, err)
		}
	})
}
