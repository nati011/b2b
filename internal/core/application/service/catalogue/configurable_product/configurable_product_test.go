package configurable_product

import (
	"context"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/catalogue/configurable_product"
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

func Test_Get(t *testing.T) {
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

func Test_GetByParam(t *testing.T) {

}

func Test_GetAll(t *testing.T) {

}
