package retailer

import (
	"context"
	"errors"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type CreateRequest struct {
	RetailerId string
	UserId     string
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
}

type DB interface {
	Writer
}
