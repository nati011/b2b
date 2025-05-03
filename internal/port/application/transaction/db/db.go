package db

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type GetResponse struct {
	Id        int
	Date      time.Time
	Amount    float64
	PartnerId int
	TxRef     string
	Status    string
}

type GetAllResponse struct {
	List []GetResponse
}

type CreateRequest struct {
	Amount    float64
	PartnerId int
	TxRef     string
	Status    string
	OrderId   int
}

type UpdateRequest struct {
	Id     int
	Status string
}

type UpdateByTransactionRefRequest struct {
	TransactionRef string
	Status         string
}

type Reader interface {
	GetByID(context.Context, int) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
	GetByDate(context.Context, time.Time) (GetAllResponse, error)
	GetByTxRef(context.Context, string) (GetAllResponse, error)
	GetByPartnerId(context.Context, int) (GetAllResponse, error)
	GetByStatus(context.Context, string) (GetAllResponse, error)
	UpdateStatus(context.Context, *UpdateRequest) error
	UpdateByTransactionRef(context.Context, *UpdateByTransactionRefRequest) error
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}
