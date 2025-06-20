package config

import (
	"context"
	"errors"
	"log"

	port "b2b.nati011.github.com/internal/port/domain/config"
)

var (
	ErrUnknown = errors.New(" unknown error")
)

type GetOrderExpiryResponse struct {
	ExpiryDurationInMinues int
}

type SetOrderExpiryRequest struct {
	ExpiryDurationInMinues int
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
	DEFAULT_ORDER_EXPIRY_MINUTES = 720
)

func (c *ConfigService) SetDefaults(ctx context.Context) error {
	err := c.SetOrderExpiryConfig(ctx, &SetOrderExpiryRequest{
		ExpiryDurationInMinues: DEFAULT_ORDER_EXPIRY_MINUTES,
	})
	if err != nil {
		log.Printf("failed to set default order expiry err: %v", err)
		return ErrUnknown
	}
	return nil
}

func (c *ConfigService) GetOrderExpiryConfig(ctx context.Context) (GetOrderExpiryResponse, error) {
	resp, err := c.DB.GetOrderExpiry(ctx)
	if err != nil {
		log.Printf("failed to get order expiry err: %v", err)
		return GetOrderExpiryResponse{}, ErrUnknown
	}
	return GetOrderExpiryResponse{
		ExpiryDurationInMinues: resp.ExpiryDurationInMinues,
	}, nil
}

func (c *ConfigService) SetOrderExpiryConfig(ctx context.Context, req *SetOrderExpiryRequest) error {
	err := c.DB.SetOrderExpiryConfig(ctx, &port.SetOrderExpiryRequest{
		ExpiryDurationInMinues: req.ExpiryDurationInMinues,
	})
	if err != nil {
		log.Printf("failed to set order expiry err: %v", err)
		return ErrUnknown
	}
	return nil
}

func (c *ConfigService) ResetOrderExpiryConfig(ctx context.Context) error {
	if err := c.SetDefaults(ctx); err != nil {
		return ErrUnknown
	}
	return nil
}
