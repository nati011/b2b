package category

import (
	"context"
	"errors"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type CreateRequest struct {
	Name string
	Desc string
}

type GetResponse struct {
	Id   int
	Name string
	Desc string
}

type GetAllResponse struct {
	List []GetResponse
}

type Reader interface {
	GetAll(ctx context.Context) (GetAllResponse, error)
	Get(ctx context.Context, id int) (GetResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	Remove(ctx context.Context, id int) error
}

type DB interface {
	Reader
	Writer
}
