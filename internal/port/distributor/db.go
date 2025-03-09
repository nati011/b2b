package distributor

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type CreateRequest struct {
	FirstName  string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
	Password   string
}

type GetResponse struct {
	Id         int
	FirstName  string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByIdRequest struct {
	Id int
}

type Reader interface {
	GetById(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}
