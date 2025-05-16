package email_template

import (
	"context"
)

type CreateRequest struct {
	Name         string
	HtmlTemplate string
}

type GetResponse struct {
	Id           int
	Name         string
	HtmlTemplate string
}

type GetAllResponse struct {
	List []GetResponse
}

type UpdateRequest struct {
	Id           int
	Name         string
	HtmlTemplate string
}

type Reader interface {
	Get(context.Context, int) (GetResponse, error)
	GetByName(context.Context, string) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
	Update(context.Context, *UpdateRequest) error
}

type DB interface {
	Reader
	Writer
}
