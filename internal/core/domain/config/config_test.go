package config

import (
	"context"
	"os"
	"testing"
)

var container TestContainer

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	container = NewTestContainer()
}

func teardown() {
	container.Teardown()
}

func Test_GetOrderExpiryConfig_happypath(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(teardown)
	resp, err := container.ConfigService.GetOrderExpiryConfig(ctx)
	if err != nil {
		t.Errorf("Failed to get order expiry config err:%v", err)
	}
	if resp.ExpiryDurationInHours != DEFAULT_ORDER_EXPIRY_HOURS {
		t.Errorf("Expected ExpiryDurationInHours: %v Got: %v", DEFAULT_ORDER_EXPIRY_HOURS, resp.ExpiryDurationInHours)
	}
}

func Test_SetOrderExpiryConfig_unhappypath(t *testing.T) {
	ctx := context.Background()
	t.Cleanup(teardown)
	expiryHours := 11
	err := container.ConfigService.SetOrderExpiryConfig(ctx, &SetOrderExpiryRequest{
		ExpiryDurationInHours: expiryHours,
	})
	if err != nil {
		t.Errorf("Failed to set order expiry config err:%v", err)
	}
	resp, err := container.ConfigService.GetOrderExpiryConfig(ctx)
	if err != nil {
		t.Errorf("Failed to get order expiry config err:%v", err)
	}
	if resp.ExpiryDurationInHours != expiryHours {
		t.Errorf("Expected ExpiryDurationInHours: %v Got: %v", expiryHours, resp.ExpiryDurationInHours)
	}
}
