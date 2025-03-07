package configurable_product

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/configurable_product"

	"b2b.nati011.github.com/internal/core/domain/product"
)

var testContainer configurable_product.TestContainer
var productService product.Provider
var configurableProductService configurable_product.Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = configurable_product.NewPackageIntegrationTestContainer()
	productService = testContainer.ProductService
	configurableProductService = testContainer.ConfigurableProductService
}

func Test_Create_ValidateProduct_happyPath(t *testing.T) {
	//create product
	ctx := context.Background()
	in := &product.CreateRequest{
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
	}

	id, err := productService.Create(ctx, in)
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
			"test",
			"test",
		},
	}
	_, err = configurableProductService.Create(ctx, in_cp)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}
}

func Test_Create_ValidateProduct_unhappyPath(t *testing.T) {
	// create configurable product
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
			"test",
			"test",
		},
	}
	_, err := configurableProductService.Create(ctx, in_cp)
	wantErr := configurable_product.ErrProductNotFound
	if err != wantErr {
		t.Errorf("Expcetd err: %v Got err: %v", wantErr, err)
	}
}
func Test_Create_ValidateAttribute_keys_happyPath(t *testing.T) {
	//create product
	ctx := context.Background()
	in := &product.CreateRequest{
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
	}

	id, err := productService.Create(ctx, in)
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
			"test",
			"test",
		},
	}
	_, err = configurableProductService.Create(ctx, in_cp)
	if err != nil {
		t.Errorf("Failed to create product err: %v", err)
	}
}

func Test_Create_ValidateAttribute_keys_unhappyPath(t *testing.T) {
	//create product
	ctx := context.Background()
	in := &product.CreateRequest{
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
	}

	id, err := productService.Create(ctx, in)
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
			"test",
			"test",
		},
	}
	_, err = configurableProductService.Create(ctx, in_cp)
	wantErr := configurable_product.ErrAttributeKeysDoNotExistInProduct
	if err != wantErr {
		t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
	}
}

func Test_Create_PopulateAttributeValues(t *testing.T) {
	//create product
	ctx := context.Background()
	in := &product.CreateRequest{
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
	}

	id, err := productService.Create(ctx, in)
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
			"test",
			"test",
		},
	}
	cp_id, err := configurableProductService.Create(ctx, in_cp)
	if err != nil {
		t.Fatalf("Failed to create config prod %v", err)
	}
	wantAttributes := map[string]string{
		"test": "test",
	}
	resp, err := configurableProductService.Get(ctx, cp_id)
	if err != nil {
		t.Fatalf("Failed to fetch product err: %v", err)
	}
	for _, i := range resp.Attributes {
		if i != wantAttributes["test"] {
			t.Errorf("Failed to populate attributes")
		}
	}
}

func Test_Update_ValidateProduct_happyPath(t *testing.T) {
	// create product
	ctx := context.Background()
	in := &product.CreateRequest{
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
	}

	id, err := productService.Create(ctx, in)
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
			"test",
			"test",
		},
	}
	config_product_id, err := configurableProductService.Create(ctx, in_cp)
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
	err = configurableProductService.Update(ctx, updated_cp)
	if err != nil {
		t.Fatalf("Failed to update product err: %v", err)
	}
}

func Test_Update_ValidateProduct_unhappyPath(t *testing.T) {
	// create product
	ctx := context.Background()
	in := &product.CreateRequest{
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
	}

	id, err := productService.Create(ctx, in)
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
			"test",
			"test",
		},
	}
	config_product_id, err := configurableProductService.Create(ctx, in_cp)
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
	err = configurableProductService.Update(ctx, updated_cp)
	wantErr := configurable_product.ErrProductNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v Got err %v", wantErr, err)
	}
}
