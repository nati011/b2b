package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	invoice "b2b.nati011.github.com/internal/core/domain/invoice"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var invoiceService invoice.Provider
var db *sql.DB
var orderId = 1

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	invoiceService = invoice.NewInvoice(
		invoice_db.NewPostgres(db),
	)
	//create order
	// orderService := order.NewOrderService(
	// 	order_db.NewPostgres(db),
	// 	invoiceService,
	// )
	// // retailerService := retailer.NewRetailerService()

	// orderId, err = orderService.Place(ctx,
	// 	&order.PlaceRequest{
	// 		RetailerId: 1,
	// 		Items: []order.Item{
	// 			{
	// 				ProductId: 1,
	// 				Quantity:  19},
	// 		},
	// 	})
	// if err != nil {
	// 	panic("failed to create order")
	// }

}

func teardown() {
	db_test_container.Teardown(db)
}

func Test_Timeout(t *testing.T) {}

func Test_read(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			OrderId:    1,
			ExternalId: "test",
			Status:     "test",
			SubTotal:   1,
			TaxAmount:  1,
			LineItems: []invoice.Item{
				{ProductId: 1, ProductName: "1", ProductQuantity: 1, ProductPrice: 1},
			},
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
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    orderId,
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
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    orderId,
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
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    orderId,
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
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			OrderId:    1,
			ExternalId: "test",
			Status:     "test",
			SubTotal:   1,
			TaxAmount:  1,
			LineItems: []invoice.Item{
				{1, "1", 1, 1},
			},
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
		t.Cleanup(teardown)
		ctx := context.Background()
		in := &invoice.CreateRequest{
			OrderId:    1,
			ExternalId: "test",
			Status:     "test",
			SubTotal:   1,
			TaxAmount:  1,
			LineItems: []invoice.Item{
				{1, "1", 1, 1},
			},
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
			t.Errorf("Expected OrderId: %v Got: %v", 1, resp.OrderId)
		}
		if resp.ExternalId != "test" {
			t.Errorf("Expected ExternalId: %v Got: %v", "test", resp.ExternalId)
		}
		if resp.Status != "test" {
			t.Errorf("Expected Status: %v Got: %v", "test", resp.Status)
		}
		if resp.SubTotal != 1 {
			t.Errorf("Expected subtotal: %v Got: %v", 1, resp.SubTotal)
		}
		if resp.TaxAmount != 1 {
			t.Errorf("Expected taxAmount: %v Got: %v", 1, resp.TaxAmount)
		}
		if len(resp.LineItems) != 1 {
			t.Errorf("Expected lineItems len: %v Got %v", 1, len(resp.LineItems))
		}
	})
	t.Run("UpdateStatus", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			OrderId:    orderId,
			ExternalId: "test",
			Status:     "test",
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
		t.Cleanup(teardown)
		ctx := context.Background()
		//setup
		in := &invoice.CreateRequest{
			ExternalId: "test",
			Status:     "test",
			OrderId:    orderId,
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

		if resp.ExternalId != "new_externalId" {
			t.Errorf("Expected extId: %v Got: %v", "new_externalId", resp.Status)
		}
	})
}
