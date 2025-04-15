package order

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

var container TestContainer
var retailer_id int
var product_id int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	container = NewPackageIntegrationTestContainer()
	retailer_id, _ = container.RetailerService.Create(ctx, &retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",

		FirstName: "test",
		LastName:  "test",

		Email: "test@gmail.com",
	})

	product_id, _ = container.ProductService.Create(ctx, &product.CreateRequest{
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
	container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     product_id,
		Amount: 100,
	})
}

func teardown() {
	container.Teardown()
}

func Test_Place_Order_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: retailer_id,
			Items: []Item{
				{
					ProductId: product_id,
					Quantity:  19},
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

	t.Run("pendingStatusByDefault", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: retailer_id,
			Items: []Item{
				{
					ProductId: product_id,
					Quantity:  19},
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

		if resp.Status != PENDING_STATUS {
			t.Errorf("Expected status: %v, Got status: %v", PENDING_STATUS, resp.Status)
		}
	})

}

func Test_Place_Order_unhappyPath(t *testing.T) {
	t.Run("retailerIdMandatory", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &PlaceRequest{
			Items: []Item{
				{
					ProductId: product_id,
					Quantity:  19},
			},
		}
		_, err := container.OrderService.Place(ctx, in)
		wantErr := ErrRetailerIdNotSupplied
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleastOneItemMandatory", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: retailer_id,
			Items:      []Item{},
		}
		_, err := container.OrderService.Place(ctx, in)
		wantErr := ErrAtleastOneOrderItemNeeded
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("ItemParamsComplete", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: retailer_id,
			Items: []Item{
				{
					ProductId: product_id},
			},
		}
		_, err := container.OrderService.Place(ctx, in)
		wantErr := ErrItemMemberProductIdOrQuantityEmpty
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Cancel_Order_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := &PlaceRequest{
		RetailerId: retailer_id,
		Items: []Item{
			{
				ProductId: product_id,
				Quantity:  19},
		},
	}
	id, err := container.OrderService.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}

	err = container.OrderService.Cancel(ctx, id)
	if err != nil {
		t.Fatalf("Failed to cancel order err:%v", err)
	}

	got, err := container.OrderService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}

	if got.Status != CANCELD_STATUS {
		t.Errorf("Expected status: %v Got: %v", CANCELD_STATUS, got.Status)
	}
}

func Test_Cancel_Order_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//check
		err := container.OrderService.Cancel(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("alreadyCanceled", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: retailer_id,
			Items: []Item{
				{
					ProductId: product_id,
					Quantity:  19},
			},
		}
		id, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}

		err = container.OrderService.Cancel(ctx, id)
		if err != nil {
			t.Fatalf("Failed to cancel order err:%v", err)
		}

		//cancel again
		err = container.OrderService.Cancel(ctx, id)
		wantErr := ErrAlreadyCanceled
		if err != wantErr {
			t.Errorf("Expected err: %v Got err:%v", wantErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := &PlaceRequest{
		RetailerId: retailer_id,
		Items: []Item{
			{
				ProductId: product_id,
				Quantity:  19},
		},
	}
	id, err := container.OrderService.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	//check
	got, err := container.OrderService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to fetch order err: err %v", err)
	}

	if got.Id != id {
		t.Errorf("Expected id: %v Got: %v", id, got.Id)
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//check
		_, err := container.OrderService.Get(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_All_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := &PlaceRequest{
		RetailerId: retailer_id,
		Items: []Item{
			{
				ProductId: product_id,
				Quantity:  19},
		},
	}
	_, err := container.OrderService.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	pagination := &Pagination{
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
}

func Test_Get_All_unhappyPath(t *testing.T) {
	t.Run("emptyGetContent", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		pagination := &Pagination{
			Limit:  1,
			Offset: 0,
		}
		_, err := container.OrderService.GetAll(ctx, pagination)
		wantErr := ErrEmptyGetResponse
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_By_Param_happyPath(t *testing.T) {
	t.Run("getByRetailerId", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: retailer_id,
			Items: []Item{
				{
					ProductId: product_id,
					Quantity:  19},
			},
		}
		_, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		got, err := container.OrderService.GetByParam(ctx, &GetByParamRequest{
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

	t.Run("getByStatus", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: retailer_id,
			Items: []Item{
				{
					ProductId: product_id,
					Quantity:  19},
			},
		}
		_, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		got, err := container.OrderService.GetByParam(ctx, &GetByParamRequest{
			Status: PENDING_STATUS,
		})
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		wantLen := 1
		if len(got.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})

	t.Run("aggregate", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: retailer_id,
			Items: []Item{
				{
					ProductId: product_id,
					Quantity:  19},
			},
		}
		_, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		got, err := container.OrderService.GetByParam(ctx, &GetByParamRequest{
			Status:     PENDING_STATUS,
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
}

func Test_Get_By_Param_unhappyPath(t *testing.T) {
	t.Run("emptyGetContent", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := container.OrderService.GetByParam(ctx, &GetByParamRequest{
			Status:     PENDING_STATUS,
			RetailerId: 1,
		})
		wantErr := ErrEmptyGetResponse
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
	})
}

func Test_Update_Status(t *testing.T) {
	t.Run("update_status", func(t *testing.T) {
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: retailer_id,
			Items: []Item{
				{
					ProductId: product_id,
					Quantity:  19},
			},
		}
		id, err := container.OrderService.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}

		//update
		_, err = container.OrderService.UpdateStatus(ctx, &UpdateRequest{
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
