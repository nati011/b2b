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

type UpdateRequest struct {
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string
	Id          int
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
	GetByName(ctx context.Context, name string) (GetResponse, error)
	GetByTin(ctx context.Context, tin string) (GetResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req CreateRequest) (int, error)
	Update(ctx context.Context, req *UpdateRequest) error
}

type DB interface {
	Reader
	Writer
}
