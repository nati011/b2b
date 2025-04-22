package payment_partner

import (
	"context"
	"errors"
	"log"

	port "b2b.nati011.github.com/internal/port/application/partner/db"
)

const (
	CALLBACK_URL = "/payment_call_back"

	INACTIVE_STATUS = "INACTIVE"
	ACTIVE_STATUS   = "ACTIVE"
)

var (
	ErrNameIsNotSupplied            = errors.New("oopsy, name is not supplied")
	ErrIconIsNotSupplied            = errors.New("oopsy, icon is not supplied")
	ErrSecretIsNotSupplied          = errors.New("oppsy, secret is not supplied")
	ErrUrlIsNotSupplied             = errors.New("oopsy, init payment url is not supplied")
	ErrIdNotFound                   = errors.New("oopsy, id not found")
	ErrPaymentOptionaAlreadyActive  = errors.New("oopsy, payment option is already active")
	ErrPaymentOptionAlreadyInactive = errors.New("oopsy, payment option is already inactive")
	ErrEmptyGetContent              = errors.New("oopsy, empty get content")
	ErrUnknown                      = errors.New("oopsy, unknown error has occured")
)

type CreateRequest struct {
	Name    string
	Icon    string
	Status  string
	BaseUrl string
	Secret  string
}

type GetResponse struct {
	Id      int
	Name    string
	Icon    string
	Status  string
	BaseUrl string
}

type GetSecretResponse struct {
	Name    string
	BaseUrl string
	Secret  string
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByParamRequest struct {
	Name   string
	Status string
}

type Provider interface {
	Create(context.Context, *CreateRequest) (int, error)
	Get(context.Context, int) (GetResponse, error)
	GetPartnerSecret(context.Context, int) (GetSecretResponse, error)
	Activate(context.Context, int) error
	Deactivate(context.Context, int) error
	GetAll(context.Context) (GetAllResponse, error)
	GetActive(context.Context) (GetAllResponse, error)
	GetByParam(context.Context, *GetByParamRequest) (GetAllResponse, error)
}

type PartnerService struct {
	DB port.DB
}

func NewPartner(db port.DB) Provider {
	return &PartnerService{
		DB: db,
	}
}

func (p *PartnerService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	err := validateName(req.Name)
	if err != nil {
		return 0, err
	}
	err = validateIcon(req.Icon)
	if err != nil {
		return 0, err
	}
	err = validateBaseURL(req.BaseUrl)
	if err != nil {
		return 0, err
	}

	err = validateSecret(req.Secret)
	if err != nil {
		return 0, err
	}

	id, err := p.DB.Create(ctx, &port.CreateRequest{
		Name:    req.Name,
		Icon:    req.Icon,
		Status:  INACTIVE_STATUS,
		BaseUrl: req.BaseUrl,
		Secret:  req.Secret,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}
	log.Printf("Created id %v", id)
	return id, nil
}

func (p *PartnerService) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := p.DB.GetByID(ctx, id)
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

func (p *PartnerService) GetPartnerSecret(ctx context.Context, id int) (GetSecretResponse, error) {
	resp, err := p.DB.GetPartnerSecret(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetSecretResponse{}, ErrIdNotFound
		default:
			return GetSecretResponse{}, ErrUnknown
		}
	}
	return GetSecretResponse(resp), nil
}

func (p *PartnerService) Activate(ctx context.Context, id int) error {
	//check if id exists
	got, err := p.Get(ctx, id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return err
		default:
			return ErrUnknown
		}
	}
	//check if already active
	if got.Status == ACTIVE_STATUS {
		return ErrPaymentOptionaAlreadyActive
	}

	_, err = p.DB.UpdateStatus(ctx, id, ACTIVE_STATUS)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (p *PartnerService) Deactivate(ctx context.Context, id int) error {
	//check if id exists
	got, err := p.Get(ctx, id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return err
		default:
			return ErrUnknown
		}
	}
	//check if already active
	if got.Status == INACTIVE_STATUS {
		return ErrPaymentOptionAlreadyInactive
	}

	_, err = p.DB.UpdateStatus(ctx, id, INACTIVE_STATUS)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (p *PartnerService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp, err := p.DB.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	var resp_val GetAllResponse
	for _, i := range resp.List {
		resp_val.List = append(resp_val.List, GetResponse(i))
	}
	return resp_val, nil
}

func (p *PartnerService) GetActive(ctx context.Context) (GetAllResponse, error) {
	resp, err := p.DB.GetByStatus(ctx, ACTIVE_STATUS)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	var resp_val GetAllResponse
	for _, i := range resp.List {
		resp_val.List = append(resp_val.List, GetResponse(i))
	}
	return resp_val, nil
}

func (p *PartnerService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	var resp GetAllResponse
	if req.Name != "" {
		resp_name, err := p.DB.GetByName(ctx, req.Name)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range resp_name.List {
			resp.List = append(resp.List, GetResponse(i))
		}
	}

	if req.Status != "" {
		resp_name, err := p.DB.GetByStatus(ctx, req.Status)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		for _, i := range resp_name.List {
			resp.List = append(resp.List, GetResponse(i))
		}
	}
	if len(resp.List) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}
	return resp, nil
}
