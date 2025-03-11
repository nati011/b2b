package product

import (
	"context"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/domain/product"
)

var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewProduct(
		db.NewMock(),
		nil,
	)
}

func Test_Create_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
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
		}

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Errorf("Failed to create product err: %v", err)
		}

		//get
		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Errorf("Expected err:%v Got err: %v", nil, err)
		}
		if resp.Id == 0 {
			t.Errorf("No product created")
		}
	})

	t.Run("inactive_by_default", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Error("Failed to create product", err)
		}

		got, err := service.Get(ctx, id)
		if err != nil {
			t.Errorf("Expected err:%v Got err: %v", nil, err)
		}
		if got.Id == 0 {
			t.Errorf("No product created")
		}
		expectedActiveStatus := false
		if got.IsActive != expectedActiveStatus {
			t.Errorf("expected active status: %v, got active status: %v", expectedActiveStatus, got.IsActive)
		}
	})

}

func Test_Create_unhappyPath(t *testing.T) {
	t.Run("name_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
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
		_, err := service.Create(ctx, in)
		wantErr := ErrNameNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("name_duplicate", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
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
		_, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		in_new := &CreateRequest{
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
		_, err = service.Create(ctx, in_new)
		wantErr := ErrNameDuplicate
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("desc_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "t",
			ExternalID: "1",
			Images: []string{
				"test",
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrDescNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_two_images_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Two ImageTest",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrImagesMustBeAtleastTwo
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("attribute_values_cannot_be_empty", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Value EmptyTest",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrAttributeValuesCannotBeEmpty
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("price_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Price Test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Attributes: map[string]string{
				"test": "test",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrPriceNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("price_cannot_be_zero", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Price ZeroTest",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 0.00,
			Attributes: map[string]string{
				"test": "tets",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrPriceNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	in := &CreateRequest{
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

	_, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	//get
	got, err := service.GetAll(ctx)
	if err != nil {
		t.Errorf("Expected err:%v Got err: %v", nil, err)
	}
	wantNum := 1
	if len(got.List) != wantNum {
		t.Errorf("Expected len: %v, Got len: %v", wantNum, len(got.List))
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("no_product_found", func(t *testing.T) {
		ctx := context.Background()
		//get
		wantErr := ErrIdNotFound
		_, err := service.Get(ctx, 99)
		if err != wantErr {
			t.Errorf("Expected err:%v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_All_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	in := &CreateRequest{
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

	_, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	in_two := &CreateRequest{
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

	_, err = service.Create(ctx, in_two)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	//get-all
	resp, err := service.GetAll(ctx)
	if err != nil {
		t.Errorf("Expected err:%v Got err: %v", nil, err)
	}
	expected_resp_len := 2
	if len(resp.List) != expected_resp_len {
		t.Errorf("Expected resp len:%v Got resp len: %v", expected_resp_len, len(resp.List))
	}
}

func Test_Get_All_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//get-all
	wantErr := ErrEmptyGetContent
	_, err := service.GetAll(ctx)
	if err != wantErr {
		t.Errorf("Expected err:%v Got err: %v", ErrEmptyGetContent, err)
	}
}

func Test_Get_by_param_happyPath(t *testing.T) {

	t.Run("byName", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
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

		_, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := service.GetByParam(ctx, &GetByParamRequest{
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

	t.Run("byExternalId", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "External Id",
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

		_, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := service.GetByParam(ctx, &GetByParamRequest{
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

	t.Run("byPriceRange", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Price Range",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 500.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}

		_, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := service.GetByParam(ctx, &GetByParamRequest{
			PriceMin: 100,
			PriceMax: 1000,
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got err: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got err: %v", wantLen, len(got.List))
		}
	})

	t.Run("aggregate-Fetch", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Aggregate Fetch",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 900.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}

		in_2 := &CreateRequest{
			Name:       "test1",
			Desc:       "test1",
			ExternalID: "1233",
			Images: []string{
				"test",
				"test",
			},
			Price: 950.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}

		_, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		_, err = service.Create(ctx, in_2)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := service.GetByParam(ctx, &GetByParamRequest{
			PriceMin: 800,
			PriceMax: 1000,
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got: %v", nil, err)
		}
		wantLen := 2
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})

}

func Test_Get_by_param_unhappyPath(t *testing.T) {
	t.Run("noMatch", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.GetByParam(ctx, &GetByParamRequest{})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})
}

func Test_Update_happyPath(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		//setup
		ctx := context.Background()
		in := &CreateRequest{
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		//update
		update_in := &UpdateRequest{
			Id:         id,
			Name:       "test1",
			ExternalID: "321",
			Price:      200,
			Desc:       "updated",
			Images: []string{
				"updated",
				"updated",
			},
		}
		_, err = service.Update(ctx, update_in)
		if err != nil {
			t.Fatalf("Failed to update")
		}

		got, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get")
		}

		if got.Name != update_in.Name {
			t.Fatalf("Expected name: %v to Got name: %v", got.Name, update_in.Name)
		}

		if got.ExternalID != update_in.ExternalID {
			t.Fatalf("Expected extId: %v to Got extId: %v", got.ExternalID, update_in.ExternalID)
		}

		if got.Price != float64(update_in.Price) {
			t.Fatalf("Expected price: %v to Got price: %v", got.Price, update_in.Price)

		}
	})
}

func Test_Update_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		ctx := context.Background()
		//update
		update_in := &UpdateRequest{
			Id:         99,
			Name:       "test1",
			ExternalID: "321",
			Desc:       "1",
			Images: []string{
				"test",
				"test",
			},
			Price: 200,
		}
		wantErr := ErrIdNotFound
		_, err := service.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("name_mandatory", func(t *testing.T) {
		//setup
		ctx := context.Background()
		in := &CreateRequest{
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}
		//update
		update_in := &UpdateRequest{
			Id:         id,
			ExternalID: "321",
			Desc:       "1",
			Images: []string{
				"test",
				"test",
			},
			Price: 200,
		}
		wantErr := ErrNameNotSupplied
		_, err = service.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("desc_mandatory", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "testUnhappy Path",
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product %v", err)
		}
		// update
		update_in := &UpdateRequest{
			Id:         id,
			Name:       "testq",
			ExternalID: "321",
			Images: []string{
				"test",
				"test",
			},
			Price: 200,
		}
		wantErr := ErrDescNotSupplied
		_, err = service.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_two_images_mandatory", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Two Images Test",
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product %v", err)
		}
		// update
		update_in := &UpdateRequest{
			Id:         id,
			Name:       "testq",
			Desc:       "test",
			ExternalID: "321",
			Images: []string{
				"test",
			},
			Price: 200,
		}
		wantErr := ErrImagesMustBeAtleastTwo
		_, err = service.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("price_mandatory", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Price MandatoryTest",
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}
		// update
		update_in := &UpdateRequest{
			Id:         id,
			Name:       "testq",
			Desc:       "test",
			ExternalID: "321",
			Images: []string{
				"test",
				"test",
			},
		}
		wantErr := ErrPriceNotSupplied
		_, err = service.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("duplicate_name_not_allowed", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Duplicate Name tEST",
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		in_2 := &CreateRequest{
			Name:       "duplicate",
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

		_, err = service.Create(ctx, in_2)
		if err != nil {
			t.Fatalf("Failed to create product")
		}
		// update
		update_in := &UpdateRequest{
			Id:         id,
			Name:       "duplicate",
			ExternalID: "321",
			Desc:       "1",
			Images: []string{
				"test",
				"test",
			},
			Price: 200,
		}
		wantErr := ErrNameDuplicate
		_, err = service.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("price_cannot_be_zero", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "zERO PRICE",
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}
		// update
		update_in := &UpdateRequest{
			Id:         id,
			Name:       "testq",
			Desc:       "test",
			ExternalID: "321",
			Images: []string{
				"test",
				"test",
			},
			Price: 0,
		}
		wantErr := ErrPriceNotSupplied
		_, err = service.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})
}

func Test_goods_receiving_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	in := &CreateRequest{
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

	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	err = service.ReceiveGoods(ctx, &GoodsReceivingRequest{
		Id:     id,
		Amount: 2,
	})
	if err != nil {
		t.Fatalf("Failed to perform goods receiving err: %v", err)
	}

	//check stock
	got, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get product")
	}

	wantStock := 2
	if got.Stock != wantStock {
		t.Errorf("Expected stock: %v Got stock: %v", wantStock, got.Stock)
	}
}

func Test_goods_receiving_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		//setup
		ctx := context.Background()
		wantErr := ErrIdNotFound
		err := service.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     99,
			Amount: 2,
		})
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})
}

func Test_dispatch_happyPath(t *testing.T) {
	// setup
	ctx := context.Background()
	in := &CreateRequest{
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

	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	err = service.ReceiveGoods(ctx, &GoodsReceivingRequest{
		Id:     id,
		Amount: 2,
	})
	if err != nil {
		t.Fatalf("Failed to initiate goods receiving")
	}

	err = service.Dispatch(ctx, &DispatchRequest{
		Id:     id,
		Amount: 2,
	})
	if err != nil {
		t.Fatalf("Failed to initiate dispatch")
	}

	// check stock
	got, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get product")
	}

	wantStock := 2
	if got.Stock != wantStock {
		t.Errorf("Expected stock: %v Got stock: %v", wantStock, got.Stock)
	}
}

func Test_dispatch_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		ctx := context.Background()
		err := service.Dispatch(ctx, &DispatchRequest{
			Id:     99,
			Amount: 2,
		})
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Activate_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	in := &CreateRequest{
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

	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	err = service.Activate(ctx, id)
	if err != nil {
		t.Fatalf("Failed to activate product")
	}

	//verify status
	wantActiveStatus := true
	resp, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to activte status")
	}
	if resp.IsActive != wantActiveStatus {
		t.Errorf("Expected active status: %v, Want status: %v", wantActiveStatus, resp.IsActive)
	}
}

func Test_Activate_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		ctx := context.Background()
		err := service.Activate(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Fatalf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("already_active", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		err = service.Activate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to activate product")
		}

		err = service.Activate(ctx, id)
		if err != ErrAlreadyActive {
			t.Fatalf("Failed to activate product")
		}
	})
}

func Test_Deactvate_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	in := &CreateRequest{
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

	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	//verify
	got, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get")
	}
	wantActiveStatus := false
	if got.IsActive != wantActiveStatus {
		t.Errorf("Expected is active status: %v, Got is active status: %v", wantActiveStatus, got.IsActive)
	}
}

func Test_Deactvate_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		ctx := context.Background()
		err := service.Activate(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Want err: %v", wantErr, err)
		}
	})

	t.Run("already_inactive", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
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

		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		err = service.Deactivate(ctx, id)
		expectedErr := ErrAlreadyInactive
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})
}
