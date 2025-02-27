package invoice

import (
	"context"
	"os"
	"testing"
)

var invoiceService Provider

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	invoiceService = NewInvoice()
}

func Test_CreateInvoice_happyPath(t *testing.T) {
	ctx := context.Background()
	in := &CreateRequest{
		ExternalId: "test",
		Status:     "test",
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
		t.Errorf("Expected id: %v Got id: %v", id, resp.Id)
	}
}

func Test_CreateInvoice_unhappyPath(t *testing.T) {
	t.Run("status_mandatory", func(t *testing.T) {
		ctx := context.Background()
		in := &CreateRequest{
			ExternalId: "test",
		}
		_, err := invoiceService.Create(ctx, in)
		expectedErr := ErrSysStatusNotSupplied
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
	}
	id, err := invoiceService.Create(ctx, in)
	if err != nil {
		t.Fatalf("Failed to create invoice err: %v", err)
	}

	//update
	err = invoiceService.Update(ctx, &UpdateByParamRequest{
		Status:     "new_status",
		ExternalId: "new_externalId",
	})
	if err != nil {
		t.Fatalf("Failed to update invoice status")
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
