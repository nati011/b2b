package resource

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/resource"
)

var (
	ErrDuplicateName = errors.New("oopsy, duplicate name")
	ErrEmptyAction   = errors.New("oopsy, empty action")
	ErrEmptyName     = errors.New("oopsy, empty name")
	ErrIdNotFound    = errors.New("oopsy, id not found")
)

type CreateRequest struct {
	Action string
	Name   string
}

type UpdateRequest struct {
	Id     int
	Action string
	Name   string
}

type GetRequest struct {
	Id   int
	Name string
}

type GetResponse struct {
	Id     int
	Action string
	Name   string
}

type GetAllResponse struct {
	List []GetResponse
}

type Provider interface {
	Get(context.Context, *GetRequest) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)

	Create(context.Context, *CreateRequest) (int, error)
	Update(context.Context, *UpdateRequest) (int, error)
	Delete(context.Context, int) error
}

type ResourceProvider struct {
	db port.DB
}

func NewResource() Provider {
	return &ResourceProvider{}
}

func (r *ResourceProvider) Create(ctx context.Context, req *CreateRequest) (int, error) {
	return 0, nil
}

func (r *ResourceProvider) Update(ctx context.Context, req *UpdateRequest) (int, error) {
	return 0, nil
}

func (r *ResourceProvider) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *ResourceProvider) Get(ctx context.Context, req *GetRequest) (GetResponse, error) {
	return GetResponse{}, nil
}

func (r *ResourceProvider) GetAll(ctx context.Context) (GetAllResponse, error) {
	return GetAllResponse{
		List: []GetResponse{},
	}, nil
}
