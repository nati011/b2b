package invoice

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSysNoRows  = errors.New("no rows found")
	ErrSysUnknown = errors.New("unknown eror")
)

type Item struct {
	ProductId       int
	ProductName     string
	ProductQuantity int
	ProductPrice    float64
}

type GetResponse struct {
	Id           int
	Created_Date time.Time
	ExternalId   string
	Status       string
	OrderId      int
	SubTotal     float64
	LineItems    []Item
	TaxAmount    float64
}

type GetAllResponse struct {
	List []GetResponse
}

type CreateRequest struct {
	ExternalId string
	Status     string
	OrderId    int
	Subtotal   float64
	TaxAmount  float64
	LineItems  []Item
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
	GetByOrderId(ctx context.Context, orderId int) (GetResponse, error)
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
