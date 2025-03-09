package user

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSysNoRows  = errors.New("no rows")
	ErrSysUnknown = errors.New("unknown error")
)

type CreateRequest struct {
	FirstName  string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	IsActive   bool
	ExternalId string
}

type GetResponse struct {
	Id         int
	FirstName  string
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

type UpdateEmailRequest struct {
	Id    int
	Email string
}

type UpdatePhoneRequest struct {
	Id    int
	Phone string
}

type UpdateUsernameRequest struct {
	Id       int
	Username string
}

type UpdateDOBRequest struct {
	Id  int
	DOB time.Time
}

type UpdateFirstNameRequest struct {
	Id        int
	FirstName string
}

type UpdateIsActiveRequest struct {
	Id       int
	IsActive bool
}

type GetAssignedRoleResponse struct {
	Id int
}

type GetAllAssignedRoleResponse struct {
	List []GetAssignedRoleResponse
}

type Reader interface {
	GetByID(ctx context.Context, id int) (GetResponse, error)
	GetByEmail(ctx context.Context, email string) (GetAllResponse, error)
	GetByPhone(ctx context.Context, phone string) (GetAllResponse, error)
	GetByUsername(ctx context.Context, username string) (GetAllResponse, error)
	GetByActiveStatus(ctx context.Context, status bool) (GetAllResponse, error)
	GetByExternalId(ctx context.Context, extId string) (GetAllResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
	CreateAndActivate(context.Context, *CreateRequest) (int, error)
	UpdateFirstName(context.Context, *UpdateFirstNameRequest) (int, error)
	UpdateEmail(context.Context, *UpdateEmailRequest) (int, error)
	UpdatePhone(context.Context, *UpdatePhoneRequest) (int, error)
	UpdateUsername(context.Context, *UpdateUsernameRequest) (int, error)
	UpdateDOB(context.Context, *UpdateDOBRequest) (int, error)
	UpdateIsActiveStatus(context.Context, *UpdateIsActiveRequest) (int, error)

	Delete(context.Context, int) error
	AssignRole(ctx context.Context, id int, roleId int) error
	RemoveAssignedRole(ctx context.Context, id int, roleId int) error
	GetAllAssignedRole(ctx context.Context, id int) (GetAllAssignedRoleResponse, error)
}

type DB interface {
	Reader
	Writer
}
