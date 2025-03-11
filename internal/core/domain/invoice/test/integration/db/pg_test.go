package db

import (
	"context"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/invoice"
)

var invoiceService invoice.Provider

func TestMain(m *testing.M) {}

func setup() {}

func teardown() {}

func Test_Timeout(t *testing.T) {}

func Test_read(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
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
	})

	t.Run("GetAll", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
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
	})

	t.Run("GetByStatus", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    1,
		}
		_, err := invoiceService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create invoice err: %v", err)
		}
		resp, err := invoiceService.GetByParam(ctx, &invoice.GetByParamRequest{
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
		in := &invoice.CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    1,
		}
		_, err := invoiceService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create invoice err: %v", err)
		}
		resp, err := invoiceService.GetByParam(ctx, &invoice.GetByParamRequest{
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

	t.Run("GetByOrderId", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    1,
		}
		_, err := invoiceService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create invoice err: %v", err)
		}
		resp, err := invoiceService.GetByParam(ctx, &invoice.GetByParamRequest{
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

func Test_write(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		ctx := context.Background()
		in := &invoice.CreateRequest{
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
	})
	t.Run("UpdateStatus", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    1,
		}
		id, err := invoiceService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create invoice err: %v", err)
		}

		//update
		err = invoiceService.Update(ctx, &invoice.UpdateByParamRequest{
			Id:     id,
			Status: "new_status",
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
	})
	t.Run("UpdateExternalId", func(t *testing.T) {
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    1,
		}
		id, err := invoiceService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create invoice err: %v", err)
		}

		//update
		err = invoiceService.Update(ctx, &invoice.UpdateByParamRequest{
			Id:         id,
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
	})

}
