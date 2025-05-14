package configurable_product

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/configurable_product"

	"b2b.nati011.github.com/internal/core/domain/product"
)

var container configurable_product.TestContainer

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = configurable_product.NewPackageIntegrationTestContainer()
}

func Test_Create_ValidateProduct_happyPath(t *testing.T) {
	//create product
	t.Cleanup(container.Teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
		Name:       "test",
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
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}
	//create configurable product
	in_cp := &configurable_product.CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"test",
			"test",
		},
		Products: []int{
			id,
		},
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
	}
	_, err = container.ConfigurableProductService.Create(ctx, in_cp)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}
}

func Test_Create_ValidateProduct_unhappyPath(t *testing.T) {
	// create configurable product
	t.Cleanup(container.Teardown)
	setup()
	ctx := context.Background()
	in_cp := &configurable_product.CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"test",
			"test",
		},
		Products: []int{
			99,
		},
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
	}
	_, err := container.ConfigurableProductService.Create(ctx, in_cp)
	wantErr := configurable_product.ErrProductNotFound
	if err != wantErr {
		t.Errorf("Expcetd err: %v Got err: %v", wantErr, err)
	}
}
func Test_Create_ValidateAttribute_keys_happyPath(t *testing.T) {
	//create product
	t.Cleanup(container.Teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
		Name:       "test",
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
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}

	//create configurable product
	in_cp := &configurable_product.CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"test",
		},
		Products: []int{
			id,
		},
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
	}
	_, err = container.ConfigurableProductService.Create(ctx, in_cp)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}
}

func Test_Create_ValidateAttribute_keys_unhappyPath(t *testing.T) {
	//create product
	t.Cleanup(container.Teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
		Name:       "test",
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
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}

	//create configurable product
	in_cp := &configurable_product.CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"unknownKey",
		},
		Products: []int{
			id,
		},
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
	}
	_, err = container.ConfigurableProductService.Create(ctx, in_cp)
	wantErr := configurable_product.ErrAttributeKeysDoNotExistInProduct
	if err != wantErr {
		t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
	}
}

// func Test_Create_PopulateAttributeValues(t *testing.T) {
// 	//create product
// 	t.Cleanup(container.Teardown)
// 	setup()
// 	ctx := context.Background()
// 	in := &product.CreateRequest{
// 		Name:       "test",
// 		Desc:       "test",
// 		ExternalID: "123",
// 		Images: []string{
// 			"test",
// 			"test",
// 		},
// 		Price: 100.00,
// 		Attributes: map[string]string{
// 			"test": "test",
// 		},
// 	}

// 	id, err := container.ProductService.Create(ctx, in)
// 	if err != nil {
// 		t.Errorf("Failed to create product err: %v", err)
// 	}

// 	//create configurable product
// 	in_cp := &configurable_product.CreateRequest{
// 		Name:       "test",
// 		Desc:       "test",
// 		ExternalId: "test",
// 		AttributeKeys: []string{
// 			"test",
// 		},
// 		Products: []int{
// 			id,
// 		},
// 		Images: []string{
// 			"test",
// 			"test",
// 		},
// 	}
// 	cp_id, err := container.ConfigurableProductService.Create(ctx, in_cp)
// 	if err != nil {
// 		t.Fatalf("Failed to create config prod %v", err)
// 	}
// 	wantAttributes := map[string]string{
// 		"test": "test",
// 	}
// 	resp, err := container.ConfigurableProductService.Get(ctx, cp_id)
// 	if err != nil {
// 		t.Fatalf("Failed to fetch product err: %v", err)
// 	}
// 	for _, i := range resp.Attributes {
// 		if i["test"] != wantAttributes["test"] {
// 			t.Errorf("Failed to populate attributes")
// 		}
// 	}
// }

func Test_Update_ValidateProduct_happyPath(t *testing.T) {
	// create product
	t.Cleanup(container.Teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
		Name:       "test",
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
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}
	// create configurable product
	in_cp := &configurable_product.CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"test",
			"test",
		},
		Products: []int{
			id,
		},
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
	}
	config_product_id, err := container.ConfigurableProductService.Create(ctx, in_cp)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}
	// create configurable product
	updated_cp := &configurable_product.UpdateRequest{
		Id: config_product_id,
		Product: []int{
			id,
		},
	}
	err = container.ConfigurableProductService.Update(ctx, updated_cp)
	if err != nil {
		t.Fatalf("Failed to update product err: %v", err)
	}
}

func Test_Update_ValidateProduct_unhappyPath(t *testing.T) {
	// create product
	t.Cleanup(container.Teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
		Name:       "test",
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
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}
	// create configurable product
	in_cp := &configurable_product.CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"test",
			"test",
		},
		Products: []int{
			id,
		},
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
	}
	config_product_id, err := container.ConfigurableProductService.Create(ctx, in_cp)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}
	// create configurable product
	updated_cp := &configurable_product.UpdateRequest{
		Id: config_product_id,
		Product: []int{
			99,
		},
	}
	err = container.ConfigurableProductService.Update(ctx, updated_cp)
	wantErr := configurable_product.ErrProductNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v Got err %v", wantErr, err)
	}
}
