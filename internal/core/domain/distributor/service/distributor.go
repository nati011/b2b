package distributor

import (
	"context"
	"errors"

	distributorDTO "b2b.nati011.github.com/internal/core/domain/distributor/model/dto"
	port "b2b.nati011.github.com/internal/port/distributor"
)

var (
	ErrEmptyGetContent = errors.New("oopsy, no distributor found")
	ErrUnknown         = errors.New("oopsy, unkown error")
)

type Provider interface {
	Create(ctx context.Context, req *distributorDTO.RegisterDistributorRequest) (id int, err error)
	CreateBusinessInformation(ctx context.Context, req *distributorDTO.UpdateBusinessRequest) (id int, err error)
	Get(ctx context.Context, id int) (distributorDTO.GetResponse, error)
	GetAll(ctx context.Context) (distributorDTO.GetAllResponse, error)
	GetByParam(ctx context.Context, req *distributorDTO.GetByParamRequest) (distributorDTO.GetAllResponse, error)
}

type DistributorService struct {
	db port.DB
}

func (d *DistributorService) CreateBusinessInformation(ctx context.Context, req *distributorDTO.UpdateBusinessRequest) (id int, err error) {
	panic("unimplemented")
}

func (d *DistributorService) Create(ctx context.Context, req *distributorDTO.RegisterDistributorRequest) (id int, err error) {
	id, err = d.db.Create(ctx, &port.CreateRequest{
		FullName:   req.FullName,
		Email:      req.Email,
		Phone:      req.Password,
		Password:   req.Password,
		DOB:        req.DOB,
		IsActive:   false,
		Username:   req.Username,
		ExternalId: req.ExternalId,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	print(id)
	print("_____________________________")
	return id, nil
}

func (d *DistributorService) Get(ctx context.Context, id int) (distributorDTO.GetResponse, error) {
	panic("unimplemented")
}

func (d *DistributorService) GetAll(ctx context.Context) (distributorDTO.GetAllResponse, error) {
	resp, err := d.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return distributorDTO.GetAllResponse{}, ErrEmptyGetContent
		default:
			return distributorDTO.GetAllResponse{}, ErrUnknown
		}
	}

	resp_val := distributorDTO.GetAllResponse{}
	for _, i := range resp.List {
		resp_val.List = append(resp_val.List, distributorDTO.GetResponse(i))
	}
	return resp_val, nil
}

func (d *DistributorService) GetByParam(ctx context.Context, req *distributorDTO.GetByParamRequest) (distributorDTO.GetAllResponse, error) {
	panic("unimplemented")
}

func NewDistributorService(db port.DB) Provider {
	return &DistributorService{
		db: db,
	}
}
