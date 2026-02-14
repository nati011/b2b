package db

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"b2b.nati011.github.com/internal/core/domain/config"
	"b2b.nati011.github.com/internal/core/domain/config/test"
	db_test_container "b2b.nati011.github.com/internal/core/util/test_container/db"
)

var container test.IntegrationTestContainer
var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	db = db_test_container.Setup()
	container = test.NewIntegrationTestContainer(db)
}

func teardown() {
	container.Teardown()
	db_test_container.Teardown(db)
}

func Test_GetOrderExpiryConfig_happypath(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(teardown)
	resp, err := container.ConfigService.GetOrderExpiryConfig(ctx)
	if err != nil {
		t.Fatalf("Failed to get order expiry config err: %v", err)
	}
	if resp.ExpiryDurationInMinues != config.DEFAULT_ORDER_EXPIRY_MINUTES {
		t.Errorf("Expected ExpiryDurationInHours: %v Got: %v", config.DEFAULT_ORDER_EXPIRY_MINUTES, resp.ExpiryDurationInMinues)
	}
}

func Test_SetOrderExpiryConfig_unhappypath(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(teardown)
	expiryHours := 11
	err := container.ConfigService.SetOrderExpiryConfig(ctx, &config.SetOrderExpiryRequest{
		ExpiryDurationInMinues: expiryHours,
	})
	if err != nil {
		t.Errorf("Failed to set order expiry config err:%v", err)
	}
	resp, err := container.ConfigService.GetOrderExpiryConfig(ctx)
	if err != nil {
		t.Errorf("Failed to get order expiry config err:%v", err)
	}
	if resp.ExpiryDurationInMinues != expiryHours {
		t.Errorf("Expected ExpiryDurationInHours: %v Got: %v", expiryHours, resp.ExpiryDurationInMinues)
	}
}
