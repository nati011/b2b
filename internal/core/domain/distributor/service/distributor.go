package distributor

import (
	"context"
	"errors"

	distributorDTO "b2b.nati011.github.com/internal/core/domain/distributor/model/dto"
	authPort "b2b.nati011.github.com/internal/port/application/auth/provider"
	port "b2b.nati011.github.com/internal/port/distributor"
)

var (
	ErrEmptyGetContent = errors.New("oopsy, no distributor found")
	ErrUnknown         = errors.New("oopsy, unkown error")
	SUCCESS_MESSAGE    = "Ahoy!"
)

type Provider interface {
	Create(ctx context.Context, req *distributorDTO.RegisterDistributorRequest) (response distributorDTO.RegisterDistributorResponse, err error)
	GetAll(ctx context.Context) (distributorDTO.GetAllResponse, error)
	GetByParam(ctx context.Context, req *distributorDTO.GetByParamRequest) (distributorDTO.GetResponse, error)
}

type DistributorService struct {
	db          port.DB
	authService authPort.Provider
}

func (d *DistributorService) Create(ctx context.Context, req *distributorDTO.RegisterDistributorRequest) (response distributorDTO.RegisterDistributorResponse, err error) {
	resp := distributorDTO.RegisterDistributorResponse{}
	distribtor := port.CreateRequest{
		FirstName:  req.FirstName,
		Email:      req.Email,
		DOB:        req.DOB,
		Username:   req.Username,
		ExternalId: req.ExternalId,
	}

	id, err := d.db.Create(ctx, &distribtor)
	if err != nil {

		return resp, err

	}

	user := authPort.RegisterUserRequest{
		FirstName:   req.FirstName,
		Email:       req.Email,
		BirthDate:   req.DOB,
		PhoneNumber: req.PhoneNumber,
		Username:    req.Username,
		ExternalId:  req.ExternalId,
	}
	_, err = d.authService.CreateNewClient(ctx, user)
	if err != nil {
		return resp, err

	}

	resp = distributorDTO.RegisterDistributorResponse{
		Message:       SUCCESS_MESSAGE,
		DistributorId: id,
	}
	return resp, nil
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

func (d *DistributorService) GetByParam(ctx context.Context, req *distributorDTO.GetByParamRequest) (distributorDTO.GetResponse, error) {
	resp, err := d.db.GetById(ctx, req.Id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return distributorDTO.GetResponse{}, ErrEmptyGetContent
		default:
			return distributorDTO.GetResponse{}, ErrUnknown
		}
	}

	resp_val := distributorDTO.GetResponse{
		Id:        resp.Id,
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
	}

	return resp_val, nil
}

func NewDistributorService(db port.DB, authService authPort.Provider) Provider {
	return &DistributorService{
		db:          db,
		authService: authService,
	}
}
