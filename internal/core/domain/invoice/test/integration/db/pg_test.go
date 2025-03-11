package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	invoice_db "b2b.nati011.github.com/internal/adapter/secondary/domain/invoice/db"
	invoice "b2b.nati011.github.com/internal/core/domain/invoice"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var invoiceService invoice.Provider
var pgContainer *postgres.PostgresContainer
var db *sql.DB
var orderId = 1

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	var err error
	ctx := context.Background()

	pgContainer, err = RunContainer(ctx)
	if err != nil {
		panic(err)
	}

	connectionString, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	db, err = sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// ddl
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db.sql")
	if err != nil {
		log.Fatalf("Error running ddl migration: %v", err)
	}

	// functions
	err = runMigration(db, "/home/natanel/personal/b2b_clean/b2b/migration/core_db_functions.sql")
	if err != nil {
		log.Fatalf("Error running stored func migration: %v", err)
	}
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
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("could not begin transaction: %v", err)
	}

	// Get all table names
	var tables []string
	rows, err := tx.Query("SELECT tablename FROM pg_tables WHERE schemaname = 'public';")
	if err != nil {
		tx.Rollback()
		log.Fatalf("could not fetch table names: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			tx.Rollback()
			log.Fatalf("could not scan table name: %v", err)
		}
		tables = append(tables, table)
	}

	// Prepare the TRUNCATE statement
	if len(tables) > 0 {
		truncateQuery := "TRUNCATE TABLE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE;"
		_, err = tx.Exec(truncateQuery)
		if err != nil {
			tx.Rollback()
			log.Fatalf("could not truncate tables: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("could not commit transaction: %v", err)
	}
}

func runMigration(db *sql.DB, filename string) error {
	// Read the SQL file
	sqlBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	// Execute the SQL
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("could not execute SQL: %w", err)
	}

	return nil
}

func RunContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	return postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),

		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
}
func Test_Timeout(t *testing.T) {}

func Test_read(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
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
			ExternalId: "test",
			Status:     "test",
			OrderId:    orderId,
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
			OrderId:    orderId,
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

		if resp.Status != "new_status" {
			t.Errorf("Expected status: %v Got: %v", "new_status", resp.Status)
		}

		if resp.ExternalId != "new_externalId" {
			t.Errorf("Expected extId: %v Got: %v", "new_externalId", resp.Status)
		}
	})
}
