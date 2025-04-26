package transaction

import (
	"context"
	"errors"
	"time"

	port "b2b.nati011.github.com/internal/port/application/transaction/db"
)

var (
	ErrUserIdNotSupplied           = errors.New("oopsy, user id mandatory")
	ErrPartnerIdNotSupplied        = errors.New("oopsy, partnerId mandatory")
	ErrAmountIsNotSupplied         = errors.New("oopsy, amount mandatory")
	ErrAmountMustBeGreaterThanZero = errors.New("oopsy, amount must be greater than zero")
	ErrIdNotFound                  = errors.New("oopsy, id not found")
	ErrEmptyGetContent             = errors.New("oopsy, empty get content")
	ErrUnknown                     = errors.New("oopsy, unkown error")
	ErrPartnerDoesNotExist         = errors.New("oopsy, partner does not exist")
	ErrUserDoesNotExist            = errors.New("oopsy, user does not exist")
)

type GetResponse struct {
	Id        int
	Date      time.Time
	Amount    float64
	PartnerId int
	TxRef     string
	Status    string
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
		case port.ErrSysNoRows:
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
		case port.ErrSysNoRows:
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

func (t *TransactionService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	ret_resp := GetAllResponse{}
	resp, err := t.DB.GetByPartnerId(ctx, req.PartnerId)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
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
			case port.ErrSysNoRows:
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
			case port.ErrSysNoRows:
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
