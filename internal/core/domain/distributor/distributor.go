package distributor

import (
	"context"
	"errors"

<<<<<<< HEAD
	"b2b.nati011.github.com/internal/core/application/user"
	port "b2b.nati011.github.com/internal/port/domain/distributor"
)

var (
	ErrUnknown         = errors.New("oopsy, unknown error")
	ErrInvalidTin      = errors.New("oopsy, tin invalid")
	ErrDuplicateTin    = errors.New("oopsy, tin already in use")
	ErrIdNotFound      = errors.New("oopsy, id not found")
	ErrEmptyGetContent = errors.New("oopsy, empty get content")
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
}

=======
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
>>>>>>> 8f0b9404 (init distributor refactor)
type GetResponse struct {
	Id          int
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
<<<<<<< HEAD
	Region      string
	Woreda      string
=======
	Woreda      string
	UserId      int
>>>>>>> 8f0b9404 (init distributor refactor)
}

type GetAllResponse struct {
	List []GetResponse
}

<<<<<<< HEAD
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

type CreateUserRequest struct {
	Distributor_Id int
	FirstName      string
	LastName       string
	Username       string
	Email          string
	Phone          string
}

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	Get(ctx context.Context, id int) (GetResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	Update(ctx context.Context, req *UpdateRequest) (int, error)
	GetAllUsers(ctx context.Context, id int) (GetAllUsers, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (int, error)
}

type DistributorService struct {
	DB          port.DB
	UserService user.Provider
}

func NewDistributorService(up user.Provider, db port.DB) Provider {
=======
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
>>>>>>> 8f0b9404 (init distributor refactor)
	return &DistributorService{
		UserService: up,
		DB:          db,
	}
}

func (d *DistributorService) CreateUser(ctx context.Context, req *CreateUserRequest) (int, error) {
	err := d.validateDistributor(ctx, req.Distributor_Id)
	if err != nil {
		return 0, err
	}

	// create user
	user_id, err := d.UserService.Create(ctx, &user.CreateRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
	})
	if err != nil {
		switch err {
		case user.ErrEmailNotValid,
			user.ErrPhoneNotValid,
			user.ErrPhoneOrEmailMandatory,
			user.ErrUsernameMandatory,
			user.ErrFirstNameMandatory:

			return 0, err
		default:
			return 0, ErrUnknown
		}
	}

	id, err := d.DB.CreateDistributorUser(ctx, &port.CreateUserAgentRequest{
		User_id:        user_id,
		Distributor_Id: req.Distributor_Id,
	})
	if err != nil {
		switch err {
		default:
			return id, ErrUnknown
		}
	}
	return id, nil
}

func (d *DistributorService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	err := d.validateTin(ctx, req.Tin)
	if err != nil {
		return 0, err
	}

	// create user
	user_id, err := d.UserService.Create(ctx, &user.CreateRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Username:  req.Username,
	})
	if err != nil {
		switch err {
		case user.ErrUnknown:
			return 0, ErrUnknown
		default:
			return 0, err
		}
	}

	// create distributor
	id, err := d.DB.Create(ctx, port.CreateRequest{
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
			d.UserService.Remove(ctx, user_id)
			return 0, ErrUnknown
		}
	}

	return id, nil
}
func (d *DistributorService) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := d.DB.Get(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
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
func (d *DistributorService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	resp := port.GetAllResponse{}

	if req.Name != "" {
		resp_name, err := d.DB.GetByName(ctx, req.Name)
		if err != nil {
			switch err {
			case ErrIdNotFound:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		if len(resp_name.List) != 0 {
			resp.List = append(resp.List, resp_name.List...)
		}
	}

	if req.Tin != "" {
		resp_tin, err := d.DB.GetByTin(ctx, req.Tin)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
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
func (d *DistributorService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp := port.GetAllResponse{}

	resp_name, err := d.DB.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
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
func (d *DistributorService) Update(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate id
	_, err := d.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return 0, err
		default:
			return 0, ErrUnknown
		}
	}

	if req.Name != "" {
		err = d.DB.UpdateName(ctx, &port.UpdateNameRequest{
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

		err = d.DB.UpdateTin(ctx, &port.UpdateTinRequest{
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
func (d *DistributorService) GetAllUsers(ctx context.Context, id int) (GetAllUsers, error) {
	var response_ids = []int{}
	users, err := d.DB.GetAllUserAgents(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetAllUsers{}, ErrEmptyGetContent
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
