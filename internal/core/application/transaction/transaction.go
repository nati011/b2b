package transaction

import (
	"context"
	"errors"
	"time"

	port "b2b.nati011.github.com/internal/port/application/transaction/db"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

var (
	ErrUserIdNotSupplied           = errors.New("¯\\_(o_o)_/¯, user id mandatory")
	ErrPartnerIdNotSupplied        = errors.New("¯\\_(o_o)_/¯, partnerId mandatory")
	ErrAmountIsNotSupplied         = errors.New("¯\\_(o_o)_/¯, amount mandatory")
	ErrAmountMustBeGreaterThanZero = errors.New("¯\\_(o_o)_/¯, amount must be greater than zero")
	ErrIdNotFound                  = errors.New("¯\\_(o_o)_/¯, id not found")
	ErrEmptyGetContent             = errors.New("¯\\_(o_o)_/¯, empty get content")
	ErrUnknown                     = errors.New("¯\\_(o_o)_/¯, unkown error")
	ErrPartnerDoesNotExist         = errors.New("¯\\_(o_o)_/¯, partner does not exist")
	ErrUserDoesNotExist            = errors.New("¯\\_(o_o)_/¯, user does not exist")
	ErrTransactionRefNotSupplied   = errors.New("¯\\_(o_o)_/¯, txRef mandatory")
)

var (
	PENDING_STATUS   = "PENDING"
	COMPLETED_STATUS = "COMPLETED"
)

type GetResponse struct {
	Id        int
	Date      time.Time
	Amount    float64
	PartnerId int
	TxRef     string
	Status    string
}

type UpdateRequest struct {
	Id     int
	Status string
}

type UpdateByTransactionRefRequest struct {
	TransactionRef string
	Status         string
}

type GetAllResponse struct {
	List []GetResponse
}

type CreateRequest struct {
	Amount    float64
	PartnerId int
	TxRef     string
	Status    string
}

type GetByParamRequest struct {
	Date      time.Time
	PartnerId int
	TxRef     string
	Status    string
}

type Provider interface {
	Create(context.Context, *CreateRequest) (int, error)
	Get(context.Context, int) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
	GetByParam(context.Context, *GetByParamRequest) (GetAllResponse, error)
	UpdateStatus(context.Context, *UpdateRequest) error
	UpdateByTransactionRef(context.Context, *UpdateByTransactionRefRequest) error
}

type TransactionService struct {
	DB port.DB
}

func NewTransactionService(db port.DB,
) Provider {
	return &TransactionService{
		DB: db,
	}
}

func (t *TransactionService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	err := validateAmount(req.Amount)
	if err != nil {
		return 0, err
	}
	err = validateTransactionRef(req.TxRef)
	if err != nil {
		return 0, err
	}
	err = validatePartnerId(req.PartnerId)
	if err != nil {
		return 0, err
	}
	id, err := t.DB.Create(ctx, &port.CreateRequest{
		PartnerId: req.PartnerId,
		Amount:    req.Amount,
		TxRef:     req.TxRef,
		Status:    req.Status,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}
	return id, nil
}

func (t *TransactionService) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := t.DB.GetByID(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	return GetResponse(resp), nil
}

func (t *TransactionService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp, err := t.DB.GetAll(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	ret_resp := GetAllResponse{}
	for _, i := range resp.List {
		ret_resp.List = append(ret_resp.List, GetResponse(i))
	}
	return ret_resp, nil
}

func (t *TransactionService) UpdateStatus(ctx context.Context, req *UpdateRequest) error {
	err := t.DB.UpdateStatus(ctx, &port.UpdateRequest{
		Id:     req.Id,
		Status: req.Status,
	})
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (t *TransactionService) UpdateByTransactionRef(ctx context.Context, req *UpdateByTransactionRefRequest) error {
	err := t.DB.UpdateByTransactionRef(ctx, &port.UpdateByTransactionRefRequest{
		TransactionRef: req.TransactionRef,
		Status:         req.Status,
	})
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (t *TransactionService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	ret_resp := GetAllResponse{}

	resp, err := t.DB.GetByPartnerId(ctx, req.PartnerId)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}

	for _, i := range resp.List {
		found := false
		for _, j := range ret_resp.List {
			if j.Id == i.Id {
				found = true
			}
		}
		if !found {
			ret_resp.List = append(ret_resp.List, GetResponse(i))
		}
	}

	if !req.Date.IsZero() {
		resp, err := t.DB.GetByDate(ctx, req.Date)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		for _, i := range resp.List {
			found := false
			for _, j := range ret_resp.List {
				if j.Id == i.Id {
					found = true
				}
			}
			if !found {
				ret_resp.List = append(ret_resp.List, GetResponse(i))
			}
		}
	}

	if req.TxRef != "" {
		resp, err := t.DB.GetByTxRef(ctx, req.TxRef)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		for _, i := range resp.List {
			found := false
			for _, j := range ret_resp.List {
				if j.Id == i.Id {
					found = true
				}
			}
			if !found {
				ret_resp.List = append(ret_resp.List, GetResponse(i))
			}
		}
	}

	if req.Status != "" {
		resp, err := t.DB.GetByStatus(ctx, req.Status)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		for _, i := range resp.List {
			found := false
			for _, j := range ret_resp.List {
				if j.Id == i.Id {
					found = true
				}
			}
			if !found {
				ret_resp.List = append(ret_resp.List, GetResponse(i))
			}
		}
	}

	if len(ret_resp.List) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}
	return ret_resp, nil
}
