package catalogue

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/category"
	product "b2b.nati011.github.com/internal/core/domain/product"
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

func Test_Add_Category_To_Product_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//create category
	category_id, err := testContainer.CategoryService.Create(ctx, &category.CreateRequest{
		Name: "test",
	})
	if err != nil {
		t.Errorf("Failed to create category")
	}

	//create product with category
	in_product := &product.CreateRequest{
		Name:       "test",
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
		CategoryId: []int{
			category_id,
		},
	}
	product_id, err := testContainer.ProductService.Create(ctx, in_product)
	if err != nil {
		t.Fatalf("Failed to create product err: %v", err)
	}

	product_resp, err := testContainer.ProductService.Get(ctx, product_id)
	if err != nil {
		t.Fatalf("Failed to get category err: %v", err)
	}
	found := false
	for _, i := range product_resp.CategoryId {
		if i == category_id {
			found = true
		}
	}
	if !found {
		t.Errorf("Category missing from product category list")
	}
}

func Test_Add_Category_To_Product_unhappyPath(t *testing.T) {
	t.Run("category_not_found", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()

		//create product with category
		in_product := &product.CreateRequest{
			Name:       "test",
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
			CategoryId: []int{
				99,
			},
		}
		wantErr := product.ErrCategoryNotFound
		_, err := testContainer.ProductService.Create(ctx, in_product)
		if err != wantErr {
			t.Errorf("Failed to create product err: %v", err)
		}

	})
}

func Test_Update_product_Category_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//create category
	category_id, err := testContainer.CategoryService.Create(ctx, &category.CreateRequest{
		Name: "test",
	})
	if err != nil {
		t.Fatalf("Failed to create category %v", err)
	}

	//create product with category
	in_product := &product.CreateRequest{
		Name:       "test",
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
		CategoryId: []int{
			category_id,
		},
	}
	product_id, err := testContainer.ProductService.Create(ctx, in_product)
	if err != nil {
		t.Fatalf("Failed to create product err: %v", err)
	}

	//create category
	new_category_id, err := testContainer.CategoryService.Create(ctx, &category.CreateRequest{
		Name: "test2",
	})
	if err != nil {
		t.Errorf("Failed to create category")
	}

	//update product with category
	update_product := &product.UpdateRequest{
		Id:         product_id,
		Name:       "test1",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"test",
			"test",
		},
		Price: 100.00,
		CategoryId: []int{
			category_id,
			new_category_id,
		},
	}

	_, err = testContainer.ProductService.Update(ctx, update_product)
	if err != nil {
		t.Fatalf("Failed to update product err: %v", err)
	}

	got, err := testContainer.ProductService.Get(ctx, product_id)
	if err != nil {
		t.Fatalf("Failed to update product err: %v", err)
	}

	//verify
	expectedCategories := []int{
		category_id,
		new_category_id,
	}

	Notfound := false
	for _, i := range expectedCategories {
		if i != got.CategoryId[0] && i != got.CategoryId[1] {
			Notfound = true
		}
	}

	if Notfound {
		t.Errorf("Expected categories: %v Got categories %v", expectedCategories, got.CategoryId)
	}
}

func Test_Update_Product_Category_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//create category
	category_id, err := testContainer.CategoryService.Create(ctx, &category.CreateRequest{
		Name: "test",
	})
	if err != nil {
		t.Errorf("Failed to create category")
	}

	//create product with category
	in_product := &product.CreateRequest{
		Name:       "test",
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
		CategoryId: []int{
			category_id,
		},
	}
	product_id, err := testContainer.ProductService.Create(ctx, in_product)
	if err != nil {
		t.Fatalf("Failed to create product err: %v", err)
	}

	//update product with category
	update_product := &product.UpdateRequest{
		Id:         product_id,
		Name:       "test1",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"test",
			"test",
		},
		Price: 100.00,
		CategoryId: []int{
			category_id,
			99,
		},
	}

	wantErr := product.ErrCategoryNotFound
	_, err = testContainer.ProductService.Update(ctx, update_product)
	if err != wantErr {
		t.Fatalf("Expected err: %v Got err %v", wantErr, err)
	}
}

func Test_Get_Products_By_Category_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//create category
	category_id, err := testContainer.CategoryService.Create(ctx, &category.CreateRequest{
		Name: "test",
	})
	if err != nil {
		t.Errorf("Failed to create category")
	}

	//create product with category
	in_product := &product.CreateRequest{
		Name:       "test",
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
		CategoryId: []int{
			category_id,
		},
	}
	product_id, err := testContainer.ProductService.Create(ctx, in_product)
	if err != nil {
		t.Fatalf("Failed to create product err: %v", err)
	}

	//create product with category
	new_in_product := &product.CreateRequest{
		Name:       "test2",
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
		CategoryId: []int{
			category_id,
		},
	}
	new_product_id, err := testContainer.ProductService.Create(ctx, new_in_product)
	if err != nil {
		t.Fatalf("Failed to create product err: %v", err)
	}

	expectedProducts := []int{
		product_id,
		new_product_id,
	}

	got, err := testContainer.ProductService.GetByParam(ctx, &product.GetByParamRequest{
		CategoryId: []int{
			category_id,
		}})
	if err != nil {
		t.Errorf("Failed to get products with categories: %v", err)
	}

	NotFound := false
	for _, i := range got.List {
		if i.Id != expectedProducts[0] && i.Id != expectedProducts[1] {
			NotFound = true
		}
	}

	if NotFound {
		t.Errorf("Expected products: %v Got Products %v", expectedProducts, nil)
	}

}

func Test_Get_Products_By_Category_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()

	//create product with category
	in_product := &product.CreateRequest{
		Name:       "test2",
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
	}
	_, err := testContainer.ProductService.Create(ctx, in_product)
	if err != nil {
		t.Fatalf("Failed to create product err: %v", err)
	}

	wantErr := product.ErrCategoryNotFound
	_, err = testContainer.ProductService.GetByParam(ctx, &product.GetByParamRequest{
		CategoryId: []int{
			99,
		}})
	if err != wantErr {
		t.Errorf("Expected err: %v Got err: %v", wantErr, err)
	}
}

func Test_Get_Categories_of_product_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	//create category
	category_id, err := testContainer.CategoryService.Create(ctx, &category.CreateRequest{
		Name: "test",
	})
	if err != nil {
		t.Errorf("Failed to create category")
	}

	//create product with category
	in_product := &product.CreateRequest{
		Name:       "test",
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
		CategoryId: []int{
			category_id,
		},
	}
	product_id, err := testContainer.ProductService.Create(ctx, in_product)
	if err != nil {
		t.Fatalf("Failed to create product err: %v", err)
	}

	product, err := testContainer.ProductService.Get(ctx, product_id)
	if err != nil {
		t.Fatalf("Failed to fetch categories")
	}

	if product.CategoryId[0] != category_id {
		t.Errorf("Expected categoryId: %v got: %v", category_id, product.CategoryId[0])
	}

}
