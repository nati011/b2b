package distributor

import (
	"context"
	"errors"
	"log"

	"b2b.nati011.github.com/internal/core/application/user"
	distributorApproval "b2b.nati011.github.com/internal/core/domain/distributor_approval"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/distributor"
)

var (
	ErrUnknown                    = errors.New(" unknown error")
	ErrInvalidTin                 = errors.New(" tin invalid")
	ErrDuplicateTin               = errors.New(" tin already in use")
	ErrIdNotFound                 = errors.New(" id not found")
	ErrEmptyGetContent            = errors.New(" empty get content")
	ErrDistributorAlreadyInactive = errors.New(" distributor already inactive")
	ErrDistributorAlreadyActive   = errors.New(" distributor already active")
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

type GetResponse struct {
	Id          int
	Name        string
	Tin         string
	Latitude    string
	Longitude   string
	GeneralZone string
	Region      string
	Woreda      string
	IsActive    bool
	Verdict     string
}

type GetAllResponse struct {
	List       []GetResponse
	TotalCount int64
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

type CreateUserRequest struct {
	DistributorId int
	FirstName     string
	LastName      string
	Username      string
	Email         string
	Phone         string
}

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	Get(ctx context.Context, id int) (GetResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	Update(ctx context.Context, req *UpdateRequest) (int, error)
	GetAllUsers(ctx context.Context, id int) (GetAllUsers, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (int, error)
	Activate(ctx context.Context, id int) error
	Dectivate(ctx context.Context, id int) error
	ApproveOnboardingRequest(ctx context.Context, distributorId int) error
	RejectOnboardingRequest(ctx context.Context, distributorId int, comment string) error
}

type DistributorService struct {
	DB                         port.DB
	UserService                user.Provider
	DistributorApprovalService distributorApproval.Provider
}

func NewDistributorService(up user.Provider, db port.DB, dap distributorApproval.Provider) Provider {
	return &DistributorService{
		UserService:                up,
		DB:                         db,
		DistributorApprovalService: dap,
	}
}

func (d *DistributorService) CreateUser(ctx context.Context, req *CreateUserRequest) (int, error) {
	err := d.validateDistributor(ctx, req.DistributorId)
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
		Distributor_Id: req.DistributorId,
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
		IsActive:    resp.IsActive,
		Verdict:     resp.Verdict,
	}, nil
}
func (d *DistributorService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	resp := port.GetAllResponse{}

	if req.Name != "" {
		resp_name, err := d.DB.GetByName(ctx, req.Name)
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
		resp_tin, err := d.DB.GetByTin(ctx, req.Tin)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
				return GetAllResponse{}, ErrEmptyGetContent
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
			IsActive:    i.IsActive,
			Verdict:     i.Verdict,
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
			Verdict:     i.Verdict,
			IsActive:    i.IsActive,
		})
	}
	service_resp.TotalCount = resp_name.TotalCount
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
		case port_commons.ErrSysNoRows:
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

func (d *DistributorService) Activate(ctx context.Context, id int) error {
	resp, err := d.Get(ctx, id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return err
		default:
			return ErrUnknown
		}
	}

	if resp.IsActive {
		return ErrDistributorAlreadyActive
	}

	err = d.DB.Activate(ctx, id)
	if err != nil {
		log.Printf("failed to activate distributor id:%v", id)
		return ErrUnknown
	}
	return nil
}

func (d *DistributorService) Dectivate(ctx context.Context, id int) error {
	resp, err := d.Get(ctx, id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return err
		default:
			return ErrUnknown
		}
	}

	if !resp.IsActive {
		return ErrDistributorAlreadyInactive
	}

	err = d.DB.Dectivate(ctx, id)
	if err != nil {
		log.Printf("failed to deactivate distributor id:%v", id)
		return ErrUnknown
	}
	return nil
}

func (d *DistributorService) ApproveOnboardingRequest(ctx context.Context, distributorId int) error {
	err := d.DistributorApprovalService.Approve(ctx, &distributorApproval.ApprovalRequest{
		DistributorId: distributorId,
	})
	if err != nil {
		switch err {
		case ErrUnknown:
			log.Print("failed to approve distributor")
			return ErrUnknown
		default:
			return err
		}
	}
	return nil
}

func (d *DistributorService) RejectOnboardingRequest(ctx context.Context, distributorId int, comment string) error {
	err := d.DistributorApprovalService.Reject(ctx, &distributorApproval.RejectionRequest{
		DistributorId: distributorId,
		Comment:       comment,
	})
	if err != nil {
		switch err {
		case ErrUnknown:
			log.Print("failed to reject distributor")
			return ErrUnknown
		default:
			return err
		}
	}
	return nil
}
