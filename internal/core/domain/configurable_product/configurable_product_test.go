package configurable_product

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/product"
)

var product_id int
var container TestContainer

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	container = NewPackageIntegrationTestContainer()
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
			"test":         "test",
			"another_test": "another_test",
		},
	})
}

func Test_Create_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
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

	t.Run("inactiveByDefault", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()

		in := &CreateRequest{
			Name:       "test1",
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

		wantIsAvailableStatus := false
		if got.IsAvailable != wantIsAvailableStatus {
			t.Errorf("Expected is available status: %v Got status: %v", wantIsAvailableStatus, got.IsAvailable)
		}
	})
}

func Test_Create_unhappyPath(t *testing.T) {
	t.Run("name_mandatory", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		in := &CreateRequest{
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
		wantErr := ErrNameIsNotSupplied
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("duplicate_name", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
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

		in_new := &CreateRequest{
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
		_, err = container.ConfigurableProductService.Create(ctx, in_new)
		wantErr := ErrNameDuplicate
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("desc_mandatory", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
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
		wantErr := ErrDescIsNotSupplied
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_one_attributeKey_mandatory", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		in := &CreateRequest{
			Name:          "test",
			Desc:          "test",
			ExternalId:    "test",
			AttributeKeys: []string{},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
				"test",
			},
		}
		_, err := container.ConfigurableProductService.Create(ctx, in)
		wantErr := ErrAttributeKeysMustBeAtleastOne
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_one_product_mandatory", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
			},
			Products: []int{},
			Images: []string{
				"test",
				"test",
			},
		}
		_, err := container.ConfigurableProductService.Create(ctx, in)
		wantErr := ErrProductsMustBeAtleastOne
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_two_images_mandatory", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
				"test",
			},
			Products: []int{
				product_id,
			},
			Images: []string{
				"test",
			},
		}
		_, err := container.ConfigurableProductService.Create(ctx, in)
		wantErr := ErrImagesMustBeAtleastTwo
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Avail_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	setup()
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
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

	err = container.ConfigurableProductService.Avail(ctx, id)
	if err != nil {
		t.Errorf("Failed to avail err: %v", err)
	}

	expectedStatus := false
	got, err := container.ConfigurableProductService.Get(ctx, id)
	if err != nil {
		t.Errorf("Expected status: %v Got status: %v", expectedStatus, got.IsAvailable)
	}

}

func Test_Avail_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		err := container.ConfigurableProductService.Avail(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("already_available", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
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

		err = container.ConfigurableProductService.Avail(ctx, id)
		if err != nil {
			t.Errorf("Failed to avail err: %v", err)
		}

		err = container.ConfigurableProductService.Avail(ctx, id)
		expectedErr := ErrAlreadyAvailable
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})
}

func Test_Disable_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	setup()
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalId: "test",
		AttributeKeys: []string{
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

	err = container.ConfigurableProductService.Avail(ctx, id)
	if err != nil {
		t.Errorf("Failed to avail err: %v", err)
	}

	err = container.ConfigurableProductService.Disable(ctx, id)
	if err != nil {
		t.Fatalf("Failed to disavail err: %v", err)
	}

	got, err := container.ConfigurableProductService.Get(ctx, id)
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
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		err := container.ConfigurableProductService.Disable(ctx, 99)
		expectedErr := ErrIdNotFound
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

	t.Run("already_disabled", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalId: "test",
			AttributeKeys: []string{
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

		err = container.ConfigurableProductService.Disable(ctx, id)
		expectedErr := ErrAlreadyUnavailable
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	setup()
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

	//Get
	got, err := container.ConfigurableProductService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to Get err: %v", err)
	}

	if got.Id != id {
		t.Errorf("Expected resp Id: %v Got: %v", id, got.Id)
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		_, err := container.ConfigurableProductService.Get(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v: %v", wantErr, err)
		}
	})

	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		// Get
		wantErr := ErrIdNotFound
		_, err := container.ConfigurableProductService.Get(ctx, 99)
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

}

func Test_GetByParam_happyPath(t *testing.T) {
	t.Run("byName", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
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
		got, err := container.ConfigurableProductService.GetByParam(ctx, &GetByParamRequest{
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
		t.Cleanup(container.Teardown)
		setup()
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
		got, err := container.ConfigurableProductService.GetByParam(ctx, &GetByParamRequest{
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
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		_, err := container.ConfigurableProductService.GetByParam(ctx, &GetByParamRequest{
			Name: "test",
		})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Fatalf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_All_happyPath(t *testing.T) {
	t.Cleanup(container.Teardown)
	setup()
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
		t.Fatalf("Failed to get by param err: %v", err)
	}
	wantLen := 1
	if len(got.List) != wantLen {
		t.Errorf("Expected len: %v Got len: %v", wantLen, len(got.List))
	}
}

func Test_Get_All_unhappyPath(t *testing.T) {
	t.Run("emptyContent", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		wantErr := ErrEmptyGetContent
		_, err := container.ConfigurableProductService.GetAll(ctx)
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Update_happyPath(t *testing.T) {
	t.Run("name", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		setup()
		in := &CreateRequest{
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
		err = container.ConfigurableProductService.Update(ctx, &UpdateRequest{
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

	t.Run("desc", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		setup()
		in := &CreateRequest{
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
		err = container.ConfigurableProductService.Update(ctx, &UpdateRequest{
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

	t.Run("extId", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		setup()
		in := &CreateRequest{
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
		err = container.ConfigurableProductService.Update(ctx, &UpdateRequest{
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

	t.Run("product", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		setup()
		in := &CreateRequest{
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
		err = container.ConfigurableProductService.Update(ctx, &UpdateRequest{
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

	t.Run("isAvailableStatus", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		setup()
		in := &CreateRequest{
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

		//isAvailableStatus
		err = container.ConfigurableProductService.Update(ctx, &UpdateRequest{
			Id:                id,
			IsAvailableStatus: false,
		})
		if err != nil {
			t.Fatalf("Failed to update err: %v", err)
		}

		resp, err := container.ConfigurableProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantIsAvailableStatus := false
		if resp.IsAvailable != wantIsAvailableStatus {
			t.Errorf("Expected product: %v Got: %v", wantIsAvailableStatus, resp.IsAvailable)
		}
	})

	t.Run("images", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		ctx := context.Background()
		setup()
		in := &CreateRequest{
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
		err = container.ConfigurableProductService.Update(ctx, &UpdateRequest{
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

	// t.Run("attributes", func(t *testing.T) {
	// 	t.Cleanup(container.Teardown)
	// 	ctx := context.Background()
	// 	setup()
	// 	in := &CreateRequest{
	// 		Name:       "test",
	// 		Desc:       "test",
	// 		ExternalId: "test",
	// 		AttributeKeys: []string{
	// 			"another_test",
	// 		},
	// 		Products: []int{
	// 			product_id,
	// 		},
	// 		Images: []string{
	// 			"test",
	// 			"test",
	// 		},
	// 	}
	// 	id, err := container.ConfigurableProductService.Create(ctx, in)
	// 	if err != nil {
	// 		t.Fatalf("Failed to create err: %v", err)
	// 	}

	// 	//attribute keys
	// 	in_update := &UpdateRequest{
	// 		Id: id,
	// 		AttributeKeys: []string{
	// 			"test",
	// 		},
	// 		Product: []int{
	// 			product_id,
	// 		},
	// 	}
	// 	err = container.ConfigurableProductService.Update(ctx, in_update)
	// 	if err != nil {
	// 		t.Fatalf("Failed to update err: %v", err)
	// 	}

	// 	resp, err := container.ConfigurableProductService.Get(ctx, id)
	// 	if err != nil {
	// 		t.Fatalf("Failed to get err: %v", err)
	// 	}
	// 	want := map[string]string{
	// 		"test": "test",
	// 	}
	// 	for got_key, _ := range resp.Attributes[0] {
	// 		for want_key, _ := range want {
	// 			if got_key != want_key {
	// 				t.Errorf("Expected key %v Got: %v", want_key, got_key)
	// 			}
	// 		}
	// 	}
	// })
}

func Test_Update_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		t.Cleanup(container.Teardown)
		setup()
		ctx := context.Background()
		err := container.ConfigurableProductService.Update(ctx, &UpdateRequest{
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
