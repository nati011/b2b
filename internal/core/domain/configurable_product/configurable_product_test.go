package configurable_product

import (
	"context"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/domain/catalogue/configurable_product"
	"b2b.nati011.github.com/internal/core/domain/product"
)

var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewConfigurableProductService(
		db.NewMock(),
		product.NewPackageIntegrationTestContainer().ProductService,
	)
}
func Test_Create_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		got, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if got.Id != id {
			t.Errorf("Expected id: %v Got id: %v", id, got.Id)
		}
	})

	t.Run("inactiveByDefault", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		got, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantIsAvailableStatus := false
		if got.IsAvailable != wantIsAvailableStatus {
			t.Errorf("Expected is available status: %v Got status: %v", wantIsAvailableStatus, got.IsAvailable)
		}
	})
}

func Test_Create_unhappyPath(t *testing.T) {
	t.Run("name_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrNameIsNotSupplied
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("duplicate_name", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		_, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		in_new := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
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
			Name:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrDescIsNotSupplied
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_one_attributeKey_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:          "test",
			Desc:          "test",
			ExternalId:    "test",
			AttributeKeys: []string{},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrAttributeKeysMustBeAtleastOne
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_one_product_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"1",
			},
			Products: []int{},
			Images: []string{
				"test",
				"test",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrProductsMustBeAtleastOne
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_two_images_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"1",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
			},
		}
		_, err := service.Create(ctx, in)
		wantErr := ErrImagesMustBeAtleastTwo
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Avail_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"1",
		},
		Products: []int{
			1,
		},
		Images: []string{
			"test",
			"test",
		},
	}
	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = service.Avail(ctx, id)
	if err != nil {
		t.Errorf("Failed to avail err: %v", err)
	}

	expectedStatus := false
	got, err := service.Get(ctx, id)
	if err != nil {
		t.Errorf("Expected status: %v Got status: %v", expectedStatus, got.IsAvailable)
	}

}

func Test_Avail_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		ctx := context.Background()
		err := service.Avail(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("already_available", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"1",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		err = service.Avail(ctx, id)
		if err != nil {
			t.Errorf("Failed to avail err: %v", err)
		}

		err = service.Avail(ctx, id)
		expectedErr := ErrAlreadyAvailable
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})
}

func Test_Disable_happyPath(t *testing.T) {
	// setup
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"1",
		},
		Products: []int{
			1,
		},
		Images: []string{
			"test",
			"test",
		},
	}
	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	err = service.Avail(ctx, id)
	if err != nil {
		t.Errorf("Failed to avail err: %v", err)
	}

	err = service.Disable(ctx, id)
	if err != nil {
		t.Fatalf("Failed to disavail err: %v", err)
	}

	got, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}

	wantStatus := false
	if got.IsAvailable != wantStatus {
		t.Errorf("Expected status: %v Got status: %v", wantStatus, got.IsAvailable)
	}
}

func Test_Disable_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		ctx := context.Background()
		err := service.Disable(ctx, 99)
		expectedErr := ErrIdNotFound
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

	t.Run("already_disabled", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"1",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		err = service.Disable(ctx, id)
		expectedErr := ErrAlreadyUnavailable
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	//setup
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"test",
			"test",
		},
		Products: []int{
			1,
		},
		Images: []string{
			"test",
			"test",
		},
	}
	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//Get
	got, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to Get err: %v", err)
	}

	if got.Id != id {
		t.Errorf("Expected resp Id: %v Got: %v", id, got.Id)
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		// Get
		got, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to Get err: %v", err)
		}

		if got.Id != id {
			t.Errorf("Expected resp Id: %v Got: %v", id, got.Id)
		}
	})

	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		// Get
		wantErr := ErrIdNotFound
		_, err := service.Get(ctx, 99)
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

}

func Test_GetByParam_happyPath(t *testing.T) {
	t.Run("byName", func(t *testing.T) {
		//setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//get
		got, err := service.GetByParam(ctx, &GetByParamRequest{
			Name: "test",
		})
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}
		if got.List[0].Id != id {
			t.Errorf("Expected Id: %v, Got Id: %v", id, got.List[0].Id)
		}
	})

	t.Run("byExtId", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		// get
		got, err := service.GetByParam(ctx, &GetByParamRequest{
			ExternalId: "test",
		})
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}
		if got.List[0].Id != id {
			t.Errorf("Expected Id: %v, Got Id: %v", id, got.List[0].Id)
		}
	})
}

func Test_GetByParam_unhappyPath(t *testing.T) {
	t.Run("emptyGetContent", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.GetByParam(ctx, &GetByParamRequest{
			Name: "test",
		})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_All_happyPath(t *testing.T) {
	// setup
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
			"test",
			"test",
		},
		Products: []int{
			1,
		},
		Images: []string{
			"test",
			"test",
		},
	}
	_, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	// get
	got, err := service.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get by param err: %v", err)
	}
	wantLen := 1
	if len(got.List) != wantLen {
		t.Errorf("Expected len: %v Got len: %v", wantLen, len(got.List))
	}
}

func Test_Get_All_unhappyPath(t *testing.T) {
	t.Run("emptyContent", func(t *testing.T) {
		ctx := context.Background()
		wantErr := ErrEmptyGetContent
		_, err := service.GetAll(ctx)
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Update_happyPath(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//updateName
		err = service.Update(ctx, &UpdateRequest{
			Id:   id,
			Name: "updated",
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantName := "updated"
		if resp.Name != wantName {
			t.Errorf("Expected name: %v Got: %v", wantName, resp.Name)
		}
	})

	t.Run("desc", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//updateDesc
		err = service.Update(ctx, &UpdateRequest{
			Id:   id,
			Desc: "updated",
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantName := "updated"
		if resp.Desc != wantName {
			t.Errorf("Expected name: %v Got: %v", wantName, resp.Name)
		}
	})

	t.Run("extId", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//updateExtId
		err = service.Update(ctx, &UpdateRequest{
			Id:         id,
			ExternalId: "updated",
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantName := "updated"
		if resp.ExternalId != wantName {
			t.Errorf("Expected extId: %v Got: %v", "updated", resp.Name)
		}
	})

	t.Run("product", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//wantProduct
		err = service.Update(ctx, &UpdateRequest{
			Id: id,
			Product: []int{
				9,
			},
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantProduct := 9
		if resp.Products[0] != wantProduct {
			t.Errorf("Expected product: %v Got: %v", wantProduct, resp.Products[0])
		}
	})

	t.Run("isAvailableStatus", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//isAvailableStatus
		err = service.Update(ctx, &UpdateRequest{
			Id:                id,
			IsAvailableStatus: false,
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantIsAvailableStatus := false
		if resp.IsAvailable != wantIsAvailableStatus {
			t.Errorf("Expected product: %v Got: %v", wantIsAvailableStatus, resp.IsAvailable)
		}
	})

	t.Run("images", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//images
		err = service.Update(ctx, &UpdateRequest{
			Id: id,
			Images: []string{
				"updated",
				"updated",
			},
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantImages := "updated"
		if resp.Images[0] != wantImages && resp.Images[1] != wantImages {
			t.Errorf("Expected image: %v Got: %v", wantImages, resp.Images[0])
		}
	})

	t.Run("attributes", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
				"test",
			},
			Products: []int{
				1,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//attribute keys
		err = service.Update(ctx, &UpdateRequest{
			Id: id,
			AttributeKeys: []string{
				"updated",
			},
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantAttributeKey := "updated"
		if resp.Attributes[wantAttributeKey] == "" {
			t.Errorf("Expected attributeKey != emptyString Got: emptySting")
		}
	})
}

func Test_Update_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		err := service.Update(ctx, &UpdateRequest{
			Id: 99,
			AttributeKeys: []string{
				"updated",
			},
		})
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}
