package config

import "context"

type GetOrderExpiryResponse struct {
	ExpiryDurationInHours int
}

type SetOrderExpiryRequest struct {
	ExpiryDurationInHours int
}

type Reader interface {
	//order expiry
	GetOrderExpiry(context.Context) (GetOrderExpiryResponse, error)
}

type Writer interface {
	SetOrderExpiryConfig(ctx context.Context, req *SetOrderExpiryRequest) error
}

type DB interface {
	Reader
	Writer
}
