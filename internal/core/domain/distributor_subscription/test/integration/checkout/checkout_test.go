package integration

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/distributor_subscription"
	test "b2b.nati011.github.com/internal/core/domain/distributor_subscription/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var container test.TestContainer
var db *sql.DB
var distributorId = 1
var DigitalPaymentPartnerId int
var ManualPaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	container = test.NewIntegrationTestContainer(
		db,
	)
	ctx := context.Background()
	var err error
	ManualPaymentPartnerId, _ = container.PartnerService.Create(ctx, &payment_partner.CreateRequest{
		Name:          "chapa",
		Icon:          "etst",
		BaseURL:       "https://api.chapa.co",
		Secret:        "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
		PaymentMethod: payment_partner.PAYMENT_METHOD_MANUAL,
	})
	DigitalPaymentPartnerId, err = container.PartnerService.Create(ctx,
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
	// testContainer.Teardown(db)
	db_test_container.Teardown(db)
}

func Test_Checkout_happypath(t *testing.T) {
	t.Run("generate_checkout_url_upon_placement_digital", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := container.SubscriptionService.CreatePlan(ctx, &distributor_subscription.CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
		subs_resp, err := container.SubscriptionService.Place(ctx, &distributor_subscription.PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   DigitalPaymentPartnerId,
		})
		if err != nil {
			t.Fatalf("failed to place order err: %v", err)
		}
		if strings.Split(subs_resp.CheckoutUrl, " ") == nil {
			t.Errorf("Expected checkoutUrl different from nil got: %v", subs_resp.CheckoutUrl)
		}
	})

	t.Run("checkout_upon_placement_manual", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := container.SubscriptionService.CreatePlan(ctx, &distributor_subscription.CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
		_, err = container.SubscriptionService.Place(ctx, &distributor_subscription.PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   ManualPaymentPartnerId,
		})
		if err != nil {
			t.Fatalf("failed to place order err: %v", err)
		}
	})

	t.Run("generate_checkout_url_upon_init_payment", func(t *testing.T) {
		t.Cleanup(teardown)
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := container.SubscriptionService.CreatePlan(ctx, &distributor_subscription.CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
		subs_resp, err := container.SubscriptionService.Place(ctx, &distributor_subscription.PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   ManualPaymentPartnerId,
		})
		if err != nil {
			t.Fatalf("failed to place order err: %v", err)
		}
		if strings.Split(subs_resp.CheckoutUrl, " ") == nil {
			t.Errorf("Expected checkoutUrl different from nil got: %v", subs_resp.CheckoutUrl)
		}
		init_resp, err := container.SubscriptionService.InitPayment(ctx, subs_resp.Id)
		if err != nil {
			t.Fatalf("Failed to init payment err: %v", err)
		}
		if strings.Split(init_resp.CheckoutUrl, " ") == nil {
			t.Errorf("Expected checkoutUrl different from nil got: %v", init_resp.CheckoutUrl)
		}
	})
}
