package db

import (
	"context"
	"database/sql"
	"math"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/product"
	test_container "b2b.nati011.github.com/internal/core/domain/product/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var container test_container.TestContainer
var db *sql.DB
var categoryId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = test_container.NewDBIntegrationTestContainer(
		db,
	)
	ctx := context.Background()
	categoryId, _ = container.CategoryService.Create(ctx, &category.CreateRequest{
		Name: "test",
	})
	db = db_test_container.Setup()
}

func teardown() {
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

func Test_read(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get product err %v", err)
		}
		wantId := id
		if got.Id != wantId {
			t.Errorf("Expected Id: %v, Got: %v", wantId, id)
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		in_two := &product.CreateRequest{
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

		_, err = container.ProductService.Create(ctx, in_two)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		//get-all
		resp, err := container.ProductService.GetAll(ctx)
		if err != nil {
			t.Errorf("Expected err:%v Got err: %v", nil, err)
		}
		expected_resp_len := 2
		if len(resp.List) != expected_resp_len {
			t.Errorf("Expected resp len:%v Got resp len: %v", expected_resp_len, len(resp.List))
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
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

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := container.ProductService.GetByParam(ctx, &product.GetByParamRequest{
			Name: "test",
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got err: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got err: %v", wantLen, len(got.List))
		}
	})

	t.Run("GetByExternalId", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
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

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := container.ProductService.GetByParam(ctx, &product.GetByParamRequest{
			ExternalID: "123",
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got err: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got err: %v", wantLen, len(got.List))
		}
	})

	t.Run("GetByDistributorId", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
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
			DistributorId: 1,
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}

		got, err := container.ProductService.GetByParam(ctx, &product.GetByParamRequest{
			DistributorId: 1,
		})
		if err != nil {
			t.Fatalf("Failed to get product err: %v", err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})

	t.Run("GetByCategory", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
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
			DistributorId: 1,
			CategoryId:    []int{categoryId},
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}

		got, err := container.ProductService.GetByParam(ctx, &product.GetByParamRequest{
			CategoryId: []int{categoryId},
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got err: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})

	t.Run("GetByPriceRange", func(t *testing.T) {
		t.Cleanup(teardown)
		// setup
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
			DistributorId: 1,
			CategoryId:    []int{categoryId},
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}

		got, err := container.ProductService.GetByParam(ctx, &product.GetByParamRequest{
			PriceMin: math.MinInt,
			PriceMax: math.MaxInt,
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got err: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})
}

func Test_write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Errorf("Failed to create product err: %v", err)
		}

		//get
		resp, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Errorf("Expected err:%v Got err: %v", nil, err)
		}
		if resp.Id == 0 {
			t.Errorf("No product created")
		}
	})

	t.Run("updateName", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err %v", err)
		}

		//update
		update_in := &product.UpdateRequest{
			Id:   id,
			Name: "test1",
		}
		_, err = container.ProductService.Update(ctx, update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}

		if got.Name != update_in.Name {
			t.Fatalf("Expected name: %v Got: %v", got.Name, update_in.Name)
		}
	})

	t.Run("updateExternalId", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		//update
		update_in := &product.UpdateRequest{
			Id:         id,
			ExternalID: "test1",
		}
		_, err = container.ProductService.Update(ctx, update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}

		if got.ExternalID != update_in.ExternalID {
			t.Fatalf("Expected externalId: %v Got: %v", got.ExternalID, update_in.ExternalID)
		}
	})

	t.Run("updatePrice", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		//update
		update_in := &product.UpdateRequest{
			Id:    id,
			Price: 2,
		}
		_, err = container.ProductService.Update(ctx, update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}

		if got.Price != float64(update_in.Price) {
			t.Fatalf("Expected price: %v Got: %v", got.Price, update_in.Price)
		}
	})
	t.Run("updateDesc", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		//update
		update_in := &product.UpdateRequest{
			Id:   id,
			Desc: "new desc",
		}
		_, err = container.ProductService.Update(ctx, update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}

		if got.Desc != update_in.Desc {
			t.Fatalf("Expected desc: %v Got: %v", got.Desc, update_in.Desc)
		}
	})
	t.Run("updateImages", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		//update
		update_in := &product.UpdateRequest{
			Id: id,
			Images: []string{
				"new image",
				"new image",
				"new image",
			},
		}
		_, err = container.ProductService.Update(ctx, update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}

		if len(got.Images) != len(update_in.Images) {
			t.Fatalf("Expected images len: %v Got: %v", len(got.Images), len(update_in.Images))
		}
	})

	t.Run("updateActiveStatus", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		err = container.ProductService.Activate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}
		wantActiveStatus := true
		if got.IsActive != wantActiveStatus {
			t.Fatalf("Expected active status: %v Got: %v", wantActiveStatus, got.IsActive)
		}

		err = container.ProductService.Deactivate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err = container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}
		wantActiveStatus = false
		if got.IsActive != wantActiveStatus {
			t.Fatalf("Expected active status: %v Got: %v", wantActiveStatus, got.IsActive)
		}
	})

	t.Run("updateCategory", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		//update
		update_in := &product.UpdateRequest{
			Id: id,
			CategoryId: []int{
				categoryId,
				categoryId,
				categoryId,
			},
		}
		_, err = container.ProductService.Update(ctx, update_in)
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}

		if len(got.CategoryId) != len(update_in.CategoryId) {
			t.Fatalf("Expected categories len: %v Got: %v", len(got.CategoryId), len(update_in.CategoryId))
		}
	})
	t.Run("goodsReceiving", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		goodsReceivingAmount := 5
		err = container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
			Id:     id,
			Amount: goodsReceivingAmount,
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}

		if got.Stock != goodsReceivingAmount {
			t.Fatalf("Expected stock: %v Got: %v", goodsReceivingAmount, got.Stock)
		}
	})

	t.Run("dispatch", func(t *testing.T) {
		t.Cleanup(teardown)
		//setup
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

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		goodsReceivingAmount := 5
		err = container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
			Id:     id,
			Amount: goodsReceivingAmount,
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		dispachAmount := 5
		err = container.ProductService.Dispatch(ctx, &product.DispatchRequest{
			Id:     id,
			Amount: dispachAmount,
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		got, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}

		if got.Stock != 0 {
			t.Fatalf("Expected stock: %v Got: %v", dispachAmount, got.Stock)
		}
	})
}
