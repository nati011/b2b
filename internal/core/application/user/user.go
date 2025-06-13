package user

import (
	"context"
	"errors"
	"log"
	"time"

	auth "b2b.nati011.github.com/internal/core/application/auth"
	role "b2b.nati011.github.com/internal/core/application/role"
	port "b2b.nati011.github.com/internal/port/application/user"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

var (
	ErrIdNotFound            = errors.New("¯\\_(•_•)_/¯, id not found")
	ErrEmailNotFound         = errors.New("¯\\_(•_•)_/¯, email not found")
	ErrPhoneNotFound         = errors.New("¯\\_(•_•)_/¯, phone not found")
	ErrUsernameNotFound      = errors.New("¯\\_(•_•)_/¯, username not found")
	ErrActiveStatusNotFound  = errors.New("¯\\_(•_•)_/¯, active_status not found")
	ErrFirstNameMandatory    = errors.New("¯\\_(•_•)_/¯, FirstName not supplied")
	ErrLastNameMandatory     = errors.New("¯\\_(•_•)_/¯, LastName not supplied")
	ErrUsernameMandatory     = errors.New("¯\\_(•_•)_/¯, Username not supplied")
	ErrPhoneOrEmailMandatory = errors.New("¯\\_(•_•)_/¯, phone or email mandatory")
	ErrEmptyGetContent       = errors.New("¯\\_(•_•)_/¯, content is empty")
	ErrUserAlreadyActive     = errors.New("¯\\_(•_•)_/¯, user already active")
	ErrUserAlreadyInactive   = errors.New("¯\\_(•_•)_/¯, user already inactive")
	ErrRoleDoesNotExist      = errors.New("¯\\_(•_•)_/¯, role does not exist")
	ErrRoleNotAssigned       = errors.New("¯\\_(•_•)_/¯, role not assigned")
	ErrRoleAlreadyAssigned   = errors.New("¯\\_(•_•)_/¯, role already assigned")
	ErrNoRoleAssigned        = errors.New("¯\\_(•_•)_/¯, no roles assigned")
	ErrResourceDoesNotExist  = errors.New("¯\\_(•_•)_/¯, resource does not exist")
	ErrUnknown               = errors.New("¯\\_(•_•)_/¯, unknown error")
	ErrPhoneNotValid         = errors.New("¯\\_(•_•)_/¯, phone number not valid")
	ErrEmailNotValid         = errors.New("¯\\_(•_•)_/¯, email not valid")
	ErrEmailTaken            = errors.New("¯\\_(•_•)_/¯, email already taken")
	ErrUserNameTaken         = errors.New("¯\\_(•_•)_/¯, username already taken")
	ErrPasswordMandatory     = errors.New("¯\\_(•_•)_/¯, password not supplied")
)

type CreateRequest struct {
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
	Password   string
}

type GetAssignedRoleResponse struct {
	Id int
}

type GetAllAssignedRoleResponse struct {
	List []GetAssignedRoleResponse
}

type UserProvider struct {
	UserId     int
	ProviderId string
}

type GetUserProviderResponse struct {
	List []UserProvider
}

type UpdateRequest struct {
	Id         int
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
}

type GetByParam struct {
	Email      string
	Phone      string
	Username   string
	IsActive   bool
	ExternalId string
}

type GetResponse struct {
	Id         int
	FirstName  string
	LastName   string
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

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (id int, err error)
	Get(ctx context.Context, id int) (resp GetResponse, err error)
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
	ResetPassword(ctx context.Context, email string) error
}

type UserService struct {
	db           port.DB
	role_service role.Provider
	auth_service auth.Provider
}

func NewUser(db port.DB, roleService role.Provider, authService auth.Provider) Provider {
	return &UserService{
		db:           db,
		role_service: roleService,
		auth_service: authService,
	}
}

func (u *UserService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	var providerResponse auth.RegisterUserResponse
	err := create_validateUserInfo(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Phone,
		req.Username,
		req.DOB,
	)
	if err != nil {
		return 0, err
	}

	if req.Password == "" {
		providerResponse, err = u.auth_service.CreateNewClientWithOutPassword(ctx, &auth.RegisterUserWithoutPasswordRequest{
			Email:       req.Email,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			PhoneNumber: req.Phone,
			Username:    req.Username,
		})
		if err != nil {
			log.Printf("failed to create auth client %v", err)
			switch err {
			case auth.ErrUsernameTaken:
				return 0, ErrUserNameTaken
			case auth.ErrEmailTaken:
				return 0, ErrEmailTaken
			case auth.ErrEmailNotSupplied:
				return 0, ErrEmailNotFound
			case auth.ErrUsernameNotSupplied:
				return 0, ErrUsernameNotFound
			case auth.ErrInvalidEmail:
				return 0, ErrEmailNotValid
			case auth.ErrFirstNameNotSupplied:
				return 0, ErrFirstNameMandatory
			case auth.ErrLastNameNotSupplied:
				return 0, ErrLastNameMandatory
			default:
				return 0, ErrUnknown
			}
		}
	} else {
		providerResponse, err = u.auth_service.CreateNewClientWithPassword(ctx, &auth.RegisterUserRequest{
			Email:       req.Email,
			Password:    req.Password,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			PhoneNumber: req.Phone,
			Username:    req.Username,
		})
		if err != nil {
			log.Printf("failed to create auth client %v", err)
			switch err {
			case auth.ErrUsernameTaken:
				return 0, ErrUserNameTaken
			case auth.ErrEmailTaken:
				return 0, ErrEmailTaken
			case auth.ErrEmailNotSupplied:
				return 0, ErrEmailNotFound
			case auth.ErrUsernameNotSupplied:
				return 0, ErrUsernameNotFound
			case auth.ErrInvalidEmail:
				return 0, ErrEmailNotValid
			case auth.ErrFirstNameNotSupplied:
				return 0, ErrFirstNameMandatory
			case auth.ErrLastNameNotSupplied:
				return 0, ErrLastNameMandatory
			default:
				return 0, ErrUnknown
			}
		}
	}

	user_id, err := u.db.CreateAndActivate(ctx, &port.CreateRequest{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Phone:      req.Phone,
		Username:   req.Username,
		DOB:        req.DOB,
		ExternalId: req.ExternalId,
	})
	if err != nil {
		log.Printf("Failed to create and activate user: %v err:%v", req.Email, err)
		u.auth_service.DeleteClient(ctx, providerResponse.Id)
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	log.Printf("Registered User %v", user_id)
	err = u.db.CreateUserProvider(ctx, &port.CreateUserProviderRequest{
		UserId:     user_id,
		ProviderId: providerResponse.Id,
	})

	if err != nil {
		switch err {
		default:
			u.Remove(ctx, user_id)
			return 0, ErrUnknown
		}
	}

	return user_id, nil
}

func (u *UserService) GetAll(ctx context.Context) (GetAllResponse, error) {
	res, err := u.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	prep_resp := GetAllResponse{}
	for _, i := range res.List {
		prep_resp.List = append(prep_resp.List, GetResponse{
			Id:         i.Id,
			FirstName:  i.FirstName,
			LastName:   i.LastName,
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

func (u *UserService) Get(ctx context.Context, id int) (GetResponse, error) {
	res, err := u.db.GetByID(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	return GetResponse(res), nil
}

func (u *UserService) GetUserProvider(ctx context.Context, id int) (GetUserProviderResponse, error) {
	user_providers, err := u.db.GetUserProvider(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetUserProviderResponse{}, ErrIdNotFound
		default:
			return GetUserProviderResponse{}, ErrUnknown
		}
	}

	resp := GetUserProviderResponse{}
	for _, i := range user_providers.List {
		resp.List = append(resp.List, UserProvider{
			UserId:     i.UserId,
			ProviderId: i.ProviderId,
		})
	}
	return resp, nil
}

func (u *UserService) getUserProviderByEmail(ctx context.Context, email string) (UserProvider, error) {
	user_provider, err := u.db.GetUserProviderByEmail(ctx, email)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return UserProvider{}, ErrIdNotFound
		default:
			return UserProvider{}, ErrUnknown
		}
	}

	return (UserProvider)(user_provider), nil
}

func (u *UserService) ResetPassword(ctx context.Context, email string) error {
	user_provider, err := u.getUserProviderByEmail(ctx, email)
	if err != nil {
		return err
	}
	req := auth.InitClientCredentialsResetRequest{
		UserId: user_provider.ProviderId,
		Email:  email,
	}

	err = u.auth_service.InitClientCredentialsReset(ctx, req)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetByParam(ctx context.Context, req *GetByParam) (GetAllResponse, error) {
	prep_resp := GetAllResponse{}

	if req.Email != "" {
		res_email, err := u.db.GetByEmail(ctx, req.Email)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
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
					prep_resp.List = append(prep_resp.List, (GetResponse)(i))
				}
			}
		}
	}

	if req.Phone != "" {
		res_phone, err := u.db.GetByPhone(ctx, req.Phone)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
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
					prep_resp.List = append(prep_resp.List, (GetResponse)(i))
				}
			}
		}
	}

	if req.Username != "" {
		res_username, err := u.db.GetByUsername(ctx, req.Username)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
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
					prep_resp.List = append(prep_resp.List, (GetResponse)(i))
				}
			}
		}
	}

	res_isActive, err := u.db.GetByActiveStatus(ctx, req.IsActive)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
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
				prep_resp.List = append(prep_resp.List, (GetResponse)(i))
			}
		}
	}

	if req.ExternalId != "" {
		res_externalId, err := u.db.GetByEmail(ctx, req.ExternalId)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
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
					prep_resp.List = append(prep_resp.List, (GetResponse)(i))
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
	_, err := u.Get(ctx, id)
	if err != nil {
		return err
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
	_, err := u.Get(ctx, id)
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

func (u *UserService) AssignRole(ctx context.Context, id int, roleId int) error {
	//validate id
	user, err := u.Get(ctx, id)
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
		Id: roleId,
	})
	if err != nil {
		switch err {
		case role.ErrIdNotFound:
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
		if i.Id == roleId {
			alreadyAssigned = true
		}
	}

	if alreadyAssigned {
		return ErrRoleAlreadyAssigned
	}

	//assign
	err = u.db.AssignRole(ctx, id, roleId)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}

	//assign for auth
	err = u.auth_service.AssignRole(ctx, user.Id, roleId)
	if err != nil {
		log.Printf("Failed to assign user auth role: %v", err)
		err := u.db.RemoveAssignedRole(ctx, user.Id, roleId)
		if err != nil {
			log.Printf("Failed to rollback assigned role")
		}
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (u *UserService) GetAllAssignedRoles(ctx context.Context, id int) (GetAllAssignedRoleResponse, error) {
	//validate id
	_, err := u.Get(ctx, id)
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
		case port_commons.ErrSysNoRows:
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

func (u *UserService) RemoveAssignedRole(ctx context.Context, id int, roleId int) error {
	//validate id
	_, err := u.Get(ctx, id)
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
		Id: roleId,
	})
	if err != nil {
		switch err {
		case role.ErrIdNotFound:
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
		if i.Id == roleId {
			alreadyAssigned = true
		}
	}

	if !alreadyAssigned {
		return ErrRoleNotAssigned
	}

	//remove
	err = u.db.RemoveAssignedRole(ctx, id, roleId)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}

	//remove from authService
	err = u.auth_service.RemoveRole(ctx, id, roleId)
	if err != nil {
		log.Printf("Failed to remove assign roleId: %v err:%v", roleId, err)
		u.AssignRole(ctx, id, roleId)
		return ErrUnknown
	}
	return nil
}

func (u *UserService) HasRole(ctx context.Context, id int, role_id int) (bool, error) {
	//validate id
	_, err := u.Get(ctx, id)
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
		case role.ErrIdNotFound:
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
	_, err := u.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}

	//validate input
	err = update_validateUserInfo(
		ctx,
		req.FirstName,
		req.Email,
		req.Phone,
		req.Username,
		req.DOB,
	)
	if err != nil {
		return GetResponse{}, err
	}

	if req.FirstName != "" {
		_, err := u.db.UpdateFirstName(ctx,
			&port.UpdateFirstNameRequest{
				Id:        req.Id,
				FirstName: req.FirstName,
			})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}

	if req.Email != "" {
		_, err := u.db.UpdateEmail(ctx, &port.UpdateEmailRequest{
			Id:    req.Id,
			Email: req.Email,
		})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}

	if req.Phone != "" {
		_, err := u.db.UpdatePhone(ctx, &port.UpdatePhoneRequest{
			Id:    req.Id,
			Phone: req.Phone,
		})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}

	if req.Username != "" {
		_, err := u.db.UpdateUsername(ctx, &port.UpdateUsernameRequest{
			Id:       req.Id,
			Username: req.Username,
		})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}
	if !req.DOB.IsZero() {
		_, err = u.db.UpdateDOB(ctx, &port.UpdateDOBRequest{
			Id:  req.Id,
			DOB: req.DOB,
		})
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
	}

	resp, err := u.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	return (GetResponse)(resp), nil
}

func (u *UserService) Remove(ctx context.Context, id int) error {
	//validate id
	_, err := u.Get(ctx, id)
	if err != nil {
		return err
	}

	resp, err := u.GetUserProvider(ctx, id)
	if err != nil {
		switch err {
		case ErrEmptyGetContent:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	for _, i := range resp.List {
		err = u.auth_service.DeleteClient(ctx, i.ProviderId)
		if err != nil {
			switch err {
			default:
				return ErrUnknown
			}
		}
	}
	//remove
	log.Printf("Deleting Id %v", id)
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
	user, err := u.Get(ctx, id)
	if err != nil {
		switch err {
		case ErrUnknown:
			return false, ErrUnknown
		default:
			return false, ErrIdNotFound
		}
	}
	return user.IsActive, nil
}
