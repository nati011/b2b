package distributor

import (
	"context"
	"errors"

	authPort "b2b.nati011.github.com/internal/port/application/auth/provider"
	port "b2b.nati011.github.com/internal/port/distributor"
)

var (
	ErrEmptyGetContent = errors.New("oopsy, no distributor found")
	ErrUnknown         = errors.New("oopsy, unkown error")
	SUCCESS_MESSAGE    = "Ahoy!"
)

type Provider interface {
	Create(ctx context.Context, req *port.RegisterDistributorRequest) (response port.RegisterDistributorResponse, err error)
	GetAll(ctx context.Context) (port.GetAllResponse, error)
	GetByParam(ctx context.Context, req *port.GetByParamRequest) (port.GetResponse, error)
}

type DistributorService struct {
	db          port.DB
	authService authPort.Provider
}

func (d *DistributorService) Create(ctx context.Context, req *port.RegisterDistributorRequest) (response port.RegisterDistributorResponse, err error) {
	resp := port.RegisterDistributorResponse{}
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

	resp = port.RegisterDistributorResponse{
		Message:       SUCCESS_MESSAGE,
		DistributorId: id,
	}
	return resp, nil
}

func (d *DistributorService) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp, err := d.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return port.GetAllResponse{}, ErrEmptyGetContent
		default:
			return port.GetAllResponse{}, ErrUnknown
		}
	}

	resp_val := port.GetAllResponse{}
	for _, i := range resp.List {
		resp_val.List = append(resp_val.List, port.GetResponse(i))
	}
	return resp_val, nil
}

func (d *DistributorService) GetByParam(ctx context.Context, req *port.GetByParamRequest) (port.GetResponse, error) {
	resp, err := d.db.GetById(ctx, req.Id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return port.GetResponse{}, ErrEmptyGetContent
		default:
			return port.GetResponse{}, ErrUnknown
		}
	}

	resp_val := port.GetResponse{
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
