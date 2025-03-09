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
	Create(ctx context.Context, req *distributorDTO.RegisterDistributorRequest) (response string, err error)
	CreateBusinessInformation(ctx context.Context, req *distributorDTO.UpdateBusinessRequest) (id int, err error)
	Get(ctx context.Context, id int) (distributorDTO.GetResponse, error)
	GetAll(ctx context.Context) (distributorDTO.GetAllResponse, error)
	GetByParam(ctx context.Context, req *distributorDTO.GetByParamRequest) (distributorDTO.GetAllResponse, error)
}

type DistributorService struct {
	db          port.DB
	authService authPort.Provider
}

func (d *DistributorService) CreateBusinessInformation(ctx context.Context, req *distributorDTO.UpdateBusinessRequest) (id int, err error) {
	panic("unimplemented")
}

func (d *DistributorService) Create(ctx context.Context, req *distributorDTO.RegisterDistributorRequest) (response string, err error) {
	distribtor := port.CreateRequest{
		FirstName:  req.FirstName,
		Email:      req.Email,
		DOB:        req.DOB,
		Username:   req.Username,
		ExternalId: req.ExternalId,
	}
	_, err = d.db.Create(ctx, &distribtor)
	if err != nil {

		return "", err

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
		return "", err

	}
	return SUCCESS_MESSAGE, nil
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

func NewDistributorService(db port.DB, authService authPort.Provider) Provider {
	return &DistributorService{
		db:          db,
		authService: authService,
	}
}
