package email

import (
	"context"
	"errors"
)

var (
	ErrSysDuplicateName_L1 = errors.New("duplicate name")
	ErrSysUnknown_L1       = errors.New("unknown")
)

type CreateRequest struct {
	Name         string
	HtmlTemplate string
}

type CreateResponse struct {
	Name string
}

type GetResponse struct {
	Name         string
	HtmlTemplate string
}

type GetAllResponse struct {
	List []base
}
type base struct {
	Name         string
	HtmlTemplate string
}

type UpdateRequest struct {
	Name         string
	HtmlTemplate string
}

type UpdateResponse struct {
	Name string
}
type template struct {
	name string
	html string
}

type Reader interface {
	Get(context.Context, string) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (CreateResponse, error)
	Update(context.Context, *UpdateRequest) (UpdateResponse, error)
}

type ReaderWriter interface {
	Reader
	Writer
}
