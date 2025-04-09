package product

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

var container order.TestContainer
var retailer_id int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	container = order.NewPackageIntegrationTestContainer()
	retailer_id, _ = container.RetailerService.Create(ctx, &retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",

		FirstName: "test",
		LastName:  "test",

		Email: "test@gmail.com",
	})
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
		_, err := container.OrderService.Place(ctx, in)
		wantErr := order.ErrItemMemberProductNotFound
		if err != wantErr {
			t.Errorf("Expected err : %v Got: %v", wantErr, err)
		}
	})
}
