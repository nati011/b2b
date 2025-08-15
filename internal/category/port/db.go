package category

import (
	"context"
)

type CreateRequest struct {
	Name string
}

type UpdateRequest struct {
	Id   int
	Name string
}

type GetResponse struct {
	Id   int
	Name string
}

type GetAllResponse struct {
	List []GetResponse
}

type Reader interface {
	GetAll(ctx context.Context) (GetAllResponse, error)
	Get(ctx context.Context, id int) (GetResponse, error)
}

type Writer interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	Remove(ctx context.Context, id int) error
	Update(ctx context.Context, req *UpdateRequest) error
}

type DB interface {
	Reader
	Writer
}
