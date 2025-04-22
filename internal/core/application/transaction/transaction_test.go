package transaction

import (
	"context"
	"os"
	"testing"
	"time"

	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/user"
)

var testContainer TestContainer
var user_id = 0
var partner_id = 0

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = NewPackageIntegrationTestContainer()
	ctx := context.Background()
	//create partner
	partner_id, _ = testContainer.PartnerService.Create(ctx, &partner.CreateRequest{
		Name:    "test",
		Icon:    "test",
		BaseUrl: "test",
	})

	//create user
	parsedTime, _ := time.Parse("2006-01-02 15:04:05", "2024-09-19 14:00:00")
	in := user.CreateRequest{
		FirstName:  "natnael jemaneh asefa",
		LastName:   "test",
		Email:      "natnaeljemaneh001@gmail.com",
		Phone:      "+251949184879",
		Username:   "test",
		DOB:        parsedTime,
		ExternalId: "123",
	}
	user_id, _ = testContainer.UserService.Create(ctx, &in)
}

func Test_Save_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			User_Id:    user_id,
			Amount:     1,
			Partner_Id: partner_id,
		}
		id, err := testContainer.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := testContainer.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		if resp.Id != id {
			t.Errorf("Expected id: %v Got: %v", id, resp.Id)
		}
	})

	t.Run("timestamp_date", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			User_Id:    user_id,
			Amount:     1,
			Partner_Id: partner_id,
		}
		id, err := testContainer.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		resp, err := testContainer.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		wantTimeStamp := time.Now()
		if resp.Date.IsZero() {
			t.Errorf("Expected date: %v != 0", wantTimeStamp)
		}
	})
}

func Test_Save_unhappyPath(t *testing.T) {
	t.Run("userId_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			Amount:     1,
			Partner_Id: partner_id,
		}
		_, err := testContainer.TransactionService.Create(ctx, in)
		if err != ErrUserIdNotSupplied {
			t.Fatalf("Failed to create err: %v", err)
		}
	})

	t.Run("amount_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			User_Id:    1,
			Partner_Id: partner_id,
		}
		_, err := testContainer.TransactionService.Create(ctx, in)
		if err != ErrAmountIsNotSupplied {
			t.Fatalf("Failed to create err: %v", err)
		}
	})

	t.Run("amount_greater_than_zero", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			User_Id:    user_id,
			Amount:     -1,
			Partner_Id: partner_id,
		}
		_, err := testContainer.TransactionService.Create(ctx, in)
		if err != ErrAmountMustBeGreaterThanZero {
			t.Fatalf("Failed to create err: %v", err)
		}
	})

	t.Run("partnerId_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			User_Id: user_id,
			Amount:  1,
		}
		_, err := testContainer.TransactionService.Create(ctx, in)
		if err != ErrPartnerIdNotSupplied {
			t.Fatalf("Failed to create err: %v", err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		User_Id:    user_id,
		Amount:     1,
		Partner_Id: partner_id,
	}
	id, err := testContainer.TransactionService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := testContainer.TransactionService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	if resp.Id != id {
		t.Errorf("Expected Id: %v Want: %v", id, resp.Id)
	}

}

func Test_Get_unhappyPath(t *testing.T) {
	ctx := context.Background()
	_, err := testContainer.TransactionService.Get(ctx, 99)
	wantErr := ErrIdNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v Got err: %v", wantErr, err)
	}

}

func Test_GetAll_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		User_Id:    user_id,
		Amount:     1,
		Partner_Id: partner_id,
	}
	_, err := testContainer.TransactionService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := testContainer.TransactionService.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	wantLen := 1
	if len(resp.List) < wantLen {
		t.Errorf("Expected length: %v Got: %v", wantLen, len(resp.List))
	}
}

func Test_GetAll_unhappyPath(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		_, err := testContainer.TransactionService.GetAll(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
	})
}

func Test_GetByParam_happyPath(t *testing.T) {
	t.Run("userId", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			User_Id:    user_id,
			Amount:     1,
			Partner_Id: partner_id,
		}
		id, err := testContainer.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		resp_param, err := testContainer.TransactionService.GetByParam(ctx, &GetByParamRequest{
			User_Id: resp.User_Id,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if wantLen > len(resp_param.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp_param.List))
		}
	})

	t.Run("partner", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			User_Id:    user_id,
			Amount:     1,
			Partner_Id: partner_id,
		}
		id, err := testContainer.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		resp_param, err := testContainer.TransactionService.GetByParam(ctx, &GetByParamRequest{
			Partner_Id: resp.Partner_Id,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantLen := 1
		if wantLen > len(resp_param.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp_param.List))
		}
	})

	t.Run("date", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			User_Id:    user_id,
			Amount:     1,
			Partner_Id: partner_id,
		}
		id, err := testContainer.TransactionService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := testContainer.TransactionService.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		resp_param, err := testContainer.TransactionService.GetByParam(ctx, &GetByParamRequest{
			Date: resp.Date,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantLen := 1
		if wantLen != len(resp_param.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp_param.List))
		}
	})
}

func Test_GetByParam_unhappyPath(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		_, err := testContainer.TransactionService.GetByParam(ctx, &GetByParamRequest{})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
	})
}
