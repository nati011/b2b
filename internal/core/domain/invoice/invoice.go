package invoice

import (
	"context"
	"errors"
	"time"

	port "b2b.nati011.github.com/internal/port/invoice"
)

var (
	ErrSysStatusNotSupplied = errors.New("status not supplied")
	ErrSysIdNotFound        = errors.New("id not found")
	ErrSysEmptyGetContent   = errors.New("empty get content")
	ErrSysUnknown           = errors.New("unknown error")
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
	DB port.DB
}

func NewInvoice(db port.DB) Provider {
	return &InvoiceService{
		DB: db,
	}
}

func (i *InvoiceService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	err := validateStatus(req.Status)
	if err != nil {
		return 0, err
	}

	id, err := i.DB.Create(ctx, &port.CreateRequest{
		ExternalId: req.ExternalId,
		Status:     req.Status,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrSysUnknown
		}
	}
	return id, nil
}

func (i *InvoiceService) Update(ctx context.Context, req *UpdateByParamRequest) error {
	//validate
	_, err := i.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case ErrSysIdNotFound:
			return ErrSysIdNotFound
		default:
			return ErrSysUnknown
		}
	}

	if req.ExternalId != "" {
		err := i.DB.UpdateExternalId(ctx, &port.UpdateExternalIdRequest{
			Id:         req.Id,
			ExternalId: req.ExternalId,
		})
		if err != nil {
			switch err {
			default:
				return ErrSysUnknown
			}
		}
	}

	if req.Status != "" {
		err = validateStatus(req.Status)
		if err != nil {
			return err
		}

		err := i.DB.UpdateStatus(ctx, &port.UpdateStatusRequest{
			Id:     req.Id,
			Status: req.Status,
		})
		if err != nil {
			switch err {
			default:
				return ErrSysUnknown
			}
		}
	}

	return nil
}

func (i *InvoiceService) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := i.DB.Get(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetResponse{}, ErrSysIdNotFound
		default:
			return GetResponse{}, ErrSysUnknown
		}
	}

	return GetResponse(resp), nil
}
func (i *InvoiceService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp := []GetResponse{}
	db_resp, err := i.DB.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetAllResponse{}, ErrSysEmptyGetContent
		default:
			return GetAllResponse{}, ErrSysUnknown
		}
	}
	for _, i := range db_resp.List {
		resp = append(resp, GetResponse(i))
	}
	return GetAllResponse{
		List: resp,
	}, nil
}

func (i *InvoiceService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	resp := []GetResponse{}

	if req.ExternalId != "" {
		db_resp, err := i.DB.GetByExternalId(ctx, req.ExternalId)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrSysUnknown
			}
		}
		for _, i := range db_resp.List {
			resp = append(resp, GetResponse(i))
		}
	}

	if req.Status != "" {
		db_resp, err := i.DB.GetByStatus(ctx, req.Status)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrSysUnknown
			}
		}
		for _, i := range db_resp.List {
			resp = append(resp, GetResponse(i))
		}
	}
	if len(resp) == 0 {
		return GetAllResponse{}, ErrSysEmptyGetContent
	}
	return GetAllResponse{
		List: resp,
	}, nil
}
