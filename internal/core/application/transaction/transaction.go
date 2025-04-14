package transaction

import (
	"context"
	"errors"
	"time"

	partner "b2b.nati011.github.com/internal/core/application/payment_partner"
	"b2b.nati011.github.com/internal/core/application/user"
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
	Id         int
	User_Id    int
	Date       time.Time
	Amount     int64
	Partner_Id int
}

type GetAllResponse struct {
	List []GetResponse
}

type CreateRequest struct {
	User_Id    int
	Amount     int64
	Partner_Id int
}

type Pagination struct {
	Limit  int
	Offset int
}

type GetByParamRequest struct {
	Date       time.Time
	Partner_Id int
	User_Id    int
}

type Provider interface {
	Create(context.Context, *CreateRequest) (int, error)
	Get(context.Context, int) (GetResponse, error)
	GetAll(context.Context, *Pagination) (GetAllResponse, error)
	GetByParam(context.Context, *GetByParamRequest, *Pagination) (GetAllResponse, error)
}

type TransactionService struct {
	DB             port.DB
	PartnerService partner.Provider
	UserService    user.Provider
}

func NewTransactionService(db port.DB,
	ps partner.Provider,
	us user.Provider,
) Provider {
	return &TransactionService{
		DB:             db,
		PartnerService: ps,
		UserService:    us,
	}
}

func (t *TransactionService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	err := t.validateUserId(ctx, req.User_Id)
	if err != nil {
		return 0, err
	}
	err = validateAmount(req.Amount)
	if err != nil {
		return 0, err
	}
	err = t.validatePartnerId(ctx, req.Partner_Id)
	if err != nil {
		return 0, err
	}

	id, err := t.DB.Create(ctx, &port.CreateRequest{
		User_Id:    req.User_Id,
		Partner_Id: req.Partner_Id,
		Amount:     req.Amount,
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

func (t *TransactionService) GetAll(ctx context.Context, pagination *Pagination) (GetAllResponse, error) {
	resp, err := t.DB.GetAll(ctx, &port.Pagination{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	})
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
	// if len(ret_resp.List) == 0 {
	// 	return GetAllResponse{}, ErrEmptyGetContent
	// }
	return ret_resp, nil
}

func (t *TransactionService) GetByParam(ctx context.Context, req *GetByParamRequest, pagination *Pagination) (GetAllResponse, error) {
	ret_resp := GetAllResponse{}
	if req.Partner_Id != 0 {

		resp, err := t.DB.GetByPartnerId(ctx, req.Partner_Id, &port.Pagination{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		})
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

	if req.User_Id != 0 {
		resp, err := t.DB.GetByUserId(ctx, req.User_Id, &port.Pagination{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		})
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

	if !req.Date.IsZero() {
		resp, err := t.DB.GetByDate(ctx, req.Date, &port.Pagination{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
		})
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
