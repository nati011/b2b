package distributor

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/distributor"
)

var (
	ErrEmptyGetContent = errors.New("oopsy, no distributor found")
	ErrUnknown         = errors.New("oopsy, unkown error")
)

type RegisterDistributorRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmed_password"`
	FullName        string `json:"full_name"`
	Username        string `json:"username"`
}

type BusinessLocation struct {
	GeneralZone string `json:"general_zone"`
	Region      string `json:"region"`
	Woreda      string `json:"woreda"`
}

type UpdateBusinessRequest struct {
	DistributorId int              `json:"distributorId"`
	Name          string           `json:"name"`
	Tin           int              `json:"tin"`
	Region        BusinessLocation `json:"region"`
}

type RegisterDistributorResponse struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type GetResponse struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmed_password"`
	FullName        string `json:"full_name"`
	Username        string `json:"username"`
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByParamRequest struct {
	Name  string
	Email string
}

type Provider interface {
	Create(ctx context.Context, req *RegisterDistributorRequest) (id int, err error)
	CreateBusinessInformation(ctx context.Context, req *UpdateBusinessRequest) (id int, err error)
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
}

type DistributorService struct {
	DB port.DB
}

func (d *DistributorService) CreateBusinessInformation(ctx context.Context, req *UpdateBusinessRequest) (id int, err error) {
	panic("unimplemented")
}

func (d *DistributorService) Create(ctx context.Context, req *RegisterDistributorRequest) (id int, err error) {
	panic("unimplemented")
}

func (d *DistributorService) Get(ctx context.Context, id int) (GetResponse, error) {
	panic("unimplemented")
}

func (d *DistributorService) GetAll(ctx context.Context) (GetAllResponse, error) {
	panic("unimplemented")
}

func (d *DistributorService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	panic("unimplemented")
}

func NewDistributorService(db port.DB) Provider {
	return &DistributorService{
		DB: db,
	}
}
