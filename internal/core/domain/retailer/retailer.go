package retailer

import (
	"context"
	"errors"

	"b2b.nati011.github.com/internal/core/application/user"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/retailer"
)

var (
	ErrUnknown            = errors.New("¯\\_(•_•)_/¯, unknown error")
	ErrInvalidTin         = errors.New("¯\\_(•_•)_/¯, tin invalid")
	ErrDuplicateTin       = errors.New("¯\\_(•_•)_/¯, tin already in use")
	ErrIdNotFound         = errors.New("¯\\_(•_•)_/¯, id not found")
	ErrEmptyGetContent    = errors.New("¯\\_(•_•)_/¯, empty get content")
	ErrRetailerHasNoUsers = errors.New("oopys, retailer has no users")
	ErrPhoneMandatory     = errors.New("oopys, phonenumber mandatory")
)

type CreateRequest struct {
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string

	FirstName string
	LastName  string
	Email     string
	Phone     string
	Username  string
	Password  string
}

type GetResponse struct {
	Id          int
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByParamRequest struct {
	Name string
	Tin  string
}

type UpdateRequest struct {
	Id   int
	Name string
	Tin  string
}

type GetAllUsers struct {
	List []int
}

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	Get(ctx context.Context, id int) (GetResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	Update(ctx context.Context, req *UpdateRequest) (int, error)
	GetAllUsers(ctx context.Context, id int) (GetAllUsers, error)
}

type RetailerService struct {
	DB          port.DB
	UserService user.Provider
}

func NewRetailerService(up user.Provider, db port.DB) Provider {
	return &RetailerService{
		UserService: up,
		DB:          db,
	}
}

func (r *RetailerService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	//validate
	if req.Tin != "" {
		err := r.validateTin(ctx, req.Tin)
		if err != nil {
			return 0, err
		}
	}

	err := validatePhone(req.Phone)
	if err != nil {
		return 0, err
	}

	// create user
	user_id, err := r.UserService.Create(ctx, &user.CreateRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Username:  req.Username,
		Phone:     req.Phone,
		Password:  req.Password,
	})
	if err != nil {
		switch err {
		case user.ErrUnknown:
			return 0, ErrUnknown
		default:
			return 0, err
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
			r.UserService.Remove(ctx, user_id)
			return 0, ErrUnknown
		}
	}

	return id, nil
}

func (r *RetailerService) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := r.DB.Get(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}

	return GetResponse{
		Id:          resp.Id,
		Name:        resp.Name,
		Tin:         resp.Tin,
		Latitude:    resp.Latitude,
		Longitude:   resp.Longitude,
		GeneralZone: resp.GeneralZone,
		Region:      resp.Region,
		Woreda:      resp.Woreda,
	}, nil
}

func (r *RetailerService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	resp := port.GetAllResponse{}

	if req.Name != "" {
		resp_name, err := r.DB.GetByName(ctx, req.Name)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		if len(resp_name.List) != 0 {
			resp.List = append(resp.List, resp_name.List...)
		}
	}

	if req.Tin != "" {
		resp_tin, err := r.DB.GetByTin(ctx, req.Tin)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		if resp_tin.Id != 0 {
			resp.List = append(resp.List, resp_tin)
		}
	}
	if len(resp.List) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}
	service_resp := GetAllResponse{}
	for _, i := range resp.List {
		service_resp.List = append(service_resp.List, GetResponse{
			Id:          i.Id,
			Name:        i.Name,
			Tin:         i.Tin,
			Latitude:    i.Latitude,
			Longitude:   i.Longitude,
			GeneralZone: i.GeneralZone,
			Region:      i.Region,
			Woreda:      i.Woreda,
		})
	}
	if len(service_resp.List) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}
	return service_resp, nil
}

func (r *RetailerService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp := port.GetAllResponse{}

	resp_name, err := r.DB.GetAll(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	if len(resp_name.List) != 0 {
		resp.List = append(resp.List, resp_name.List...)
	}
	service_resp := GetAllResponse{}
	for _, i := range resp.List {
		service_resp.List = append(service_resp.List, GetResponse{
			Id:          i.Id,
			Name:        i.Name,
			Tin:         i.Tin,
			Latitude:    i.Latitude,
			Longitude:   i.Longitude,
			GeneralZone: i.GeneralZone,
			Region:      i.Region,
			Woreda:      i.Woreda,
		})
	}
	if len(service_resp.List) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}
	return service_resp, nil
}

func (r *RetailerService) Update(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate id
	_, err := r.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return 0, err
		default:
			return 0, ErrUnknown
		}
	}

	if req.Name != "" {
		err = r.DB.UpdateName(ctx, &port.UpdateNameRequest{
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
		err = r.validateTin(ctx, req.Tin)
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

		err = r.DB.UpdateTin(ctx, &port.UpdateTinRequest{
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

	return req.Id, nil
}

func (r *RetailerService) GetAllUsers(ctx context.Context, id int) (GetAllUsers, error) {
	var response_ids = []int{}
	users, err := r.DB.GetAllUserAgents(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllUsers{}, ErrRetailerHasNoUsers
		default:
			return GetAllUsers{}, ErrUnknown
		}
	}
	for _, i := range users.List {
		response_ids = append(response_ids, i.Id)
	}
	return GetAllUsers{
		List: response_ids,
	}, nil
}
