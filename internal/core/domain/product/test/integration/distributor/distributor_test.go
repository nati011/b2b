package distributor

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/product"
)

var testContainer product.TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = product.NewPackageIntegrationTestContainer()
}

func teardown() {
	testContainer.Teardown()
}

func Test_Validate_Distributor(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(teardown)
	distributor_id := 99
	_, err := testContainer.ProductService.Create(ctx, &product.CreateRequest{
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
		DistributorId: distributor_id,
	})
	wantErr := product.ErrDistributorNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v Got: %v", wantErr, err)
	}
}
