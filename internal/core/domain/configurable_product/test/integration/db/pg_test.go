package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/configurable_product"
	test_container "b2b.nati011.github.com/internal/core/domain/configurable_product/test/integration"
	"b2b.nati011.github.com/internal/core/domain/product"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var container test_container.TestContainer
var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = test_container.NewDBIntegrationTestContainer(db)
	db = db_test_container.Setup()
}

func teardown() {
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_read(t *testing.T) {
	t.Run("GetByName", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		product_id, _ := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})
		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//get
		got, err := container.ConfigurableProductService.GetByParam(ctx, &configurable_product.GetByParamRequest{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}
		if got.List[0].Id != id {
			t.Errorf("Expected Id: %v, Got Id: %v", id, got.List[0].Id)
		}
	})

	t.Run("GetByExternalId", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
		ctx := context.Background()
		product_id, _ := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})
		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		// get
		got, err := container.ConfigurableProductService.GetByParam(ctx, &configurable_product.GetByParamRequest{
			ExternalId: "test",
		})
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}
		if got.List[0].Id != id {
			t.Errorf("Expected Id: %v, Got Id: %v", id, got.List[0].Id)
		}
	})

	t.Run("GetById", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		product_id, _ := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})
		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		// Get
		got, err := container.ConfigurableProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get err: %v", err)
		}

		if got.Id != id {
			t.Errorf("Expected resp Id: %v Got: %v", id, got.Id)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		product_id, _ := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})
		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		_, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		// get
		got, err := container.ConfigurableProductService.GetAll(ctx)
		if err != nil {
			t.Fatalf("Failed to get all err: %v", err)
		}
		wantLen := 1
		if len(got.List) != wantLen {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(got.List))
		}
	})
}

func Test_write(t *testing.T) {

	t.Run("Create", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		product_id, _ := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})

		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		got, err := container.ConfigurableProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if got.Id != id {
			t.Errorf("Expected id: %v Got id: %v", id, got.Id)
		}
	})

	t.Run("updateName", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		product_id, err := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})
		if err != nil {
			t.Fatalf("Failed to update name %v", err)
		}
		//setup
		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//updateName
		err = container.ConfigurableProductService.Update(ctx, &configurable_product.UpdateRequest{
			Id:   id,
			Name: "updated",
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := container.ConfigurableProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantName := "updated"
		if resp.Name != wantName {
			t.Errorf("Expected name: %v Got: %v", wantName, resp.Name)
		}
	})

	t.Run("updateDesc", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		//setup
		product_id, _ := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})
		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//updateDesc
		err = container.ConfigurableProductService.Update(ctx, &configurable_product.UpdateRequest{
			Id:   id,
			Desc: "updated",
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := container.ConfigurableProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantName := "updated"
		if resp.Desc != wantName {
			t.Errorf("Expected name: %v Got: %v", wantName, resp.Name)
		}
	})

	t.Run("updateExternalId", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		//setup
		product_id, _ := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})
		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//updateExtId
		err = container.ConfigurableProductService.Update(ctx, &configurable_product.UpdateRequest{
			Id:         id,
			ExternalId: "updated",
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := container.ConfigurableProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantName := "updated"
		if resp.ExternalId != wantName {
			t.Errorf("Expected extId: %v Got: %v", "updated", resp.Name)
		}
	})

	t.Run("updateProducts", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		//setup
		product_id, _ := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})
		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//wantProduct
		err = container.ConfigurableProductService.Update(ctx, &configurable_product.UpdateRequest{
			Id: id,
			Product: []int{
				product_id,
			},
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := container.ConfigurableProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantProduct := product_id
		if resp.Products[0] != wantProduct {
			t.Errorf("Expected product: %v Got: %v", wantProduct, resp.Products[0])
		}
	})

	t.Run("updateImages", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		product_id, _ := container.ProductService.Create(ctx, &product.CreateRequest{
			Name:       "testProduct",
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
		})
		in := &configurable_product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := container.ConfigurableProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//images
		err = container.ConfigurableProductService.Update(ctx, &configurable_product.UpdateRequest{
			Id: id,
			Images: []string{
				"updated",
				"updated",
			},
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := container.ConfigurableProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantImages := "updated"
		if resp.Images[0] != wantImages && resp.Images[1] != wantImages {
			t.Errorf("Expected image: %v Got: %v", wantImages, resp.Images[0])
		}
	})
}
