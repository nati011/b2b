package pg

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/domain/distributor_subscription"
	test "b2b.nati011.github.com/internal/core/domain/distributor_subscription/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var testContainer test.TestContainer
var db *sql.DB
var distributorId = 1
var DigitalPaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	testContainer = test.NewIntegrationTestContainer(
		db,
	)
	ctx := context.Background()
	var err error
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
	// testContainer.Teardown(db)
	db_test_container.Teardown(db)
}

func Test_Read(t *testing.T) {

	t.Run("GetAllPlan", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := testContainer.SubscriptionService.CreatePlan(ctx, &distributor_subscription.CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}

		got, err := testContainer.SubscriptionService.GetAllPlan(ctx)
		if err != nil {
			t.Fatalf("failed to get plans err: %v", err)
		}
		wantLen := 1
		if len(got.List) != wantLen {
			t.Fatalf("expected len: %v got: %v", wantLen, len(got.List))
		}

	})

	t.Run("GetPlan", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := testContainer.SubscriptionService.CreatePlan(ctx, &distributor_subscription.CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}

		got, err := testContainer.SubscriptionService.GetPlan(ctx, planId)
		if err != nil {
			t.Fatalf("failed to get plans err: %v", err)
		}
		wantId := planId
		if wantId != got.Id {
			t.Fatalf("expected id: %v got: %v", wantId, got.Id)
		}
	})

	t.Run("GetAllSubscriptions", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := testContainer.SubscriptionService.CreatePlan(ctx, &distributor_subscription.CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
		_, err = testContainer.SubscriptionService.Place(ctx, &distributor_subscription.PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   DigitalPaymentPartnerId,
		})
		if err != nil {
			t.Fatalf("failed to place order err: %v", err)
		}
		all_resp, err := testContainer.SubscriptionService.GetSubscriptions(ctx)
		if err != nil {
			t.Fatalf("failed to get plans err: %v", err)
		}
		wantLen := 1
		if len(all_resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(all_resp.List))
		}
	})

	t.Run("GetSubscriptionsByDistributorId", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := testContainer.SubscriptionService.CreatePlan(ctx, &distributor_subscription.CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
		sub_id, err := testContainer.SubscriptionService.Place(ctx, &distributor_subscription.PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   DigitalPaymentPartnerId,
		})
		if err != nil {
			t.Fatalf("failed to place order err: %v", err)
		}
		got, err := testContainer.SubscriptionService.GetSubscriptionByDistributorId(ctx, distributorId)
		if err != nil {
			t.Fatalf("failed to get subscriptions by distId err: %v", err)
		}
		wantSubId := sub_id.Id
		if wantSubId != got.Id {
			t.Errorf("expected sub Id: %v got :%v", wantSubId, got.Id)
		}
	})
}

func Test_Write(t *testing.T) {
	t.Run("createPlan", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := testContainer.SubscriptionService.CreatePlan(ctx, &distributor_subscription.CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
	})

	t.Run("place", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := testContainer.SubscriptionService.CreatePlan(ctx, &distributor_subscription.CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
		_, err = testContainer.SubscriptionService.Place(ctx, &distributor_subscription.PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   DigitalPaymentPartnerId,
		})
		if err != nil {
			t.Fatalf("failed to place order err: %v", err)
		}
	})

}
