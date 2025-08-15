package config

import "context"

type GetOrderExpiryResponse struct {
	ExpiryDurationInMinues int
}

type SetOrderExpiryRequest struct {
	ExpiryDurationInMinues int
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
