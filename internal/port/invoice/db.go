package invoice

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSysNoRows = errors.New("no rows found")
)

type GetResponse struct {
	Id           int
	Created_Date time.Time
	ExternalId   string
	Status       string
}

type GetAllResponse struct {
	List []GetResponse
}

type CreateRequest struct {
	ExternalId string
	Status     string
}

type UpdateExternalIdRequest struct {
	Id         int
	ExternalId string
}

type UpdateStatusRequest struct {
	Id     int
	Status string
}

type Reader interface {
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByExternalId(ctx context.Context, extId string) (GetAllResponse, error)
	GetByStatus(ctx context.Context, extId string) (GetAllResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	UpdateExternalId(ctx context.Context, req *UpdateExternalIdRequest) error
	UpdateStatus(ctx context.Context, req *UpdateStatusRequest) error
}

type DB interface {
	Reader
	Writer
}
