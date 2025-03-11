package distributor

import (
	"context"
	"errors"

	authPort "b2b.nati011.github.com/internal/port/application/auth/provider"
	port "b2b.nati011.github.com/internal/port/distributor"
)

var (
	ErrEmptyGetDistributorContent = errors.New("oopsy, no distributor found")
	ErrEmptyGetBusinessContent    = errors.New("oopsy, no business found")
	ErrUnknown                    = errors.New("oopsy, unkown error")
	SUCCESS_MESSAGE               = "Ahoy!"
)

type Provider interface {
	Create(ctx context.Context, req *port.RegisterDistributorRequest) (response port.RegisterDistributorResponse, err error)
	GetAll(ctx context.Context) (port.GetAllResponse, error)
	GetByParam(ctx context.Context, req *port.GetByParamRequest) (port.GetResponse, error)
	AddBusinessInformattion(ctx context.Context, req *port.CreateBusinessInformation) (response port.CreateBusinessResponse, err error)
	GetBusinessAll(ctx context.Context) (port.GetAllResponse, error)
	GetByDistributor(ctx context.Context, distributorId int) (port.GetBusinessResponse, error)
	GetById(ctx context.Context, id int) (port.GetResponse, error)
	UpdateBusiness(ctx context.Context, req *port.UpdateBusinessRequest) (response port.RegisterDistributorResponse, err error)
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
		Message: SUCCESS_MESSAGE,
		Id:      id,
	}
	return resp, nil
}

func (d *DistributorService) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp, err := d.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return port.GetAllResponse{}, ErrEmptyGetDistributorContent
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
			return port.GetResponse{}, ErrEmptyGetDistributorContent
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
func (d *DistributorService) AddBusinessInformattion(ctx context.Context, req *port.CreateBusinessInformation) (resp port.CreateBusinessResponse, err error) {
	resp, err = d.db.CreateBusiness(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (d *DistributorService) GetBusinessAll(ctx context.Context) (resp port.GetAllResponse, err error) {
	resp, err = d.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return resp, ErrEmptyGetBusinessContent
		default:
			return resp, ErrUnknown
		}
	}
	return resp, nil
}

func (d *DistributorService) GetByDistributor(ctx context.Context, distributorId int) (resp port.GetBusinessResponse, err error) {
	resp, err = d.db.GetByDistributorId(ctx, distributorId)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return resp, ErrEmptyGetDistributorContent
		default:
			return resp, ErrUnknown
		}
	}
	return resp, nil
}
func (d *DistributorService) GetById(ctx context.Context, id int) (resp port.GetResponse, err error) {
	resp, err = d.db.GetById(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return resp, ErrEmptyGetDistributorContent
		default:
			return resp, ErrUnknown
		}
	}
	return resp, nil
}

func (d *DistributorService) UpdateBusiness(ctx context.Context, req *port.UpdateBusinessRequest) (response port.RegisterDistributorResponse, err error) {
	_, err = d.db.GetBusinessById(
		ctx,
		req.Id,
	)

	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return response, ErrEmptyGetDistributorContent

		default:
			return response, ErrUnknown
		}
	}

	id, err := d.db.UpdateBusiness(ctx, req)
	if err != nil {

		return response, err

	}

	response = port.RegisterDistributorResponse{
		Message: SUCCESS_MESSAGE,
		Id:      id,
	}
	return response, nil
}
func NewDistributorService(db port.DB, authService authPort.Provider) Provider {
	return &DistributorService{
		db:          db,
		authService: authService,
	}
}
