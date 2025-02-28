package invoice

import (
	"context"
	"os"
	"testing"

	db "b2b.nati011.github.com/internal/adapter/secondary/invoice"
)

var invoiceService Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	invoiceService = NewInvoice(
		db.NewMock(),
	)
}

func Test_CreateInvoice_happyPath(t *testing.T) {
	ctx := context.Background()
	in := &CreateRequest{
		OrderId:    1,
		ExternalId: "test",
		Status:     "test",
		SubTotal:   1,
	}
	id, err := invoiceService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create invoice err: %v", err)
	}
	//verify
	resp, err := invoiceService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get invocie err: %v", err)
	}
	if resp.Id != id {
		t.Errorf("Expected id: %v Got: %v", id, resp.Id)
	}
	if resp.OrderId != 1 {
		t.Errorf("Expected OrderId: %v Got: %v", id, resp.Id)
	}
	if resp.ExternalId != "test" {
		t.Errorf("Expected ExternalId: %v Got: %v", id, resp.Id)
	}
	if resp.Status != "test" {
		t.Errorf("Expected ExternalId: %v Got: %v", id, resp.Id)
	}
}

func Test_CreateInvoice_unhappyPath(t *testing.T) {
	t.Run("status_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			ExternalId: "test",
			OrderId:    1,
			SubTotal:   1,
		}
		_, err := invoiceService.Create(ctx, in)
		expectedErr := ErrSysStatusNotSupplied
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})

	t.Run("orderId_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			ExternalId: "test",
			Status:     "test",
			SubTotal:   1,
		}
		_, err := invoiceService.Create(ctx, in)
		expectedErr := ErrSysOrderIdNotSupplied
		if err != expectedErr {
			t.Errorf("Expected err: %v Got err: %v", expectedErr, err)
		}
	})
}

func Test_Update_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		ExternalId: "test",
		Status:     "test",
		OrderId:    1,
	}
	id, err := invoiceService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create invoice err: %v", err)
	}

	//update
	err = invoiceService.Update(ctx, &UpdateByParamRequest{
		Id:         id,
		Status:     "new_status",
		ExternalId: "new_externalId",
	})
	if err != nil {
		t.Fatalf("Failed to update invoice status err: %v", err)
	}

	//check
	resp, err := invoiceService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to Get invocie err: %v", err)
	}

	if resp.Status != "new_status" {
		t.Errorf("Expected status: %v Got: %v", "new_status", resp.Status)
	}

	if resp.ExternalId != "new_externalId" {
		t.Errorf("Expected extId: %v Got: %v", "new_externalId", resp.Status)
	}
}

func Test_Update_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		err := invoiceService.Update(ctx, &UpdateByParamRequest{
			Id:     99,
			Status: "test",
		})
		wantErr := ErrSysIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_GetInvoice_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		ExternalId: "test",
		Status:     "test",
		OrderId:    1,
	}
	id, err := invoiceService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create invoice err: %v", err)
	}

	resp, err := invoiceService.Get(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get invoice err: %v", err)
	}
	if resp.Id != id {
		t.Errorf("Extected id: %v Got: %v", id, resp.Id)
	}
}

func Test_GetInvoice_unhappyPath(t *testing.T) {
	t.Run("idNotFound", func(t *testing.T) {
		ctx := context.Background()
		err := invoiceService.Update(ctx, &UpdateByParamRequest{
			Id:     99,
			Status: "test",
		})
		wantErr := ErrSysIdNotFound
		if err != wantErr {
			t.Errorf("Expected err: %v Got err: %v", wantErr, err)
		}
	})
}

func Test_GetByParam_happyPath(t *testing.T) {
	t.Run("getByStatus", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    1,
		}
		_, err := invoiceService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create invoice err: %v", err)
		}
		resp, err := invoiceService.GetByParam(ctx, &GetByParamRequest{
			Status: "test",
		})
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}

		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(resp.List))
		}
	})

	t.Run("getByExternalId", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    1,
		}
		_, err := invoiceService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create invoice err: %v", err)
		}
		resp, err := invoiceService.GetByParam(ctx, &GetByParamRequest{
			ExternalId: "test",
		})
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}

		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(resp.List))
		}
	})

	t.Run("getByOrderId", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    1,
		}
		_, err := invoiceService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create invoice err: %v", err)
		}
		resp, err := invoiceService.GetByParam(ctx, &GetByParamRequest{
			OrderId: 1,
		})
		if err != nil {
			t.Fatalf("Failed to get by param err: %v", err)
		}

		wantLen := 1
		if len(resp.List) != wantLen {
			t.Errorf("Expected len: %v Got: %v", wantLen, len(resp.List))
		}
		wantOrderId := 1
		if resp.List[0].OrderId != wantOrderId {
			t.Errorf("Expected orderId: %v Got: %v", wantOrderId, resp.List[0].OrderId)
		}
	})
}

func Test_GetByParam_unhappyPath(t *testing.T) {
	t.Run("emptyGetContent", func(t *testing.T) {
		ctx := context.Background()
		_, err := invoiceService.GetByParam(ctx, &GetByParamRequest{
			ExternalId: "test",
		})
		wantErr := ErrSysEmptyGetContent
		if err != wantErr {
			t.Fatalf("Expected err:%v Got err:%v", wantErr, err)
		}
	})

}

func Test_GetAll_happyPath(t *testing.T) {
	ctx := context.Background()
	//setup
	in := &CreateRequest{
		ExternalId: "test",
		Status:     "test",
		OrderId:    1,
	}
	_, err := invoiceService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create invoice err: %v", err)
	}
	resp, err := invoiceService.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get by param err: %v", err)
	}

	wantLen := 1
	if len(resp.List) != wantLen {
		t.Errorf("Expected len: %v Got: %v", wantLen, len(resp.List))
	}
}

func Test_GetAll_unhappyPath(t *testing.T) {
	t.Run("emptyGetContent", func(t *testing.T) {
		ctx := context.Background()
		_, err := invoiceService.GetAll(ctx)
		wantErr := ErrSysEmptyGetContent
		if err != wantErr {
			t.Fatalf("Expected err:%v Got err:%v", wantErr, err)
		}
	})
}
