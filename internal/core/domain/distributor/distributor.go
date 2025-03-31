package distributor

import (
	"context"
	"errors"

	"b2b.nati011.github.com/internal/core/application/auth"
	user "b2b.nati011.github.com/internal/core/application/user"
	port "b2b.nati011.github.com/internal/port/application/distributor"
)

var (
	ErrEmptyGetDistributorContent = errors.New("oopsy, no distributor found")
	ErrEmptyGetBusinessContent    = errors.New("oopsy, no business found")
	ErrUnknown                    = errors.New("oopsy, unkown error")
	SUCCESS_MESSAGE               = "Ahoy!"
	ErrInvalidTin                 = errors.New("oopsy, tin invalid")
	ErrDuplicateTin               = errors.New("oopsy, tin already in use")
	ErrIdNotFound                 = errors.New("oopsy, id not found")
	ErrEmptyGetContent            = errors.New("oopsy, empty get content")
)

type UpdateRequest struct {
	Id   int
	Name string
	Tin  string
}
type GetResponse struct {
	Id          int
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Woreda      string
	UserId      int
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByIdRequest struct {
	Id int
}

type CreateRequest struct {
}

type BusinessLocation struct {
	GeneralZone string
	Region      string
	Woreda      string
}

type UpdateBusinessRequest struct {
	Id            int
	DistributorId int
	Name          string
	Tin           string
	Location      BusinessLocation
}

type RegisterDistributorResponse struct {
	Id      int
	Message string
}

type GetByParamRequest struct {
	Id    int
	Name  string
	Email string
	Tin   string
}

type GetBusinessResponse struct {
	Id            int
	Name          string
	Tin           string
	DistributorId int
}

type CreateBusinessInformation struct {
	Name          string
	Tin           string
	Latitude      string
	Longitude     string
	GeneralZone   string
	Region        string
	Woreda        string
	DistributorId int
}

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (response RegisterDistributorResponse, err error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetResponse, error)
	GetBusinessByDistributor(ctx context.Context, distributorId int) (GetBusinessResponse, error)
	GetById(ctx context.Context, id int) (GetResponse, error)
	Update(ctx context.Context, req *UpdateBusinessRequest) (id int, err error)
}

type DistributorService struct {
	db port.DB
}

func (d *DistributorService) Create(ctx context.Context, req *CreateRequest) (resp RegisterDistributorResponse, err error) {
	//validate
	err = d.validateTin(ctx, req.Tin)
	if err != nil {
		return 0, err
	}

	// create user
	user_id, err := d.UserService.Create(ctx, &user.CreateRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	})
	if err != nil {
		switch err {
		case user.ErrEmailNotValid,
			user.ErrPhoneNotValid,
			user.ErrPhoneOrEmailMandatory,
			user.ErrFirstNameMandatory:

			return 0, err
		default:
			return 0, ErrUnknown
		}
	}

	// create retailer
	id, err := r.DB.Create(ctx, port.CreateRequest{
		Name:        req.FirstName + req.LastName,
		Tin:         req.Tin,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		GeneralZone: req.GeneralZone,
		Region:      req.Region,
		Woreda:      req.Woreda,
		UserId:      user_id,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	return id, nil
}

func (d *DistributorService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp, err := d.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetDistributorContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}

	resp_val := GetAllResponse{}
	for _, i := range resp.List {
		resp_val.List = append(resp_val.List, GetResponse(i))
	}
	return resp_val, nil
}

func (d *DistributorService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetResponse, error) {
	resp, err := d.db.GetById(ctx, req.Id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetResponse{}, ErrEmptyGetDistributorContent
		default:
			return GetResponse{}, ErrUnknown
		}
	}

	resp_val := GetResponse{
		Id:        resp.Id,
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
	}

	return resp_val, nil
}

func (d *DistributorService) GetBusinessByDistributor(ctx context.Context, distributorId int) (resp GetBusinessResponse, err error) {
	business, err := d.db.GetBusinessByDistributorId(ctx, distributorId)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return resp, ErrEmptyGetDistributorContent
		default:
			return resp, ErrUnknown
		}
	}

	resp = GetBusinessResponse{
		Id:            business.Id,
		Name:          business.Name,
		Tin:           business.Tin,
		DistributorId: business.DistributorId,
	}
	return resp, nil
}
func (d *DistributorService) GetById(ctx context.Context, id int) (resp GetResponse, err error) {
	distributor, err := d.db.GetById(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return resp, ErrEmptyGetDistributorContent
		default:
			return resp, ErrUnknown
		}
	}

	resp = GetResponse{
		Id:         distributor.Id,
		Name:       distributor.FirstName,
		LastName:   distributor.LastName,
		Email:      distributor.Email,
		Phone:      distributor.Phone,
		Username:   distributor.Username,
		DOB:        distributor.DOB,
		ExternalId: distributor.ExternalId,
	}
	return resp, nil
}

func (d *DistributorService) Update(ctx context.Context, req *UpdateBusinessRequest) (int, error) {
	_, err := d.db.GetById(
		ctx,
		req.Id,
	)

	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return 0, ErrEmptyGetDistributorContent

		default:
			return 0, ErrUnknown
		}
	}
	if req.Name != "" {
		err = d.db.UpdateName(ctx, &port.UpdateNameRequest{
			Id:   req.Id,
			Name: req.Name,
		})
		if err != nil {
			switch err {
			default:
				return 0, ErrUnknown
			}
		}
	}

	if req.Tin != "" {
		err = d.validateTin(ctx, req.Tin)
		if err != nil {
			switch err {
			case ErrInvalidTin:
				return 0, err
			case ErrDuplicateTin:
				return 0, err
			default:
				return 0, ErrUnknown
			}
		}

		err = d.db.UpdateTin(ctx, &port.UpdateTinRequest{
			Id:  req.Id,
			Tin: req.Tin,
		})
		if err != nil {
			switch err {
			default:
				return 0, ErrUnknown
			}
		}
	}

	return id, nil
}
func NewDistributorService(db port.DB, authService auth.Provider) Provider {
	return &DistributorService{
		db:          db,
		authService: authService,
	}
}
