package adapter

import (
	"context"
	"errors"
)

var (
	ErrSysUnknown = errors.New("unknown error")
	ErrNoRows     = errors.New("no rows")
)

type CreateRequest struct {
	Desc string
	Name string
}

type UpdateDescRequest struct {
	Id   int
	Desc string
}

type UpdateNameRequest struct {
	Id   int
	Name string
}

type GetResponse struct {
	Id   int
	Desc string
	Name string
}

type GetResourceResponse struct {
	Id     int
	Action string
	Name   string
}

type GetAllResponse struct {
	List []GetResponse
}

type GetAllResourcesResponse struct {
	List []GetResourceResponse
}

type Reader interface {
	GetByID(ctx context.Context, roleId int) (GetResponse, error)
	GetByName(ctx context.Context, name string) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (resourceId int, err error)
	UpdateDesc(ctx context.Context, req *UpdateDescRequest) (resourceId int, err error)
	UpdateName(ctx context.Context, req *UpdateNameRequest) (resourceId int, err error)
	Delete(ctx context.Context, roleId int) (err error)

	AddResource(ctx context.Context, roleId int, resourceId int) (err error)
	RemoveResource(ctx context.Context, roleId int, resourceId int) (err error)

	GetAllResources(ctx context.Context, roleId int) (resources GetAllResourcesResponse, err error)
}

type DB interface {
	Reader
	Writer
}
