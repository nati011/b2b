package distributor

import (
	"context"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/distributor"
)

var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewDistributorService(
		db.NewMock(),
	)
}
func Test_Create_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		got, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if got.Id != id {
			t.Errorf("Expected id: %v Got id: %v", id, got.Id)
		}
	})

	t.Run("inactiveByDefault", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		got, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantIsAvailableStatus := false
		if got.IsAvailable != wantIsAvailableStatus {
			t.Errorf("Expected is available status: %v Got status: %v", wantIsAvailableStatus, got.IsAvailable)
		}
	})
}
