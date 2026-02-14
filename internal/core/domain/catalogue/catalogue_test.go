package catalogue

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/product"
)

var container TestContainer
var category_id int
var distributorId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = NewPackageTestContainer()
	ctx := context.Background()
	category_id, _ = container.CategoryService.Create(ctx, &category.CreateRequest{
		Name: "test",
	})
	var err error
	distributorId, err = container.DistributorService.Create(ctx, &distributor.CreateRequest{
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
		panic("failed to create distributor")
	}
	err = container.DistributorService.Activate(ctx, distributorId)
	if err != nil {
		panic("failed to activate distributor")
	}
}

func tearDown() {
	container.Teardown()
}

func Test_GetAll(t *testing.T) {
	t.Run("product_list", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &product.CreateRequest{
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
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Errorf("Failed to create product err: %v", err)
		}

		got, err := container.CatalogueService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to get all ")
		}
		wantLen := 1
		if len(got.List) != wantLen {
			t.Errorf("Expected len: %v Got len:%v", wantLen, len(got.List))
		}
	})

	t.Run("configurable_product_list", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		productId, err := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test":         "test",
				"another_test": "another_test",
			},
			DistributorId: distributorId,
		})

		if err != nil {
			t.Fatalf("fauked to create product err: %v", err)
		}

		_, err = container.ConfigurableProductService.Create(ctx, &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				productId,
			},
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
		})
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		got, err := container.CatalogueService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to get all ")
		}
		wantLen := 1
		if len(got.List) != wantLen {
			t.Errorf("Expected len: %v Got len:%v", wantLen, len(got.List))
		}
	})
}

func Test_Search(t *testing.T) {
	t.Run("product", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &product.CreateRequest{
			Name:       "searchTest",
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
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}
		got, err := container.CatalogueService.Search(ctx, &SearchCatalogueRequest{
			Name:     "sear",
			PriceMax: 10000,
			PriceMin: 0,
		})
		if err != nil {
			t.Fatalf("Failed to search catalogue err: %v", err)
		}
		wantLen := 1
		if len(got.List) != wantLen {
			t.Errorf("Expected len: %v Got Len: %v", wantLen, len(got.List))
		}

	})

}
