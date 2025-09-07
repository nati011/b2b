package distributor_subscription

import (
	"context"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
)

var container TestContainer
var distributorId = 1
var DigitalPaymentPartnerId int
var ManualPaymentPartnerId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = NewTestContainer()
	ctx := context.Background()
	var err error
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

	ManualPaymentPartnerId, err = container.PartnerService.Create(ctx,
		&payment_partner.CreateRequest{
			Name:          "chapa",
			Icon:          "etst",
			BaseURL:       "https://api.chapa.co",
			Secret:        "CHASECK_TEST-KUZmLnnAPtwFg8hQPqCx7mc4o7TUbIe5",
			PaymentMethod: payment_partner.PAYMENT_METHOD_MANUAL,
		})
	if err != nil {
		panic("failed to create payment partner")
	}
}

func teardown() {
	container.Teardown()
}

func Test_Create_Happypath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	_, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
		Name:        "test",
		Price:       101,
		TermInMonth: 11,
		Description: "test",
	})
	if err != nil {
		t.Fatalf("failed to create plan err: %v", err)
	}
}

func Test_Create_unappypath(t *testing.T) {
	t.Run("name_mandatory", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
			Name:        "",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		wantErr := ErrNameMandatory
		if err != wantErr {
			t.Fatalf("expected err: %v got: %v", wantErr, err)
		}
	})

	t.Run("desc_mandatory", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "",
		})
		wantErr := ErrDescriptionMandatory
		if err != wantErr {
			t.Fatalf("expected err: %v got: %v", wantErr, err)
		}
	})

	t.Run("price_zero", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
			Name:        "test",
			Price:       0,
			TermInMonth: 11,
			Description: "test",
		})
		wantErr := ErrPriceCannotBeZero
		if err != wantErr {
			t.Fatalf("expected err: %v got: %v", wantErr, err)
		}
	})

	t.Run("price_negative", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
			Name:        "test",
			Price:       -10,
			TermInMonth: 11,
			Description: "test",
		})
		wantErr := ErrPriceCannotBeNegative
		if err != wantErr {
			t.Fatalf("expected err: %v got: %v", wantErr, err)
		}
	})

	t.Run("term_zero", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
			Name:        "test",
			Price:       10,
			TermInMonth: 0,
			Description: "test",
		})
		wantErr := ErrTermCannotBeZero
		if err != wantErr {
			t.Fatalf("expected err: %v got: %v", wantErr, err)
		}
	})

	t.Run("term_negative", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
			Name:        "test",
			Price:       10,
			TermInMonth: -10,
			Description: "test",
		})
		wantErr := ErrTermCannotBeNegative
		if err != wantErr {
			t.Fatalf("expected err: %v got: %v", wantErr, err)
		}
	})
}

func Test_GetAll_happypath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	_, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
		Name:        "test",
		Price:       101,
		TermInMonth: 11,
		Description: "test",
	})
	if err != nil {
		t.Fatalf("failed to create plan err: %v", err)
	}

	got, err := container.SubscriptionService.GetAllPlan(ctx)
	if err != nil {
		t.Fatalf("failed to get plans err: %v", err)
	}
	wantLen := 1
	if len(got.List) != wantLen {
		t.Fatalf("expected len: %v got: %v", wantLen, len(got.List))
	}
}

func Test_GetAll_unhappypath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	_, err := container.SubscriptionService.GetAllPlan(ctx)

	wantErr := ErrEmptyGetContent
	if err != wantErr {
		t.Fatalf("expected err: %v got: %v", wantErr, err)
	}
}

func Test_Get_happypath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	planId, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
		Name:        "test",
		Price:       101,
		TermInMonth: 11,
		Description: "test",
	})
	if err != nil {
		t.Fatalf("failed to create plan err: %v", err)
	}

	got, err := container.SubscriptionService.GetPlan(ctx, planId)
	if err != nil {
		t.Fatalf("failed to get plans err: %v", err)
	}
	wantId := planId
	if wantId != got.Id {
		t.Fatalf("expected id: %v got: %v", wantId, got.Id)
	}
}

func Test_Get_unhappypath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	randomPlanId := 11
	_, err := container.SubscriptionService.GetPlan(ctx, randomPlanId)

	wantErr := ErrEmptyGetContent
	if err != wantErr {
		t.Fatalf("expected err: %v got: %v", wantErr, err)
	}
}

func Test_Place_happypath(t *testing.T) {
	t.Run("place", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
		_, err = container.SubscriptionService.Place(ctx, &PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   DigitalPaymentPartnerId,
		})
		if err != nil {
			t.Fatalf("failed to place order err: %v", err)
		}
	})
}

func Test_Place_unhappypath(t *testing.T) {
	t.Run("invalidPlan", func(t *testing.T) {
		t.Cleanup(teardown)
		randomPlanId := 112
		ctx := context.Background()
		_, err := container.SubscriptionService.Place(ctx, &PlaceRequest{
			SubscriptionPlanId: randomPlanId,
			DistributorId:      distributorId,
			PaymentPartnerId:   DigitalPaymentPartnerId,
		})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("expected err: %v got: %v", wantErr, err)
		}
	})

	t.Run("paymentPartnerIdMandatory", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
		randomPaymentPartnerId := 1112
		_, err = container.SubscriptionService.Place(ctx, &PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   randomPaymentPartnerId,
		})
		wantErr := ErrPaymentPartnerIdNotSupported
		if err != wantErr {
			t.Errorf("expected err: %v got: %v", wantErr, err)
		}
	})

	t.Run("subscriptionAlreadyExists", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		planId, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
			Name:        "test",
			Price:       101,
			TermInMonth: 11,
			Description: "test",
		})
		if err != nil {
			t.Fatalf("failed to create plan err: %v", err)
		}
		_, err = container.SubscriptionService.Place(ctx, &PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   DigitalPaymentPartnerId,
		})
		if err != nil {
			t.Fatalf("failed to place order err: %v", err)
		}
		_, err = container.SubscriptionService.Place(ctx, &PlaceRequest{
			SubscriptionPlanId: planId,
			DistributorId:      distributorId,
			PaymentPartnerId:   DigitalPaymentPartnerId,
		})
		wantErr := ErrSubscriptionAlreadyExistsForDistributor
		if err != wantErr {
			t.Errorf("subscription already exists err: %v", err)
		}
	})
}

func Test_GetAllSubscriptions_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	planId, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
		Name:        "test",
		Price:       101,
		TermInMonth: 11,
		Description: "test",
	})
	if err != nil {
		t.Fatalf("failed to create plan err: %v", err)
	}
	_, err = container.SubscriptionService.Place(ctx, &PlaceRequest{
		SubscriptionPlanId: planId,
		DistributorId:      distributorId,
		PaymentPartnerId:   DigitalPaymentPartnerId,
	})
	if err != nil {
		t.Fatalf("failed to place order err: %v", err)
	}
	all_resp, err := container.SubscriptionService.GetSubscriptions(ctx)
	if err != nil {
		t.Fatalf("failed to get plans err: %v", err)
	}
	wantLen := 1
	if len(all_resp.List) != wantLen {
		t.Errorf("Expected len: %v Got: %v", wantLen, len(all_resp.List))
	}
}

func Test_GetSubscription_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	planId, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
		Name:        "test",
		Price:       101,
		TermInMonth: 11,
		Description: "test",
	})
	if err != nil {
		t.Fatalf("failed to get subscription err: %v", err)
	}
	sub, err := container.SubscriptionService.Place(ctx, &PlaceRequest{
		SubscriptionPlanId: planId,
		DistributorId:      distributorId,
		PaymentPartnerId:   DigitalPaymentPartnerId,
	})
	if err != nil {
		t.Fatalf("failed to get subsctiption err: %v", err)
	}

	all_resp, err := container.SubscriptionService.GetSubscription(ctx, sub.Id)
	if err != nil {
		t.Fatalf("failed to get plans err: %v", err)
	}
	if all_resp.Id != sub.Id {
		t.Errorf("Expected id: %v Got: %v", sub.Id, all_resp.Id)
	}
}

func Test_GetAllSubscriptions_unhappyPath(t *testing.T) {
	t.Run("emptyGetContent", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := container.SubscriptionService.GetSubscriptions(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("expected err: %v got: %v", wantErr, err)
		}
	})
}

func Test_GetSubscriptionByDistributorId_happyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	planId, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
		Name:        "test",
		Price:       101,
		TermInMonth: 11,
		Description: "test",
	})
	if err != nil {
		t.Fatalf("failed to create plan err: %v", err)
	}
	sub_id, err := container.SubscriptionService.Place(ctx, &PlaceRequest{
		SubscriptionPlanId: planId,
		DistributorId:      distributorId,
		PaymentPartnerId:   DigitalPaymentPartnerId,
	})
	if err != nil {
		t.Fatalf("failed to place order err: %v", err)
	}
	got, err := container.SubscriptionService.GetSubscriptionByDistributorId(ctx, distributorId)
	if err != nil {
		t.Fatalf("failed to get subscriptions by distId err: %v", err)
	}
	wantSubId := sub_id.Id
	if wantSubId != got.Id {
		t.Errorf("expected sub Id: %v got :%v", wantSubId, got.Id)
	}
}

func Test_GetSubscriptionByDistributorId_unhappyPath(t *testing.T) {
	t.Run("emptyGetContent", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := container.SubscriptionService.GetSubscriptionByDistributorId(ctx, distributorId)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("expected err: %v got: %v", wantErr, err)
		}
	})
}
