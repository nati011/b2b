package payment

import (
	"context"
)

type CreateRequest struct {
	Name          string
	Icon          string
	Status        string
	BaseURL       string
	Secret        string
	PaymentMethod string
}

type GetResponse struct {
	Id            int
	Name          string
	Icon          string
	Status        string
	BaseURL       string
	PaymentMethod string
}

type GetPartnerSecret struct {
	Name    string
	BaseURL string
	Secret  string
}

type UpdatePartnerSecret struct {
	Id      int
	BaseURL string
	Secret  string
}

type GetAllResponse struct {
	List []GetResponse
}

type Reader interface {
	GetPartnerSecret(context.Context, int) (GetPartnerSecret, error)
	UpdatePartnerSecret(context.Context, UpdatePartnerSecret) error
	GetByID(context.Context, int) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
	GetByStatus(context.Context, string) (GetAllResponse, error)
	GetByName(context.Context, string) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
	UpdateStatus(context.Context, int, string) (int, error)
	UpdateName(context.Context, int, string) (int, error)
	Delete(context.Context, int) error
}

type DB interface {
	Reader
	Writer
}
