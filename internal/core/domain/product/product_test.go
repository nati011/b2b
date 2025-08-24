package product

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

var container TestContainer
var category_id int
var distributorId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = NewPackageIntegrationTestContainer()
	ctx := context.Background()
	category_id, _ = container.CategoryService.Create(ctx, &category.CreateRequest{
		Name: "test",
	})
	var err error
	distributorId, err = container.DistributorService.Create(ctx, &distributor.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "username",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@gmail.com",
	})
	if err != nil {
		panic("failed to create distributor")
	}
	err = container.DistributorService.Activate(ctx, distributorId)
	if err != nil {
		panic("failed to activate distributor")
	}
}

func tearDown() {
	container.Teardown()
}

func Test_Create_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "testProduct",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
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

	t.Run("inactive_by_default", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Error("Failed to create product", err)
		}

		got, err := container.ProductService.Get(ctx, id)
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
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}
		_, err := container.ProductService.Create(ctx, in)
		wantErr := ErrNameNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("name_duplicate", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}
		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		in_new := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}
		_, err = container.ProductService.Create(ctx, in_new)
		wantErr := ErrNameDuplicate
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("desc_mandatory", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "t",
			ExternalID: "1",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}
		_, err := container.ProductService.Create(ctx, in)
		wantErr := ErrDescNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_two_images_mandatory", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Two ImageTest",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}
		_, err := container.ProductService.Create(ctx, in)
		wantErr := ErrImagesMustBeAtleastTwo
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("attribute_values_cannot_be_empty", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Value EmptyTest",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "",
			},
			DistributorId: distributorId,
		}
		_, err := container.ProductService.Create(ctx, in)
		wantErr := ErrAttributeValuesCannotBeEmpty
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("price_mandatory", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Price Test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}
		_, err := container.ProductService.Create(ctx, in)
		wantErr := ErrPriceNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("price_cannot_be_zero", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Price ZeroTest",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 0.00,
			Attributes: map[string]string{
				"test": "tets",
			},
			DistributorId: distributorId,
		}
		_, err := container.ProductService.Create(ctx, in)
		wantErr := ErrPriceNotSupplied
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
			"tets": "test",
		},
		DistributorId: distributorId,
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	//get
	got, err := container.ProductService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err %v", err)
	}

	if got.Name != in.Name {
		t.Errorf("Expected name:%v Got: %v", got.Name, in.Name)
	}
	if got.Desc != in.Desc {
		t.Errorf("Expected desc:%v Got: %v", got.Desc, in.Desc)
	}
	if got.ExternalID != in.ExternalID {
		t.Errorf("Expected extId:%v Got: %v", got.ExternalID, in.ExternalID)
	}
	if got.Price != in.Price {
		t.Errorf("Expected price:%v Got: %v", got.Price, in.Price)
	}

	for i, v := range in.Attributes {
		if got.Attributes[i] != v {
			t.Errorf("Expected attr: %v Got: %v", v, got.Attributes[i])
		}
	}

	for i, v := range in.Images {
		if got.Images[i].ImageUrl != v {
			t.Errorf("Expected attr:%v Got: %v", v, got.Images[i])
		}
	}
}

func Test_Get_unhappyPath(t *testing.T) {
	t.Run("no_product_found", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		//get
		wantErr := ErrIdNotFound
		_, err := container.ProductService.Get(ctx, 99)
		if err != wantErr {
			t.Errorf("Expected err:%v Got err: %v", wantErr, err)
		}
	})
}

func Test_Get_All_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
		DistributorId: distributorId,
	}

	_, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	in_two := &CreateRequest{
		Name:       "test2",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
		DistributorId: distributorId,
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
}

func Test_Get_All_unhappyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	//get-all
	wantErr := ErrEmptyGetContent
	_, err := container.ProductService.GetAll(ctx)
	if err != wantErr {
		t.Errorf("Expected err:%v Got err: %v", ErrEmptyGetContent, err)
	}
}

func Test_Search_happyPath(t *testing.T) {
	t.Run("search", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "searchTest",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := container.ProductService.Search(ctx, "sear")
		if err != nil {
			t.Errorf("Expected err: %v, Got err: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got err: %v", wantLen, len(got.List))
		}
	})
}

func Test_Get_by_param_happyPath(t *testing.T) {

	t.Run("byName", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := container.ProductService.GetByParam(ctx, &GetByParamRequest{
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
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "External Id",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := container.ProductService.GetByParam(ctx, &GetByParamRequest{
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

	t.Run("byDistributorId", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := container.ProductService.GetByParam(ctx, &GetByParamRequest{
			DistributorId: distributorId,
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})

	t.Run("byDistributorId", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 500.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := container.ProductService.GetByParam(ctx, &GetByParamRequest{
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

	t.Run("byCategory", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		category_id, _ = container.CategoryService.Create(ctx, &category.CreateRequest{
			Name: "test",
		})

		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 500.00,
			Attributes: map[string]string{
				"test": "test",
			},
			CategoryId:    []int{category_id},
			DistributorId: distributorId,
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}

		got, err := container.ProductService.GetByParam(ctx, &GetByParamRequest{
			CategoryId: []int{category_id},
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})

	t.Run("aggregate-Fetch", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Aggregate Fetch",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 900.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		in_2 := &CreateRequest{
			Name:       "test1",
			Desc:       "test1",
			ExternalID: "1233",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 950.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		_, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		_, err = container.ProductService.Create(ctx, in_2)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := container.ProductService.GetByParam(ctx, &GetByParamRequest{
			PriceMin: 100,
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
		t.Cleanup(tearDown)
		ctx := context.Background()
		_, err := container.ProductService.GetByParam(ctx, &GetByParamRequest{})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})
}

func Test_Update_happyPath(t *testing.T) {
	t.Run("update", func(t *testing.T) {
		t.Cleanup(tearDown)
		//setup
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
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
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
		}
		_, err = container.ProductService.Update(ctx, update_in)
		if err != nil {
			t.Fatalf("Failed to update")
		}

		got, err := container.ProductService.Get(ctx, id)
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
		t.Cleanup(tearDown)
		ctx := context.Background()
		//update
		update_in := &UpdateRequest{
			Id:         99,
			Name:       "test1",
			ExternalID: "321",
			Desc:       "1",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 200,
		}
		wantErr := ErrIdNotFound
		_, err := container.ProductService.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("atleast_two_images_mandatory", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Two Images Test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
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
		_, err = container.ProductService.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("duplicate_name_not_allowed", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "Duplicate Name tEST",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		in_2 := &CreateRequest{
			Name:       "duplicate",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		_, err = container.ProductService.Create(ctx, in_2)
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
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 200,
		}
		wantErr := ErrNameDuplicate
		_, err = container.ProductService.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("price_cannot_be_zero", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "zERO PRICE",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
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
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: -1,
		}
		wantErr := ErrPriceCannotBeNegative
		_, err = container.ProductService.Update(ctx, update_in)
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})
}

func Test_goods_receiving_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
		DistributorId: distributorId,
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
		Id:     id,
		Amount: 2,
	})
	if err != nil {
		t.Fatalf("Failed to perform goods receiving err: %v", err)
	}

	//check stock
	got, err := container.ProductService.Get(ctx, id)
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
		t.Cleanup(tearDown)
		ctx := context.Background()
		wantErr := ErrIdNotFound
		err := container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     99,
			Amount: 2,
		})
		if err != wantErr {
			t.Errorf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})
}

func Test_dispatch_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
		DistributorId: distributorId,
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
		Id:     id,
		Amount: 2,
	})
	if err != nil {
		t.Fatalf("Failed to initiate goods receiving")
	}

	err = container.ProductService.Dispatch(ctx, &DispatchRequest{
		Id:     id,
		Amount: 2,
	})
	if err != nil {
		t.Fatalf("Failed to initiate dispatch")
	}

	// check stock
	got, err := container.ProductService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get product")
	}

	wantStock := 0
	if got.Stock != wantStock {
		t.Errorf("Expected stock: %v Got stock: %v", wantStock, got.Stock)
	}
}

func Test_dispatch_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		err := container.ProductService.Dispatch(ctx, &DispatchRequest{
			Id:     99,
			Amount: 2,
		})
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})

	t.Run("requested_quantity_greater_than_stock", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     id,
			Amount: 2,
		})
		if err != nil {
			t.Fatalf("Failed to initiate goods receiving")
		}

		err = container.ProductService.Dispatch(ctx, &DispatchRequest{
			Id:     id,
			Amount: 3,
		})
		wantErr := ErrStockUnavailable
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_Activate_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
		DistributorId: distributorId,
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	err = container.ProductService.Activate(ctx, id)
	if err != nil {
		t.Fatalf("Failed to activate product")
	}

	//verify status
	wantActiveStatus := true
	resp, err := container.ProductService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to activte status")
	}
	if resp.IsActive != wantActiveStatus {
		t.Errorf("Expected active status: %v, Want status: %v", wantActiveStatus, resp.IsActive)
	}
}

func Test_Activate_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		err := container.ProductService.Activate(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Fatalf("Expected err: %v, Got err: %v", wantErr, err)
		}
	})

	t.Run("already_active", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		err = container.ProductService.Activate(ctx, id)
		if err != nil {
			t.Fatalf("Failed to activate product")
		}

		err = container.ProductService.Activate(ctx, id)
		if err != ErrAlreadyActive {
			t.Fatalf("Failed to activate product")
		}
	})
}

func Test_Deactvate_happyPath(t *testing.T) {
	t.Cleanup(tearDown)
	ctx := context.Background()
	in := &CreateRequest{
		Name:       "test",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
		DistributorId: distributorId,
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	//verify
	got, err := container.ProductService.Get(ctx, id)
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
		t.Cleanup(tearDown)
		ctx := context.Background()
		err := container.ProductService.Activate(ctx, 99)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Want err: %v", wantErr, err)
		}
	})

	t.Run("already_inactive", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		err = container.ProductService.Deactivate(ctx, id)
		expectedErr := ErrAlreadyInactive
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})
}

func Test_Reserve_happyPath(t *testing.T) {
	t.Run("deduct_available_item_upon_reservation", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}
		inital_stock := 100
		err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     id,
			Amount: inital_stock,
		})
		if err != nil {
			t.Fatalf("Failed to recieveGoods err: %v", err)
		}

		reserved_stock := 50
		err = container.ProductService.Reserve(
			ctx,
			id,
			reserved_stock)
		if err != nil {
			t.Fatalf("Failed to reserve stock err: %v", err)
		}

		product, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get product err: %v", err)
		}

		if product.AvailableStock != inital_stock-reserved_stock {
			t.Errorf("Want availableStock: %v, Got: %v", inital_stock-reserved_stock, product.AvailableStock)
		}
	})

	t.Run("add_reserved_item_upon_reservation", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}
		inital_stock := 100
		err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     id,
			Amount: inital_stock,
		})
		if err != nil {
			t.Fatalf("Failed to recieveGoods err: %v", err)
		}

		reserved_stock := 50
		err = container.ProductService.Reserve(ctx, id, reserved_stock)
		if err != nil {
			t.Fatalf("Failed to reserve stock err: %v", err)
		}

		product, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get product err: %v", err)
		}

		if product.ReservedStock != reserved_stock {
			t.Errorf("Want reservedStock: %v, Got: %v", reserved_stock, product.AvailableStock)
		}
	})
}

func Test_Reserve_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}
		inital_stock := 100
		err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     id,
			Amount: inital_stock,
		})
		if err != nil {
			t.Fatalf("Failed to recieveGoods err: %v", err)
		}

		reserved_stock := 50
		err = container.ProductService.Reserve(ctx, 9999, reserved_stock)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("qty_not_found", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}
		inital_stock := 100
		err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     id,
			Amount: inital_stock,
		})
		if err != nil {
			t.Fatalf("Failed to recieveGoods err: %v", err)
		}

		err = container.ProductService.Reserve(ctx, id, 100000)
		wantErr := ErrStockReservationQtyMustBeLessThanOrEqualToAvailableQty
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_Free_Reservation_happyPath(t *testing.T) {
	t.Run("add_available_item_upon_reservation", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}
		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}
		inital_stock := 100
		err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     id,
			Amount: inital_stock,
		})
		if err != nil {
			t.Fatalf("Failed to recieveGoods err: %v", err)
		}

		reserved_stock := 50
		err = container.ProductService.Reserve(ctx, id, reserved_stock)
		if err != nil {
			t.Fatalf("Failed to reserve stock err: %v", err)
		}

		err = container.ProductService.FreeReservation(ctx, id, reserved_stock)
		if err != nil {
			t.Fatalf("Failed to reserve stock err: %v", err)
		}

		product, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get product err: %v", err)
		}

		if product.AvailableStock != inital_stock {
			t.Errorf("Want availableStock: %v, Got: %v", product.AvailableStock, inital_stock)
		}
	})

	t.Run("deduct_reserved_item_upon_reservation", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}
		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}
		inital_stock := 100
		err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     id,
			Amount: inital_stock,
		})
		if err != nil {
			t.Fatalf("Failed to recieveGoods err: %v", err)
		}

		reserved_stock := 50
		err = container.ProductService.Reserve(ctx, id, reserved_stock)
		if err != nil {
			t.Fatalf("Failed to reserve stock err: %v", err)
		}

		err = container.ProductService.FreeReservation(ctx, id, reserved_stock)
		if err != nil {
			t.Fatalf("Failed to reserve stock err: %v", err)
		}

		product, err := container.ProductService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get product err: %v", err)
		}

		if product.ReservedStock != 0 {
			t.Errorf("Want reservedStock: %v, Got: %v", product.AvailableStock, 0)
		}
	})
}

func Test_Free_Reservation_unhappyPath(t *testing.T) {
	t.Run("id_not_found", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}
		inital_stock := 100
		err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     id,
			Amount: inital_stock,
		})
		if err != nil {
			t.Fatalf("Failed to recieveGoods err: %v", err)
		}

		reserved_stock := 50
		err = container.ProductService.Reserve(ctx, id, reserved_stock)
		if err != nil {
			t.Fatalf("Failed to Reserve err: %v", err)
		}

		err = container.ProductService.FreeReservation(ctx, 9999, reserved_stock)
		wantErr := ErrIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("qty_more_than_reservation", func(t *testing.T) {
		t.Cleanup(tearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w500&q80",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: distributorId,
		}

		id, err := container.ProductService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product err: %v", err)
		}
		inital_stock := 100
		err = container.ProductService.ReceiveGoods(ctx, &GoodsReceivingRequest{
			Id:     id,
			Amount: inital_stock,
		})
		if err != nil {
			t.Fatalf("Failed to recieveGoods err: %v", err)
		}

		err = container.ProductService.FreeReservation(ctx, id, 100000)
		wantErr := ErrFreeReservationQtyMustBeLessThanOrEqualToReservedQty
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}
