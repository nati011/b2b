package config

import (
	"context"
	"log"

	port "b2b.nati011.github.com/internal/port/domain/config"
)

type GetOrderExpiryResponse struct {
	ExpiryDurationInHours int
}

type SetOrderExpiryRequest struct {
	ExpiryDurationInHours int
}

type Provider interface {
	SetDefaults(ctx context.Context) error

	//order expiry
	GetOrderExpiryConfig(ctx context.Context) (GetOrderExpiryResponse, error)
	SetOrderExpiryConfig(ctx context.Context, req *SetOrderExpiryRequest) error
	ResetOrderExpiryConfig(ctx context.Context) error
}

type ConfigService struct {
	DB port.DB
}

func NewConfig(DB port.DB) Provider {
	cfgService := &ConfigService{
		DB: DB,
	}
	ctx := context.Background()
	cfgService.SetDefaults(ctx)
	return cfgService
}

const (
	DEFAULT_ORDER_EXPIRY_HOURS = 12
)

func (c *ConfigService) SetDefaults(ctx context.Context) error {
	err := c.SetOrderExpiryConfig(ctx, &SetOrderExpiryRequest{
		ExpiryDurationInHours: DEFAULT_ORDER_EXPIRY_HOURS,
	})
	if err != nil {
		log.Printf("failed to set default order expiry err: %v", err)
		return err
	}
	return nil
}

func (c *ConfigService) GetOrderExpiryConfig(ctx context.Context) (GetOrderExpiryResponse, error) {
	resp, err := c.DB.GetOrderExpiry(ctx)
	if err != nil {
		log.Printf("failed to set order expiry err: %v", err)
		return GetOrderExpiryResponse{}, err
	}
	return GetOrderExpiryResponse{
		ExpiryDurationInHours: resp.ExpiryDurationInHours,
	}, nil
}

func (c *ConfigService) SetOrderExpiryConfig(ctx context.Context, req *SetOrderExpiryRequest) error {
	err := c.DB.SetOrderExpiryConfig(ctx, &port.SetOrderExpiryRequest{
		ExpiryDurationInHours: req.ExpiryDurationInHours,
	})
	if err != nil {
		log.Printf("failed to set order expiry err: %v", err)
		return err
	}
	return nil
}

func (c *ConfigService) ResetOrderExpiryConfig(ctx context.Context) error {
	if err := c.SetDefaults(ctx); err != nil {
		return err
	}
	return nil
}
