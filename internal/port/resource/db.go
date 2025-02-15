package adapter

import "context"

type CreateRequest struct {
	Action string
	Name   string
}

type UpdateActionRequest struct {
	Id     string
	Action string
}

type UpdateNameRequest struct {
	Id   string
	Name string
}

type GetResponse struct {
	Id     string
	Action string
	Name   string
}

type GetAllResponse struct {
	List []GetResponse
}

type Reader interface {
	GetByID(*context.Context, int) (GetResponse, error)
	GetByName(*context.Context, string) (GetResponse, error)
	GetAll(*context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(*context.Context, *CreateRequest) (int, error)
	UpdateAction(*context.Context, *UpdateActionRequest) (int, error)
	UpdateName(*context.Context, *UpdateNameRequest) (int, error)
	Delete(*context.Context, int) error
}

type DB interface {
	Reader
	Writer
}
