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
	PartnerId      int
	TransactionRef string
}

type CreateRequest struct {
	Amount         float64
	PartnerId      int
	TransactionRef string
}

type Reader interface {
	GetByID(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByTransactionRef(txRef string, ctx context.Context) (GetResponse, error)
}

type Writer interface {
	Create(context.Context, CreateRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}
