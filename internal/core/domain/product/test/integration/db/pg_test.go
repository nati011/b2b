package db

import (
	"context"
	"database/sql"
	"math"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/category"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/product"
	test_container "b2b.nati011.github.com/internal/core/domain/product/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var container test_container.TestContainer
var db *sql.DB
var categoryId int
var distributorId int

func TestMain(m *testing.M) {

	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	container = test_container.NewDBIntegrationTestContainer(
		db,
	)
	ctx := context.Background()
	categoryId, _ = container.CategoryService.Create(ctx, &category.CreateRequest{
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
		panic("failed to create distributor err: ")
	}
}

func teardown() {
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {
}

//for performace reasons run tests separately

func Test_Get(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
			"tet":  "tets",
		},
		DistributorId: distributorId,
	}

	id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	got, err := container.ProductService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get product err %v", err)
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

func Test_GetAll(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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

	in_two := &product.CreateRequest{
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

func Test_GetAllStockLedger(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	prod_id, err := container.ProductService.Create(ctx, &product.CreateRequest{
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
	})
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	err = container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     prod_id,
		Amount: 10,
	})
	if err != nil {
		t.Fatalf("Failed to receive goods")
	}

	resp, err := container.ProductService.GetAllStockLedger(ctx)
	if err != nil {
		t.Errorf("Expected err:%v Got err: %v", nil, err)
	}

	if resp.List[0].Product_id != prod_id {
		t.Errorf("expected productId: %v Got: %v", resp.List[0].Product_id, prod_id)
	}

	if resp.List[0].Operation != "GOODS_RECEIVING" {
		t.Errorf("expected productId: %v Got: %v", resp.List[0].Product_id, prod_id)
	}
}

func Test_GetStockLedger(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	prod_id, err := container.ProductService.Create(ctx, &product.CreateRequest{
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
	})
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	err = container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     prod_id,
		Amount: 10,
	})
	if err != nil {
		t.Fatalf("Failed to receive goods")
	}

	resp, err := container.ProductService.GetStockLedger(ctx, prod_id)
	if err != nil {
		t.Errorf("Expected err:%v Got err: %v", nil, err)
	}

	if resp.List[0].Product_id != prod_id {
		t.Errorf("expected productId: %v Got: %v", resp.List[0].Product_id, prod_id)
	}

	if resp.List[0].Operation != "GOODS_RECEIVING" {
		t.Errorf("expected productId: %v Got: %v", resp.List[0].Product_id, prod_id)
	}
}

func Test_GetByName(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_GetByExternalId(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_GetByDistributorId(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
		t.Fatalf("Failed to create product err: %v", err)
	}

	got, err := container.ProductService.GetByParam(ctx, &product.GetByParamRequest{
		DistributorId: distributorId,
	})
	if err != nil {
		t.Fatalf("Failed to get product err: %v", err)
	}
	wantLen := 1
	if wantLen != len(got.List) {
		t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
	}
}

func Test_GetByCategory(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_GetByPriceRange(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_search(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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

	got, err := container.ProductService.Search(ctx, &product.SearchRequest{
		Name: "sear",
	})
	if err != nil {
		t.Errorf("Expected err: %v, Got err: %v", nil, err)
	}
	wantLen := 1
	if wantLen != len(got.List) {
		t.Errorf("Expected len: %v Got err: %v", wantLen, len(got.List))
	}
}

//for performance reasons all tests are separately run

func Test_create(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
			"tet":  "tst",
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
		t.Fatalf("Failed to get product err %v", err)
	}

	if resp.Name != in.Name {
		t.Errorf("Expected name:%v Got: %v", resp.Name, in.Name)
	}
	if resp.Desc != in.Desc {
		t.Errorf("Expected desc:%v Got: %v", resp.Desc, in.Desc)
	}
	if resp.ExternalID != in.ExternalID {
		t.Errorf("Expected extId:%v Got: %v", resp.ExternalID, in.ExternalID)
	}
	if resp.Price != in.Price {
		t.Errorf("Expected price:%v Got: %v", resp.Price, in.Price)
	}
	for i, v := range in.Attributes {
		if resp.Attributes[i] != v {
			t.Errorf("Expected attr: %v Got: %v", v, resp.Attributes[i])
		}
	}

	for i, v := range in.Images {
		if resp.Images[i].ImageUrl != v {
			t.Errorf("Expected attr:%v Got: %v", v, resp.Images[i])
		}
	}
}

func Test_updateName(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_updateExternalId(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_updatePrice(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_updateDesc(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}
func Test_updateImages(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
	update_in := &product.UpdateRequest{
		Id: id,
		Images: []string{
			"https://res.cloudinary.com/ddbdbuuqw/image/upload/v1713311534/kecw097ntniwoiub04sz.png",
			"https://res.cloudinary.com/ddbdbuuqw/image/upload/v1713311534/kecw097ntniwoiub04sz.png",
			"https://res.cloudinary.com/ddbdbuuqw/image/upload/v1713311534/kecw097ntniwoiub04sz.png",
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
}

func Test_updateActiveStatus(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_updateCategory(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_goodsReceiving(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_dispatch(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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
}

func Test_reserve(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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

	product_id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	goodsReceivingAmount := 5
	err = container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     product_id,
		Amount: goodsReceivingAmount,
	})
	if err != nil {
		t.Fatalf("Failed to update err: %v", err)
	}

	reservedAmount := 4
	err = container.ProductService.Reserve(ctx, product_id, reservedAmount)
	if err != nil {
		t.Fatalf("Failed to update err: %v", err)
	}

	got, err := container.ProductService.Get(ctx, product_id)
	if err != nil {
		t.Fatalf("Failed to Get")
	}

	if got.AvailableStock != goodsReceivingAmount-reservedAmount {
		t.Fatalf("expected available stock: %v Got: %v", goodsReceivingAmount-reservedAmount, got.AvailableStock)
	}
}

func Test_freeReservedStock(t *testing.T) {
	t.Cleanup(teardown)
	setup()
	ctx := context.Background()
	in := &product.CreateRequest{
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

	product_id, err := container.ProductService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create product")
	}

	goodsReceivingAmount := 5
	err = container.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     product_id,
		Amount: goodsReceivingAmount,
	})
	if err != nil {
		t.Fatalf("Failed to update err: %v", err)
	}

	reservedAmount := 4
	err = container.ProductService.Reserve(ctx, product_id, reservedAmount)
	if err != nil {
		t.Fatalf("Failed to update err: %v", err)
	}

	err = container.ProductService.FreeReservation(ctx, product_id, reservedAmount)
	if err != nil {
		t.Fatalf("Failed to update err: %v", err)
	}

	got, err := container.ProductService.Get(ctx, product_id)
	if err != nil {
		t.Fatalf("Failed to Get")
	}

	if got.AvailableStock != goodsReceivingAmount {
		t.Fatalf("expected available stock: %v Got: %v", goodsReceivingAmount, got.AvailableStock)
	}
}
