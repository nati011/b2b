package distributor

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/product"
)

var testContainer product.TestContainer
var distributorId int
var inactiveDistributorId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	testContainer = product.NewPackageIntegrationTestContainer()
	var err error
	distributorId, err = testContainer.DistributorService.Create(ctx, &distributor.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "username",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@gmail.com",
	})
	if err != nil {
		panic("failed to create distributor err: ")
	}
	err = testContainer.DistributorService.Activate(ctx, distributorId)
	if err != nil {
		panic("failed to ctivate distributor")
	}
	inactiveDistributorId, err = testContainer.DistributorService.Create(ctx, &distributor.CreateRequest{
		Tin:         "1111111112",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "username2",
		FirstName:   "test2",
		LastName:    "test2",
		Email:       "test1@gmail.com",
	})
	if err != nil {
		panic("failed to create distributor")
	}

}

func teardown() {
	testContainer.Teardown()
}

func Test_Validate_Distributor(t *testing.T) {
	t.Run("invalidDistributorId", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(teardown)
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
			DistributorId: 101,
		})
		wantErr := product.ErrDistributorNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("distributorIdMandatory", func(t *testing.T) {
		ctx := context.Background()
		t.Cleanup(teardown)
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
		})
		wantErr := product.ErrDistributorIdMandatory
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_inactiveDistributorCannotCreateProduct(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(teardown)
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
		DistributorId: inactiveDistributorId,
	})
	wantErr := product.ErrDistributorInactive
	if err != wantErr {
		t.Errorf("Expected err: %v Got: %v", wantErr, err)
	}
}
