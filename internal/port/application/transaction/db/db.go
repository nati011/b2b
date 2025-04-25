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
	Id         int
	Date       time.Time
	Amount     int64
	Partner_Id int
}

type GetAllResponse struct {
	List []GetResponse
}

type CreateRequest struct {
	Amount     int64
	Partner_Id int
}

type Reader interface {
	GetByID(context.Context, int) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
	GetByDate(context.Context, time.Time) (GetAllResponse, error)
	GetByPartnerId(context.Context, int) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}
