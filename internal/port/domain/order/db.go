package order

import (
	"context"
	"errors"
)

var (
	ErrSysUnknown = errors.New("unknown error")
	ErrSysNoRows  = errors.New("no rows")
)

type Item struct {
	ProductId int
	Quantity  int
}
type CreateRequest struct {
	RetailerId int
	Items      []Item
}

type GetResponse struct {
	Id         int
	RetailerId int
	Items      []Item
	Total      float32
	Status     string
}

type GetAllResponse struct {
	List []GetResponse
}

type UpdateOrderStatusRequest struct {
	Id     int
	Status string
}

type Reader interface {
	GetByID(context.Context, int) (GetResponse, error)
	GetByRetailerID(context.Context, int) (GetAllResponse, error)
	GetByStatus(context.Context, string) (GetAllResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
	UpdateOrderStatus(context.Context, *UpdateOrderStatusRequest) error
}

type DB interface {
	Reader
	Writer
}
