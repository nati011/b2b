package payment

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSysNoRows = errors.New("no rows")
)

type CreateRequest struct {
	User_Id int
	Amount  int
}

type GetResponse struct {
	Id           int
	User_Id      int
	Amount       int
	Created_Date time.Time
}

type GetAllResponse struct {
	List []GetResponse
}

type Reader interface {
	GetAll(context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
	UpdateStatus(context.Context, int, string) (int, error)
	Delete(context.Context, int) error
}

type DB interface {
	Reader
	Writer
}
