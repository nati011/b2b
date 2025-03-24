package user

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/retailer"
)

var testContainer TestContainer
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
	_, err := testContainer.RetailerService.Create(ctx, &in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}

	//check if user exists
	//attempt to deactivate
	testContainer.UserService.Activate(ctx, id)

}
