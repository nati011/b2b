package distributor

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/distributor"
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
	inactiveDistributorId, err := testContainer.DistributorService.Create(ctx, &distributor.CreateRequest{
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
	_, err = testContainer.ProductService.Create(ctx, &product.CreateRequest{
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

func Test_Deactivate_Products_Upon_Distributor_Deactivation(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	distributorId, err := testContainer.DistributorService.Create(ctx, &distributor.CreateRequest{
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
		t.Fatalf("failed to create distributor %v", err)
	}
	err = testContainer.DistributorService.Activate(ctx, distributorId)
	if err != nil {
		t.Fatalf("failed to activate distributor err: %v", err)
	}

	productId, err := testContainer.ProductService.Create(ctx, &product.CreateRequest{
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
		DistributorId: distributorId,
	})
	if err != nil {
		t.Fatalf("failed to create product %v", err)
	}

	err = testContainer.DistributorService.Dectivate(ctx, distributorId)
	if err != nil {
		t.Fatalf("failed to deactivate disributor %v", err)
	}

	product, err := testContainer.ProductService.Get(ctx, productId)
	if err != nil {
		t.Fatalf("failed to get product %v", err)
	}
	wantActiveStatus := false
	if product.IsActive != wantActiveStatus {
		t.Errorf("Expected active status: %v Got: %v", wantActiveStatus, product.IsActive)
	}
}
