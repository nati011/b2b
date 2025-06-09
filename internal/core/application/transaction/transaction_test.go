package transaction

import (
	"context"
	"os"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/application/payment_partner"
	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
)

var testContainer TestContainer
var PartnerId = 0

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = NewPackageIntegrationTestContainer()
	ctx := context.Background()
	//create partner
	var err error
	PartnerId, err = testContainer.PartnerService.Create(ctx, &partner.CreateRequest{
		Name:          "test",
		Icon:          "test",
		BaseURL:       "test",
		Status:        "test",
		Secret:        "test",
		PaymentMethod: payment_partner.PAYMENT_METHOD_DIGITAL,
	})
	if err != nil {
		panic("failed to create transaction")
	}
}

func Test_Save_happyPath(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "232",
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
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "232",
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
	t.Run("amount_mandatory", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		in := &CreateRequest{
			PartnerId: PartnerId,
			TxRef:     "232",
		}
		_, err := testContainer.TransactionService.Create(ctx, in)
		if err != ErrAmountIsNotSupplied {
			t.Fatalf("Failed to create err: %v", err)
		}
	})

	t.Run("amount_greater_than_zero", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Amount:    -1,
			PartnerId: PartnerId,
			TxRef:     "232",
		}
		_, err := testContainer.TransactionService.Create(ctx, in)
		if err != ErrAmountMustBeGreaterThanZero {
			t.Fatalf("Failed to create err: %v", err)
		}
	})

	t.Run("partnerId_mandatory", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Amount: 1,
			TxRef:  "w",
			Status: "te",
		}
		_, err := testContainer.TransactionService.Create(ctx, in)
		if err != ErrPartnerIdNotSupplied {
			t.Fatalf("Failed to create err: %v", err)
		}
	})
	t.Run("tx_ref_mandatory", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		in := &CreateRequest{
			Amount:    1,
			Status:    "2",
			PartnerId: 1,
		}
		_, err := testContainer.TransactionService.Create(ctx, in)
		if err != ErrTransactionRefNotSupplied {
			t.Fatalf("Failed to create err: %v", err)
		}
	})
}

func Test_Get_happyPath(t *testing.T) {
	t.Cleanup(testContainer.TearDown)
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		Amount:    1,
		PartnerId: PartnerId,
		TxRef:     "232",
		Status:    "tets",
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
	t.Cleanup(testContainer.TearDown)
	ctx := context.Background()
	_, err := testContainer.TransactionService.Get(ctx, 99)
	wantErr := ErrIdNotFound
	if err != wantErr {
		t.Errorf("Expected err: %v Got err: %v", wantErr, err)
	}

}

func Test_GetAll_happyPath(t *testing.T) {
	t.Cleanup(testContainer.TearDown)
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		Amount:    1,
		PartnerId: PartnerId,
		TxRef:     "232",
		Status:    "tets",
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
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		_, err := testContainer.TransactionService.GetAll(ctx)
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
	})
}

func Test_GetByParam_happyPath(t *testing.T) {

	t.Run("partner", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "232",
			Status:    "tets",
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
			PartnerId: in.PartnerId,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantLen := 1
		if wantLen > len(resp_param.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp_param.List))
		}

		if resp.PartnerId != in.PartnerId {
			t.Errorf("Expected partnerId: %v Want: %v", in.PartnerId, resp.PartnerId)
		}
	})

	t.Run("date", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "232",
			Status:    "tets",
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

	t.Run("tx_ref", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "232",
			Status:    "tets",
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
			TxRef: in.TxRef,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantLen := 1
		if wantLen != len(resp_param.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp_param.List))
		}

		if resp.TxRef != in.TxRef {
			t.Errorf("Expected transactionRef: %v Want: %v", in.TxRef, resp.TxRef)
		}
	})

	t.Run("status", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			Amount:    1,
			PartnerId: PartnerId,
			TxRef:     "test",
			Status:    "test",
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
			Status: in.Status,
		})
		if err != nil {
			t.Fatalf("Failed to get err: %v", err)
		}

		wantLen := 1
		if wantLen != len(resp_param.List) {
			t.Errorf("Expected len: %v Got len: %v", wantLen, len(resp_param.List))
		}

		if resp.Status != in.Status {
			t.Errorf("Expected status: %v Want: %v", in.TxRef, resp.TxRef)
		}
	})
}

func Test_GetByParam_unhappyPath(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		t.Cleanup(testContainer.TearDown)
		ctx := context.Background()
		_, err := testContainer.TransactionService.GetByParam(ctx, &GetByParamRequest{})
		wantErr := ErrEmptyGetContent
		if err != wantErr {
			t.Errorf("Expected err: %v Got err %v", wantErr, err)
		}
	})
}

func Test_update_status_happyPath(t *testing.T) {
	t.Cleanup(testContainer.TearDown)
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		Amount:    1,
		PartnerId: PartnerId,
		TxRef:     "test",
		Status:    "test",
	}
	id, err := testContainer.TransactionService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create err: %v", err)
	}
	in_status := &UpdateRequest{
		Id:     id,
		Status: "test",
	}
	err = testContainer.TransactionService.UpdateStatus(ctx, in_status)
	if err != nil {
		t.Errorf("Failed to update err %v", err)
	}
	resp, err := testContainer.TransactionService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	if resp.Status != in_status.Status {
		t.Errorf("Expected status: %v Want: %v", in.TxRef, resp.TxRef)
	}
}
