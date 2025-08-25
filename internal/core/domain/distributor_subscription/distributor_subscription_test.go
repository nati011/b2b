package distributor_subscription

import (
	"context"
	"os"
	"testing"
)

var container TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = NewTestContainer()
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
		})
		wantErr := ErrNameMandatory
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
		})
		wantErr := ErrTermCannotBeNegative
		if err != wantErr {
			t.Fatalf("expected err: %v got: %v", wantErr, err)
		}
	})
}

func Test_Get_happypath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	_, err := container.SubscriptionService.CreatePlan(ctx, &CreatePlanRequest{
		Name:        "test",
		Price:       101,
		TermInMonth: 11,
	})
	if err != nil {
		t.Fatalf("failed to create plan err: %v", err)
	}

	got, err := container.SubscriptionService.GetPlan(ctx)
	if err != nil {
		t.Fatalf("failed to get plans err: %v", err)
	}
	wantLen := 1
	if len(got.List) != wantLen {
		t.Fatalf("expected len: %v got: %v", wantLen, len(got.List))
	}
}

func Test_Get_unhappypath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	_, err := container.SubscriptionService.GetPlan(ctx)

	wantErr := ErrEmptyGetContent
	if err != wantErr {
		t.Fatalf("expected err: %v got: %v", wantErr, err)
	}
}
