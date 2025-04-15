package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/order"
	test_container "b2b.nati011.github.com/internal/core/domain/order/test"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var container test_container.TestContainer
var db *sql.DB
var retailer_id int
var product_id int
var setupError error

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	db = db_test_container.Setup()
	container = test_container.NewDBIntegrationTestContainer(
		db,
	)
	// setup
	retailer_id, setupError = container.RetailerService.Create(ctx, &retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "test_user1",
		FirstName:   "test",
		LastName:    "test",

		Email: "test@gmail.com",
	})

	if setupError != nil {
		print(setupError.Error())
		panic("Failed to create retailer")
	}

	product_id, setupError = container.ProductService.Create(ctx, &product.CreateRequest{
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

	if setupError != nil {
		print(setupError.Error())
		panic("Failed to create product")
	}
	container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     product_id,
		Amount: 100,
	})
}

func teardown() {
	container.Teardown(db)
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {

}

func Test_Read(t *testing.T) {
	t.Run("get_by_id", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
			Id:     product_id,
			Amount: 100,
		})

		in := &order.PlaceRequest{
			RetailerId: retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		order_id, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		got, err := container.OrderService.Get(ctx, order_id)
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		wantId := order_id
		if got.Id != wantId {
			t.Errorf("Expected id: %v Got: %v", wantId, got.Id)
		}

	})

	t.Run("get_by_retailerId", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := &order.PlaceRequest{
			RetailerId: retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		_, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		got, err := container.OrderService.GetByParam(ctx, &order.GetByParamRequest{
			RetailerId: 1,
		})
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		wantLen := 1
		if len(got.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})

	t.Run("get_by_status", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := &order.PlaceRequest{
			RetailerId: retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		order_id, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//update
		_, err = container.OrderService.UpdateStatus(ctx, &order.UpdateRequest{
			Id:     order_id,
			Status: "test",
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}
		//check
		got, err := container.OrderService.GetByParam(ctx, &order.GetByParamRequest{
			Status: "test",
		})
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		wantLen := 1
		if len(got.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})

	t.Run("get_all", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := &order.PlaceRequest{
			RetailerId: retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		_, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}

		pagination := &order.Pagination{
			Limit:  1,
			Offset: 0,
		}
		//check
		got, err := container.OrderService.GetAll(ctx, pagination)
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		wantLen := pagination.Limit
		if len(got.List) != wantLen && len(got.List) > wantLen {
			t.Errorf("Expected length: %v Want: %v", wantLen, len(got.List))
		}
	})
}

func Test_Write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := &order.PlaceRequest{
			RetailerId: retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		id, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		resp, err := container.OrderService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected Id: %v Got Id: %v", id, resp.Id)
		}
	})

	t.Run("update_status", func(t *testing.T) {
		t.Cleanup(teardown)
		setup()
		ctx := context.Background()
		in := &order.PlaceRequest{
			RetailerId: retailer_id,
			Items: []order.Item{
				{
					ProductId: product_id,
					Quantity:  1},
			},
		}
		id, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}

		//update
		_, err = container.OrderService.UpdateStatus(ctx, &order.UpdateRequest{
			Id:     id,
			Status: "New",
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		//check
		resp, err := container.OrderService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		wantStatus := "New"
		if resp.Status != wantStatus {
			t.Errorf("Expected status: %v got: %v", wantStatus, resp.Status)
		}
	})
}
