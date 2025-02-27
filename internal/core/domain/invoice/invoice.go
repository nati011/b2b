package invoice

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSysStatusNotSupplied = errors.New("status not supplied")
	ErrSysIdNotFound        = errors.New("id not found")
	ErrSysEmptyGetContent   = errors.New("empty get content")
)

type CreateRequest struct {
	ExternalId string
	Status     string
}

type GetResponse struct {
	Id           int
	Created_Date time.Time
	ExternalId   string
	Status       string
}

type GetAllResponse struct {
	List []GetResponse
}

type UpdateByParamRequest struct {
	Id         int
	Status     string
	ExternalId string
}

type GetByParamRequest struct {
	Status       string
	ExternalId   string
	Created_Date time.Time
}

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	Update(ctx context.Context, req *UpdateByParamRequest) error
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
}

type InvoiceService struct {
}

func NewInvoice() Provider {
	return &InvoiceService{}
}

func (i *InvoiceService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	err := validateStatus(req.Status)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (i *InvoiceService) Update(ctx context.Context, req *UpdateByParamRequest) error {
	//validate
	return nil
}

func (i *InvoiceService) Get(ctx context.Context, id int) (GetResponse, error) {
	return GetResponse{}, nil
}
func (i *InvoiceService) GetAll(ctx context.Context) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

func (i *InvoiceService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}
