package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"b2b.nati011.github.com/internal/core/domain/product"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var productService product.Provider
var container product.TestContainer
var pgContainer *postgres.PostgresContainer
var db *sql.DB
var categoryId = 1

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
	container = product.NewPackageIntegrationTestContainer()
	productService = container.ProductService
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

func Test_Timeout(t *testing.T) {
}

func Test_read(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		//setup
		ctx := context.Background()
		in := &product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}

		_, err := productService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		//get
		got, err := productService.GetAll(ctx)
		if err != nil {
			t.Errorf("Expected err:%v Got err: %v", nil, err)
		}
		wantNum := 1
		if len(got.List) != wantNum {
			t.Errorf("Expected len: %v, Got len: %v", wantNum, len(got.List))
		}
	})
	t.Run("GetAll", func(t *testing.T) {
		//setup
		ctx := context.Background()
		in := &product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}

		_, err := productService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		in_two := &product.CreateRequest{
			Name:       "test2",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}

		_, err = productService.Create(ctx, in_two)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		//get-all
		resp, err := productService.GetAll(ctx)
		if err != nil {
			t.Errorf("Expected err:%v Got err: %v", nil, err)
		}
		expected_resp_len := 2
		if len(resp.List) != expected_resp_len {
			t.Errorf("Expected resp len:%v Got resp len: %v", expected_resp_len, len(resp.List))
		}
	})

	t.Run("GetByName", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}

		_, err := productService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := productService.GetByParam(ctx, &product.GetByParamRequest{
			Name: "test",
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got err: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got err: %v", wantLen, len(got.List))
		}
	})

	t.Run("GetByExternalId", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}

		_, err := productService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := productService.GetByParam(ctx, &product.GetByParamRequest{
			ExternalID: "123",
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got err: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got err: %v", wantLen, len(got.List))
		}
	})

	t.Run("GetByDistributorId", func(t *testing.T) {
		// setup
		ctx := context.Background()
		in := &product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
			DistributorId: 1,
		}

		_, err := productService.Create(ctx, in)
		if err != nil {
			t.Fatalf("Failed to create product")
		}

		got, err := productService.GetByParam(ctx, &product.GetByParamRequest{
			DistributorId: 1,
		})
		if err != nil {
			t.Errorf("Expected err: %v, Got err: %v", nil, err)
		}
		wantLen := 1
		if wantLen != len(got.List) {
			t.Errorf("Expected len: %v Got err: %v", wantLen, len(got.List))
		}
	})

	t.Run("GetByCategory", func(t *testing.T) {})
	t.Run("GetByPriceRange", func(t *testing.T) {})
}

func Test_write(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		ctx := context.Background()
		in := &product.CreateRequest{
			Name:       "test",
			Desc:       "test",
			ExternalID: "123",
			Images: []string{
				"test",
				"test",
			},
			Price: 100.00,
			Attributes: map[string]string{
				"test": "test",
			},
		}

		id, err := productService.Create(ctx, in)
		if err != nil {
			t.Errorf("Failed to create product err: %v", err)
		}

		//get
		resp, err := productService.Get(ctx, id)
		if err != nil {
			t.Errorf("Expected err:%v Got err: %v", nil, err)
		}
		if resp.Id == 0 {
			t.Errorf("No product created")
		}
	})

	t.Run("updateName", func(t *testing.T) {})
	t.Run("updateExternalId", func(t *testing.T) {})
	t.Run("updatePrice", func(t *testing.T) {})
	t.Run("updateDesc", func(t *testing.T) {})
	t.Run("updateImages", func(t *testing.T) {})
	t.Run("updateActiveStatus", func(t *testing.T) {})
	t.Run("updateCategory", func(t *testing.T) {})
	t.Run("goodsReceiving", func(t *testing.T) {})
	t.Run("dispatch", func(t *testing.T) {})
}
