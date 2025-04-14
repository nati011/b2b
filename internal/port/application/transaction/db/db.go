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
	User_Id    int
	Date       time.Time
	Amount     int64
	Partner_Id int
}

type Pagination struct {
	Limit  int
	Offset int
}

type GetAllResponse struct {
	List []GetResponse
}

type CreateRequest struct {
	User_Id    int
	Amount     int64
	Partner_Id int
}

type Reader interface {
	GetByID(context.Context, int) (GetResponse, error)
	GetAll(context.Context, *Pagination) (GetAllResponse, error)
	GetByDate(context.Context, time.Time, *Pagination) (GetAllResponse, error)
	GetByUserId(context.Context, int, *Pagination) (GetAllResponse, error)
	GetByPartnerId(context.Context, int, *Pagination) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
}

type DB interface {
	Reader
	Writer
}
