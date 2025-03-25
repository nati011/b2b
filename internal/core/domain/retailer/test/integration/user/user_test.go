package user

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/retailer"
)

var testContainer retailer.TestContainer
var ctx context.Context

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	ctx = context.Background()
	testContainer = retailer.NewPackageIntegrationTestContainer()
}

func Test_create_default_admin_user_upon_retailer_registration(t *testing.T) {
	//setup
	t.Cleanup(testContainer.Cleanup)
	in := retailer.CreateRequest{
		Tin:         "1111111111",
		Latitude:    "9.0192° N",
		Longitude:   "38.7525° E",
		GeneralZone: "test",
		Region:      "test",
		Woreda:      "test",

		FirstName: "test",
		LastName:  "test",
		Email:     "test@gmail.com",
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
