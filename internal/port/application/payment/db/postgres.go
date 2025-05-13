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
	Date           time.Time
	Amount         float64
	OrderId        int
	PartnerId      int
	TransactionRef string
}

type CreateRequest struct {
	OrderId        int
	PartnerId      int
	TransactionRef string
}

type Reader interface {
	GetByID(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByTransactionRef(ctx context.Context, txRef string) (GetResponse, error)
}

type Writer interface {
	Create(context.Context, CreateRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}
