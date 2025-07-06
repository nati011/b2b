package user

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/retailer"
	test_container "b2b.nati011.github.com/internal/core/domain/retailer/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var testContainer test_container.TestContainer
var db *sql.DB
var ctx context.Context

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	ctx = context.Background()
	testContainer = test_container.NewDBIntegrationTestContainer(db)
}

func teardown() {
	db_test_container.Teardown(db)
	testContainer.Teardown(db)
}

func Test_create_default_admin_user_upon_retailer_registration(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "test",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@gmail.com",
		Phone:       "+251949184879",
	}
	id, err := testContainer.RetailerService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//get all retailer users
	users, err := testContainer.RetailerService.GetAllUsers(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get users err: %v", err)
	}
	//attempt to deactivate
	err = testContainer.UserService.Deactivate(ctx, users.List[0])
	if err != nil {
		t.Fatalf("Failed to activate user err: %v", err)
	}
}

func Test_assign_default_admin_user_role_retailer(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	in := retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "test",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@gmail.com",
		Phone:       "+251949184879",
	}
	id, err := testContainer.RetailerService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	users, err := testContainer.RetailerService.GetAllUsers(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get all users err: %v", err)
	}
}

func Test_Get_All_Users_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	t.Cleanup(teardown)
	in := retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "test",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@gmail.com",
		Phone:       "+251949184879",
	}
	id, err := testContainer.RetailerService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	_, err = testContainer.RetailerService.GetAllUsers(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get user agents %v", err)
	}
}

func Test_Get_All_Users_unhappyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	t.Cleanup(teardown)
	in := retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",
		Username:    "test",
		FirstName:   "test",
		LastName:    "test",
		Email:       "test@gmail.com",
		Phone:       "+251949184879",
	}
	id, err := testContainer.RetailerService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	retailer_users, err := testContainer.RetailerService.GetAllUsers(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get user agents %v", err)
	}

	//remove retailer
	err = testContainer.UserService.Remove(ctx, retailer_users.List[0])
	if err != nil {
		t.Fatalf("Failed to remove user agent with id: %v  err: %v", retailer_users.List[0], err)
	}

	resp, err := testContainer.RetailerService.GetAllUsers(ctx, id)
	print(resp.List)
	wantErr := retailer.ErrRetailerHasNoUsers
	if err != wantErr {
		t.Errorf("Expected err: %v got: %v", wantErr, err)
	}
}
