package payment

import (
	"context"
	"errors"
)

var (
	ErrSysNoRows = errors.New("no rows")
)

type CreateRequest struct {
	Name             string
	Icon             string
	Init_payment_url string
}

type GetResponse struct {
	Id     int
	Name   string
	Icon   string
	Status string
}

type GetAllResponse struct {
	List []GetResponse
}

type Reader interface {
	GetByID(context.Context, int) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
	GetByStatus(context.Context, string) (GetAllResponse, error)
	GetByName(context.Context, string) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
	UpdateStatus(context.Context, int, string) (int, error)
	UpdateName(context.Context, int, string) (int, error)
	Delete(context.Context, int) error
}

type DB interface {
	Reader
	Writer
}
