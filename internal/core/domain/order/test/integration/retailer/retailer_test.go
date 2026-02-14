package retailer

import (
	"context"
	"os"
	"testing"

	"math/rand"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	"b2b.nati011.github.com/internal/core/domain/order"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
)

var testContainer order.TestContainer
var retailerId int
var retailerUserId int
var productId int
var distributorId int
var DigitalPaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = order.NewPackageIntegrationTestContainer()
	ctx := context.Background()
	var err error
	retailerId, err = testContainer.RetailerService.Create(ctx, &retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "order_test",
		FirstName:   "test",
		LastName:    "test",
		Phone:       "+251949184879",
		Email:       "test@gmail.com",
	})
	if err != nil {
		panic("failed to create product")
	}

	restailerUsers, err := testContainer.RetailerService.GetAllUsers(ctx, retailerId)
	if err != nil {
		panic("failed to get all retailer users")
	}
	retailerUserId = restailerUsers.List[0]

	distributorId, err = testContainer.DistributorService.Create(ctx, &distributor.CreateRequest{
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

	err = testContainer.DistributorService.Activate(ctx, distributorId)
	if err != nil {
		panic("failed to activate distributor")
	}

	productId, err = testContainer.ProductService.Create(ctx, &product.CreateRequest{
		Name:       "testProduct",
		Desc:       "test",
		ExternalID: "123",
		Images: []string{
			"https://picsum.photos/200",
			"https://picsum.photos/200",
		},
		Price: 100.00,
		Attributes: map[string]string{
			"test": "test",
		},
		DistributorId: distributorId,
	})
	if err != nil {
		panic("failed to create product")
	}
	err = testContainer.ProductService.ReceiveGoods(ctx, &product.GoodsReceivingRequest{
		Id:     productId,
		Amount: 10000,
	})
	if err != nil {
		panic("failed to recieve goods")
	}
	DigitalPaymentPartnerId, err = testContainer.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:          "chapa",
			Icon:          "etst",
			BaseURL:       "https://api.chapa.co",
			Secret:        "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
			PaymentMethod: payment_partner.PAYMENT_METHOD_DIGITAL,
		})
	if err != nil {
		panic("failed to create payment partner")
	}
}

func teardown() {
	testContainer.Teardown()
}

func Test_validate_retailer_upon_order_placement(t *testing.T) {
	ctx := context.Background()
	randomRetailerId := rand.Int()
	in := &order.PlaceRequest{
		RetailerId: randomRetailerId,
		Items: []order.Item{
			{
				ProductId: productId,
				Quantity:  1},
		},
	}
	_, err := testContainer.OrderService.Place(ctx, in)
	wantErr := order.ErrRetailerIdNotFound
	if err != wantErr {
		t.Errorf("Expected err : %v Got: %v", wantErr, err)
	}
}

func Test_Get_By_UserId_happyPath(t *testing.T) {
	t.Run("getByUserId", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()

		got, err := testContainer.OrderService.GetRetailerOrdersWithUserContext(ctx, retailerUserId)
		if err != nil {
			switch err {
			case order.ErrEmptyGetResponse:
			default:
				t.Fatalf("Failed to fetch order err: err %v", err)
			}
		}
		wantLen := 0
		if len(got.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(got.List))
		}
	})
}

func Test_Get_By_UserId_unhappyPath(t *testing.T) {
	t.Run("getByUserId", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//check
		randomRetailerUserId := rand.Int()
		_, err := testContainer.OrderService.GetRetailerOrdersWithUserContext(ctx, randomRetailerUserId)
		wantErr := order.ErrRetailerIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_Place_By_UserId_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := &order.PlaceAsUserRequest{
		PaymentPartnerId: DigitalPaymentPartnerId,
		UserId:           retailerUserId,
		Items: []order.Item{
			{
				ProductId: productId,
				Quantity:  1},
		},
	}
	order_resp, err := testContainer.OrderService.PlaceWithUserContext(ctx, in)
	if err != nil {
		t.Fatalf("Failed to fetch order err: err %v", err)
	}

	//check
	resp, err := testContainer.OrderService.Get(ctx, order_resp.Id)
	if err != nil {
		t.Fatalf("Failed to fetch order err: err %v", err)
	}
	if resp.Id != order_resp.Id {
		t.Errorf("Expected Id: %v Got Id: %v", order_resp.Id, resp.Id)
	}
}

func Test_Place_By_UserId_unhappyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	randomRetailerUserId := rand.Int()
	in := &order.PlaceAsUserRequest{
		PaymentPartnerId: DigitalPaymentPartnerId,
		UserId:           randomRetailerUserId,
		Items: []order.Item{
			{
				ProductId: productId,
				Quantity:  1},
		},
	}
	_, err := testContainer.OrderService.PlaceWithUserContext(ctx, in)
	wantErr := order.ErrRetailerIdNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v Got: %v", wantErr, err)
	}
}
