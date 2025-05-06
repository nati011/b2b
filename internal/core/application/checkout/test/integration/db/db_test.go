package transaction

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/application/checkout/test"
	"b2b.nati011.github.com/internal/core/application/payment_partner"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var testContainer test.TestContainer
var db *sql.DB
var PaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	testContainer = test.NewDBIntegrationTestContainer(db)
	ctx := context.Background()
	var err error
	PaymentPartnerId, err = testContainer.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:    "chapa",
			Icon:    "test",
			BaseURL: "https://api.chapa.co",
			Secret:  "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		})

	if err != nil {
		panic("Failed to create payment partner")
	}

}

func Test_read(t *testing.T) {
}

func Test_write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		setup()
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		in := &checkout.CheckoutRequest{
			PaymentPartnerId: PaymentPartnerId,
			OrderId:          1,
			Amount:           100,
		}

		_, err := testContainer.CheckoutService.Checkout(ctx, in)
		if err != nil {
			t.Errorf("Failed to checkout err: %v", err)
		}

	})
}
