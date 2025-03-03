package partner

import (
	"context"
	"errors"
)

const (
	CALLBACK_URL = "/payment_call_back"

	INACTIVE_STATUS = "INACTIVE"
	ACTIVE_STATUS   = "ACTIVE"
)

var (
	ErrNameIsNotSupplied            = errors.New("oopsy, name is not supplied")
	ErrIconIsNotSupplied            = errors.New("oopsy, icon is not supplied")
	ErrUrlIsNotSupplied             = errors.New("oopsy, init payment url is not supplied")
	ErrIdNotFound                   = errors.New("oopsy, id not found")
	ErrPaymentOptionaAlreadyActive  = errors.New("oopsy, payment option already active")
	ErrPaymentOptionAlreadyInactive = errors.New("oopsy, payment option alreadt active")
	ErrEmptyGetContent              = errors.New("oopsy, empty get content")
)

type CreateRequest struct {
	Name             string
	Icon             string
	Init_payment_url string
}

type GetResponse struct {
	Id     int
	Name   string
	Icon   string
	Status string
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByParamRequest struct {
	Name   string
	Status string
}

type CheckoutResponse struct {
	Checkout_url string
}

type Provider interface {
	Create(context.Context, *CreateRequest) (int, error)
	Get(context.Context, int) (GetResponse, error)
	Activate(context.Context, int) error
	Deactivate(context.Context, int) error
	GetAll(context.Context) (GetAllResponse, error)
	GetActive(context.Context) (GetAllResponse, error)
	GetByParam(context.Context, *GetByParamRequest) (GetAllResponse, error)
}

type PartnerService struct {
}

func NewPartner() Provider {
	return &PartnerService{}
}

func (p *PartnerService) Create(context.Context, *CreateRequest) (int, error) {
	return 0, nil
}

func (p *PartnerService) Get(context.Context, int) (GetResponse, error) {
	return GetResponse{}, nil
}

func (p *PartnerService) Activate(context.Context, int) error {
	return nil
}

func (p *PartnerService) Deactivate(context.Context, int) error {
	return nil
}

func (p *PartnerService) GetAll(context.Context) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

func (p *PartnerService) GetActive(context.Context) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}

func (p *PartnerService) GetByParam(context.Context, *GetByParamRequest) (GetAllResponse, error) {
	return GetAllResponse{}, nil
}
