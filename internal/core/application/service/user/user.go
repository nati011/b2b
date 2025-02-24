package user

import (
	"context"
	"errors"
	"time"

	role "b2b.nati011.github.com/internal/core/application/service/role"
	port "b2b.nati011.github.com/internal/port/user"
)

var (
	ErrIdNotFound            = errors.New("oopsy, id not found")
	ErrEmailNotFound         = errors.New("oopsy, email not found")
	ErrPhoneNotFound         = errors.New("oopsy, phone not found")
	ErrUsernameNotFound      = errors.New("oopsy, username not found")
	ErrActiveStatusNotFound  = errors.New("oopsy, active_status not found")
	ErrFullNameMandatory     = errors.New("oopsy, fullname not supplied")
	ErrPhoneOrEmailMandatory = errors.New("oopsy, phone or email mandatory")
	ErrEmptyGetContent       = errors.New("oopsy, content is empty")
	ErrUserAlreadyActive     = errors.New("oopsy, user already active")
	ErrUserAlreadyInactive   = errors.New("oopsy, user already inactive")
	ErrRoleDoesNotExist      = errors.New("oppsy, role does not exist")
	ErrRoleNotAssigned       = errors.New("oopsy, role not assigned")
	ErrRoleAlreadyAssigned   = errors.New("oopsy, role already assigned")
	ErrNoRoleAssigned        = errors.New("oopsy, no roles assigned")
	ErrResourceDoesNotExist  = errors.New("oopsy, resource does not exist")
	ErrUnknown               = errors.New("oopsy, unknown error")
	ErrPhoneNotValid         = errors.New("oopsy, phone number not valid")
	ErrEmailNotValid         = errors.New("oopsy, email not valid")
)

type CreateRequest struct {
	FullName   string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
}

type GetResponse struct {
	Id         int
	FullName   string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	IsActive   bool
	ExternalId string
}
type GetAllResponse struct {
	List []GetResponse
}

type GetByParam struct {
	ID         int
	Email      string
	Phone      string
	Username   string
	IsActive   bool
	ExternalId string
}

type GetAssignedRoleResponse struct {
	Id int
}

type GetAllAssignedRoleResponse struct {
	List []GetAssignedRoleResponse
}

type UpdateRequest struct {
	Id         int
	FullName   string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
}

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (id int, err error)
	GetAll(ctx context.Context) (resp GetAllResponse, err error)
	GetByParam(ctx context.Context, req *GetByParam) (resp GetAllResponse, err error)
	Activate(ctx context.Context, id int) (err error)
	Deactivate(ctx context.Context, id int) (err error)
	IsActive(ctx context.Context, id int) (resp bool, err error)
	AssignRole(ctx context.Context, id int, role_id int) (err error)
	RemoveAssignedRole(ctx context.Context, id int, role_id int) (err error)
	GetAllAssignedRoles(ctx context.Context, id int) (resp GetAllAssignedRoleResponse, err error)
	HasRole(ctx context.Context, id int, role_id int) (resp bool, err error)
	Update(ctx context.Context, req *UpdateRequest) (resp GetResponse, err error)
	Remove(ctx context.Context, id int) (err error)
}

type UserService struct {
	db           port.DB
	role_service role.Provider
}

func NewUser(db port.DB, roleService role.Provider) Provider {
	return &UserService{
		db:           db,
		role_service: roleService,
	}
}

func (u *UserService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	//validate input
	err := create_validateUserInfo(
		ctx,
		req.FullName,
		req.Email,
		req.Phone,
		req.Username,
		req.DOB,
	)
	if err != nil {
		return 0, err
	}

	//create user
	user_id, err := u.db.CreateAndActivate(ctx, &port.CreateRequest{
		FullName:   req.FullName,
		Email:      req.Email,
		Phone:      req.Phone,
		Username:   req.Username,
		DOB:        req.DOB,
		ExternalId: req.ExternalId,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	return user_id, nil
}

func (u *UserService) GetAll(ctx context.Context) (GetAllResponse, error) {
	res, err := u.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	prep_resp := GetAllResponse{}
	for _, i := range res.List {
		prep_resp.List = append(prep_resp.List, GetResponse{
			Id:         i.Id,
			FullName:   i.FullName,
			Email:      i.Email,
			Phone:      i.Phone,
			Username:   i.Username,
			DOB:        i.DOB,
			IsActive:   i.IsActive,
			ExternalId: i.ExternalId,
		})
	}

	if len(prep_resp.List) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}
	return prep_resp, nil
}

func (u *UserService) GetByParam(ctx context.Context, req *GetByParam) (GetAllResponse, error) {
	prep_resp := GetAllResponse{}
	if req.ID != 0 {
		res_id, err := u.db.GetByID(ctx, req.ID)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		if res_id.Id != 0 {
			already_present := false
			for _, i := range prep_resp.List {
				if i.Id == res_id.Id {
					already_present = true
					continue
				}
			}
			if !already_present {
				prep_resp.List = append(prep_resp.List, GetResponse(res_id))
			}
		}
	}

	if req.Email != "" {
		res_email, err := u.db.GetByEmail(ctx, req.Email)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		if len(res_email.List) != 0 {
			already_present := false
			for _, i := range prep_resp.List {
				for _, j := range res_email.List {
					if i.Id == j.Id {
						already_present = true
						continue
					}
				}
			}
			if !already_present {
				for _, i := range res_email.List {
					prep_resp.List = append(prep_resp.List, GetResponse(i))
				}
			}
		}
	}

	if req.Phone != "" {
		res_phone, err := u.db.GetByPhone(ctx, req.Phone)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		if len(res_phone.List) != 0 {
			already_present := false
			for _, i := range prep_resp.List {
				for _, j := range res_phone.List {
					if i.Id == j.Id {
						already_present = true
						continue
					}
				}
			}
			if !already_present {
				for _, i := range res_phone.List {
					prep_resp.List = append(prep_resp.List, GetResponse(i))
				}
			}
		}
	}

	if req.Username != "" {
		res_username, err := u.db.GetByUsername(ctx, req.Username)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}
		if len(res_username.List) != 0 {
			already_present := false
			for _, i := range prep_resp.List {
				for _, j := range res_username.List {
					if i.Id == j.Id {
						already_present = true
						continue
					}
				}
			}
			if !already_present {
				for _, i := range res_username.List {
					prep_resp.List = append(prep_resp.List, GetResponse(i))
				}
			}
		}
	}

	res_isActive, err := u.db.GetByActiveStatus(ctx, req.IsActive)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	if len(res_isActive.List) != 0 {
		already_present := false
		for _, i := range prep_resp.List {
			for _, j := range res_isActive.List {
				if i.Id == j.Id {
					already_present = true
					continue
				}
			}
		}
		if !already_present {
			for _, i := range res_isActive.List {
				prep_resp.List = append(prep_resp.List, GetResponse(i))
			}
		}
	}

	if req.ExternalId != "" {
		res_externalId, err := u.db.GetByEmail(ctx, req.ExternalId)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		if len(res_externalId.List) != 0 {
			already_present := false
			for _, i := range prep_resp.List {
				for _, j := range res_externalId.List {
					if i.Id == j.Id {
						already_present = true
						continue
					}
				}
			}
			if !already_present {
				for _, i := range res_externalId.List {
					prep_resp.List = append(prep_resp.List, GetResponse(i))
				}
			}
		}
	}

	if len(prep_resp.List) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}
	return prep_resp, nil
}

func (u *UserService) Activate(ctx context.Context, id int) error {
	//validate id
	_, err := u.GetByParam(ctx, &GetByParam{
		ID: id,
	})
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//validate active status
	isActive, err := u.IsActive(ctx, id)
	if err != nil {
		return ErrUnknown
	}
	if isActive {
		return ErrUserAlreadyActive
	}

	//activate
	_, err = u.db.UpdateIsActiveStatus(ctx, &port.UpdateIsActiveRequest{
		Id:       id,
		IsActive: true,
	})
	if err != nil {
		return ErrUnknown
	}
	return nil
}

func (u *UserService) Deactivate(ctx context.Context, id int) error {
	//validate id
	_, err := u.GetByParam(ctx, &GetByParam{
		ID: id,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//validate active status
	isActive, err := u.IsActive(ctx, id)
	if err != nil {
		return ErrUnknown
	}
	if !isActive {
		return ErrUserAlreadyInactive
	}

	//deactivate
	_, err = u.db.UpdateIsActiveStatus(ctx, &port.UpdateIsActiveRequest{
		Id:       id,
		IsActive: false,
	})
	if err != nil {
		return ErrUnknown
	}
	return nil
}

func (u *UserService) AssignRole(ctx context.Context, id int, role_id int) error {
	//validate id
	_, err := u.GetByParam(ctx, &GetByParam{
		ID: id,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//validate if role exists
	_, err = u.role_service.Get(ctx, &role.GetRequest{
		Id: role_id,
	})
	if err != nil {
		switch err {
		case role.ErrEmptyGetContent:
			return ErrRoleDoesNotExist
		default:
			return ErrUnknown
		}
	}

	//validate if role was assigned
	assigned_roles, err := u.GetAllAssignedRoles(ctx, id)
	if err != nil {
		switch err {
		case ErrNoRoleAssigned:
		default:
			return ErrUnknown
		}
	}
	alreadyAssigned := false
	for _, i := range assigned_roles.List {
		if i.Id == role_id {
			alreadyAssigned = true
		}
	}

	if alreadyAssigned {
		return ErrRoleAlreadyAssigned
	}

	//assign
	err = u.db.AssignRole(ctx, id, role_id)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (u *UserService) GetAllAssignedRoles(ctx context.Context, id int) (GetAllAssignedRoleResponse, error) {
	//validate id
	_, err := u.GetByParam(ctx, &GetByParam{
		ID: id,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return GetAllAssignedRoleResponse{}, ErrIdNotFound
		default:
			return GetAllAssignedRoleResponse{}, ErrUnknown
		}
	}

	//get all assigned roles
	assigend_roles, err := u.db.GetAllAssignedRole(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetAllAssignedRoleResponse{}, ErrNoRoleAssigned
		default:
			return GetAllAssignedRoleResponse{}, ErrUnknown
		}
	}
	resp := GetAllAssignedRoleResponse{}
	for _, i := range assigend_roles.List {
		resp.List = append(resp.List, GetAssignedRoleResponse{
			Id: i.Id,
		})
	}

	if len(resp.List) == 0 {
		return resp, ErrNoRoleAssigned
	}
	return resp, nil
}

func (u *UserService) RemoveAssignedRole(ctx context.Context, id int, role_id int) error {
	//validate id
	_, err := u.GetByParam(ctx, &GetByParam{
		ID: id,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//validate if role exists
	_, err = u.role_service.Get(ctx, &role.GetRequest{
		Id: role_id,
	})
	if err != nil {
		switch err {
		case role.ErrEmptyGetContent:
			return ErrRoleDoesNotExist
		default:
			return ErrUnknown
		}
	}

	//validate if role was assigned
	assigned_roles, err := u.GetAllAssignedRoles(ctx, id)
	if err != nil {
		switch err {
		case ErrNoRoleAssigned:
			return ErrRoleNotAssigned
		default:
			return ErrUnknown
		}
	}
	alreadyAssigned := false
	for _, i := range assigned_roles.List {
		if i.Id == role_id {
			alreadyAssigned = true
		}
	}

	if !alreadyAssigned {
		return ErrRoleNotAssigned
	}

	//remove
	err = u.db.RemoveAssignedRole(ctx, id, role_id)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (u *UserService) HasRole(ctx context.Context, id int, role_id int) (bool, error) {
	//validate id
	_, err := u.GetByParam(ctx, &GetByParam{
		ID: id,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return false, ErrIdNotFound
		default:
			return false, ErrUnknown
		}
	}

	//validate role id
	_, err = u.role_service.Get(ctx, &role.GetRequest{
		Id: id,
	})
	if err != nil {
		switch err {
		case role.ErrEmptyGetContent:
			return false, ErrRoleDoesNotExist
		default:
			return false, nil
		}
	}

	//get all roles
	allAssignedRoles, err := u.GetAllAssignedRoles(ctx, id)
	if err != nil {
		switch err {
		case ErrNoRoleAssigned:
			return false, nil
		default:
			return false, ErrUnknown
		}
	}
	for _, i := range allAssignedRoles.List {
		if i.Id == id {
			return true, nil
		}
	}
	return false, nil
}

func (u *UserService) Update(ctx context.Context, req *UpdateRequest) (GetResponse, error) {
	//validate id
	_, err := u.GetByParam(ctx, &GetByParam{
		ID: req.Id,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}

	//validate input
	err = update_validateUserInfo(
		ctx,
		req.FullName,
		req.Email,
		req.Phone,
		req.Username,
		req.DOB,
	)
	if err != nil {
		return GetResponse{}, err
	}

	if req.FullName != "" {
		_, err := u.db.UpdateFullName(ctx,
			&port.UpdateFullNameRequest{
				Id:       req.Id,
				FullName: req.FullName,
			})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}

	if req.Email != "" {
		_, err := u.db.UpdateEmail(ctx, &port.UpdateEmailRequest{})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}

	if req.Phone != "" {
		_, err := u.db.UpdatePhone(ctx, &port.UpdatePhoneRequest{})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}

	if req.Username != "" {
		_, err := u.db.UpdateUsername(ctx, &port.UpdateUsernameRequest{})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}
	if !req.DOB.IsZero() {
		_, err = u.db.UpdateDOB(ctx, &port.UpdateDOBRequest{})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}

	resp, err := u.GetByParam(ctx, &GetByParam{
		ID: req.Id,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	return resp.List[0], nil
}

func (u *UserService) Remove(ctx context.Context, id int) error {
	//validate id
	_, err := u.GetByParam(ctx, &GetByParam{
		ID: id,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//remove
	err = u.db.Delete(ctx, id)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (u *UserService) IsActive(ctx context.Context, id int) (bool, error) {
	// validate id
	user, err := u.GetByParam(ctx, &GetByParam{
		ID: id,
	})
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return false, ErrIdNotFound
		default:
			return false, ErrUnknown
		}
	}

	//get user id
	return user.List[0].IsActive, nil
}
