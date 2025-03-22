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
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string
	UserId      int
}

type UpdateNameRequest struct {
	Name string
	Id   int
}

type UpdateTinRequest struct {
	Tin string
	Id  int
}

type GetResponse struct {
	Id          int
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string
}

type GetAllResponse struct {
	List []GetResponse
}

type Reader interface {
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByName(ctx context.Context, name string) (GetAllResponse, error)
	GetByTin(ctx context.Context, tin string) (GetResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req CreateRequest) (int, error)
	UpdateName(ctx context.Context, req *UpdateNameRequest) error
	UpdateTin(ctx context.Context, tin *UpdateTinRequest) error
}

type DB interface {
	Reader
	Writer
}
