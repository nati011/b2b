package transaction

import (
	"context"
	"os"
	"testing"
	"time"

	db "b2b.nati011.github.com/internal/adapter/secondary/application/transaction/db"
)

var service Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	service = NewTransactionService(
		db.NewMock(),
	)
}

func Test_Save_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			User_Id:    1,
			Amount:     1,
			Partner_Id: 1,
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		//check
		resp, err := service.Get(ctx, id)
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
			User_Id:    1,
			Amount:     1,
			Partner_Id: 1,
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}

		resp, err := service.Get(ctx, id)
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
			Amount: 1,
		}
		_, err := service.Create(ctx, in)
		if err != ErrUserIdNotSupplied {
			t.Fatalf("Failed to create err: %v", err)
		}
	})

	t.Run("amount_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			User_Id: 1,
		}
		_, err := service.Create(ctx, in)
		if err != ErrAmountIsNotSupplied {
			t.Fatalf("Failed to create err: %v", err)
		}
	})

	t.Run("partnerId_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			User_Id: 1,
			Amount:  1,
		}
		_, err := service.Create(ctx, in)
		if err != ErrPartnerIdNotSupplied {
			t.Fatalf("Failed to create err: %v", err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		User_Id:    1,
		Amount:     1,
		Partner_Id: 1,
	}
	id, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	if resp.Id != id {
		t.Errorf("Expected Id: %v Want: %v", id, resp.Id)
	}

}

func Test_Get_unhappyPath(t *testing.T) {
	ctx := context.Background()
	_, err := service.Get(ctx, 99)
	wantErr := ErrIdNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v Got err: %v", wantErr, err)
	}

}

func Test_GetAll_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		User_Id:    1,
		Amount:     1,
		Partner_Id: 1,
	}
	_, err := service.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	resp, err := service.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	wantLen := 1
	if len(resp.List) != wantLen {
		t.Errorf("Expected length: %v Want: %v", len(resp.List), wantLen)
	}
}

func Test_GetAll_unhappyPath(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ctx := context.Background()
		_, err := service.GetAll(ctx)
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
			User_Id:    1,
			Amount:     1,
			Partner_Id: 1,
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		resp_param, err := service.GetByParam(ctx, &GetByParamRequest{
			User_Id: resp.User_Id,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}
		wantLen := 1
		if wantLen != len(resp_param.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp_param.List))
		}
	})

	t.Run("partner", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			User_Id:    1,
			Amount:     1,
			Partner_Id: 1,
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		resp_param, err := service.GetByParam(ctx, &GetByParamRequest{
			Partner_Id: resp.Partner_Id,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantLen := 1
		if wantLen != len(resp_param.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp_param.List))
		}
	})

	t.Run("date", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			User_Id:    1,
			Amount:     1,
			Partner_Id: 1,
		}
		id, err := service.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create err: %v", err)
		}
		resp, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		resp_param, err := service.GetByParam(ctx, &GetByParamRequest{
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
		_, err := service.GetByParam(ctx, &GetByParamRequest{})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
	})
}
