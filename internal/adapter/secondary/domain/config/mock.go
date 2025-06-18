package config

import (
	"context"

	port "b2b.nati011.github.com/internal/port/domain/config"
)

type MockConfigs struct {
	OrderExpiryDurationInHours int
}

type Mock struct {
	configs MockConfigs
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) GetOrderExpiry(ctx context.Context) (port.GetOrderExpiryResponse, error) {
	return port.GetOrderExpiryResponse{
		ExpiryDurationInHours: m.configs.OrderExpiryDurationInHours,
	}, nil
}

func (m *Mock) SetOrderExpiryConfig(ctx context.Context, req *port.SetOrderExpiryRequest) error {
	m.configs.OrderExpiryDurationInHours = req.ExpiryDurationInHours
	return nil
}
