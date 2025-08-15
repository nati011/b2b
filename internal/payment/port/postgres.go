package db

import (
	"context"
	"time"
)

type GetAllResponse struct {
	List []GetResponse
}
type GetResponse struct {
	Id             int
	OrderId        int
	PartnerId      int
	TransactionRef string
	Amount         float64
	Date           time.Time
}

type CreateRequest struct {
	OrderId        int
	PartnerId      int
	TransactionRef string
	Amount         float64
}

type Reader interface {
	GetByID(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByTransactionRef(ctx context.Context, txRef string) (GetResponse, error)
	GetByOrderId(ctx context.Context, orderId int) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, CreateRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}
