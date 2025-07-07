package adapter

import (
	"context"
)

type CreateRequest struct {
	Action   string
	Name     string
	Resource string
}

type UpdateActionRequest struct {
	Id     int
	Action string
}

type UpdateNameRequest struct {
	Id   int
	Name string
}

type UpdateResouceRequest struct {
	Id       int
	Resource string
}

type GetResponse struct {
	Id       int
	Action   string
	Resource string
	Name     string
}

type GetAllResponse struct {
	List []GetResponse
}

type Reader interface {
	GetByID(context.Context, int) (GetResponse, error)
	GetByName(context.Context, string) (GetResponse, error)
	GetByResource(context.Context, string) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
	UpdateAction(context.Context, *UpdateActionRequest) (int, error)
	UpdateName(context.Context, *UpdateNameRequest) (int, error)
	UpdateResource(context.Context, *UpdateResouceRequest) (int, error)
	Delete(context.Context, int) error
}

type DB interface {
	Reader
	Writer
}
