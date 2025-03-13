package order

import (
	"context"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/domain/order/db"
)

var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewOrderService(
		db.NewMock(),
		nil,
	)
}

func Test_Place_Order_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: 1,
			Items: []Item{
				{
					ProductId: 1,
					Quantity:  19},
			},
		}
		id, err := service.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to fetch order err: err %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected Id: %v Got Id: %v", id, resp.Id)
		}
	})

	t.Run("pendingStatusByDefault", func(t *testing.T) {
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: 1,
			Items: []Item{
				{
					ProductId: 1,
					Quantity:  19},
			},
		}
		id, err := service.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		resp, err := service.Get(ctx, id)
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
		ctx := context.Background()
		in := &PlaceRequest{
			Items: []Item{
				{
					ProductId: 1,
					Quantity:  19},
			},
		}
		_, err := service.Place(ctx, in)
		wantErr := ErrRetailerIdNotSupplied
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleastOneItemMandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: 1,
			Items:      []Item{},
		}
		_, err := service.Place(ctx, in)
		wantErr := ErrAtleastOneOrderItemNeeded
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("ItemParamsComplete", func(t *testing.T) {
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: 1,
			Items: []Item{
				{
					ProductId: 1},
			},
		}
		_, err := service.Place(ctx, in)
		wantErr := ErrItemMemberProductIdOrQuantityEmpty
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Cancel_Order_happyPath(t *testing.T) {
	ctx := context.Background()
	in := &PlaceRequest{
		RetailerId: 1,
		Items: []Item{
			{
				ProductId: 1,
				Quantity:  19},
		},
	}
	id, err := service.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}

	err = service.Cancel(ctx, id)
	if err != nil {
		t.Fatalf("Failed to cancel order err:%v", err)
	}

	got, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}

	if got.Status != CANCELD_STATUS {
		t.Errorf("Expected status: %v Got: %v", CANCELD_STATUS, got.Status)
	}
}

func Test_Cancel_Order_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		//check
		err := service.Cancel(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("alreadyCanceled", func(t *testing.T) {
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: 1,
			Items: []Item{
				{
					ProductId: 1,
					Quantity:  19},
			},
		}
		id, err := service.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}

		err = service.Cancel(ctx, id)
		if err != nil {
			t.Fatalf("Failed to cancel order err:%v", err)
		}

		//cancel again
		err = service.Cancel(ctx, id)
		wantErr := ErrAlreadyCanceled
		if err != wantErr {
			t.Errorf("Expected err: %v Got err:%v", wantErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	ctx := context.Background()
	in := &PlaceRequest{
		RetailerId: 1,
		Items: []Item{
			{
				ProductId: 1,
				Quantity:  19},
		},
	}
	id, err := service.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	//check
	got, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to fetch order err: err %v", err)
	}

	if got.Id != id {
		t.Errorf("Expected id: %v Got: %v", id, got.Id)
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		//check
		_, err := service.Get(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_All_happyPath(t *testing.T) {
	ctx := context.Background()
	in := &PlaceRequest{
		RetailerId: 1,
		Items: []Item{
			{
				ProductId: 1,
				Quantity:  19},
		},
	}
	_, err := service.Place(ctx, in)
	if err != nil {
		t.Fatalf("Failed to place order err: %v", err)
	}
	//check
	got, err := service.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to fetch order err: err %v", err)
	}
	wantLen := 1
	if len(got.List) != wantLen {
		t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
	}
}

func Test_Get_All_unhappyPath(t *testing.T) {
	t.Run("emptyGetContent", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.GetAll(ctx)
		wantErr := ErrEmptyGetResponse
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_By_Param_happyPath(t *testing.T) {
	t.Run("getByRetailerId", func(t *testing.T) {
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: 1,
			Items: []Item{
				{
					ProductId: 1,
					Quantity:  19},
			},
		}
		_, err := service.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		got, err := service.GetByParam(ctx, &GetByParamRequest{
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
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: 1,
			Items: []Item{
				{
					ProductId: 1,
					Quantity:  19},
			},
		}
		_, err := service.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		got, err := service.GetByParam(ctx, &GetByParamRequest{
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
		ctx := context.Background()
		in := &PlaceRequest{
			RetailerId: 1,
			Items: []Item{
				{
					ProductId: 1,
					Quantity:  19},
			},
		}
		_, err := service.Place(ctx, in)
		if err != nil {
			t.Fatalf("Failed to place order err: %v", err)
		}
		//check
		got, err := service.GetByParam(ctx, &GetByParamRequest{
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
		ctx := context.Background()
		_, err := service.GetByParam(ctx, &GetByParamRequest{
			Status:     PENDING_STATUS,
			RetailerId: 1,
		})
		wantErr := ErrEmptyGetResponse
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
	})
}
