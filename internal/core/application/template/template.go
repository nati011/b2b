package template

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/application/email-template"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

var (
	ErrInvalidTemplate = errors.New("¯\\_(o_o)_/¯, invalid template")
	ErrInvalidName     = errors.New("¯\\_(o_o)_/¯, invalid name")
	ErrUnknown         = errors.New("¯\\_(o_o)_/¯, unknown error")
	ErrDuplicateName   = errors.New("¯\\_(o_o)_/¯, duplicate name")
	ErrIdNotFound      = errors.New("¯\\_(o_o)_/¯, id not found")
	ErrNameNotFound    = errors.New("¯\\_(o_o)_/¯, name not found")
	ErrEmptyGetContent = errors.New("¯\\_(o_o)_/¯, empty get content")
)

type CreateRequest struct {
	Name         string
	HtmlTemplate string
}

type CreateResponse struct {
	Name    string
	Message string
}

type GetResponse struct {
	Id           int
	Name         string
	HtmlTemplate string
}

type GetAllResponse struct {
	List []GetResponse
}

type Provider interface {
	Create(context.Context, *CreateRequest) (int, error)
	GetAll(context.Context) (GetAllResponse, error)
	Get(context.Context, int) (GetResponse, error)
	GetByName(ctx context.Context, name string) (GetResponse, error)
}

type Template struct {
	db port.DB
}

func NewTemplateService(db_provider port.DB) Provider {
	return &Template{
		db: db_provider,
	}
}

func (t *Template) Create(ctx context.Context, req *CreateRequest) (int, error) {
	if err := t.ValidateName(ctx, req.Name); err != nil {
		return 0, err
	}
	if err := ValidateHTML(req.HtmlTemplate); err != nil {
		return 0, err
	}
	id, err := t.db.Create(ctx, &port.CreateRequest{
		Name:         req.Name,
		HtmlTemplate: req.HtmlTemplate,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}
	return id, nil
}

func (t *Template) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := t.db.Get(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	return GetResponse(resp), nil
}

func (t *Template) GetByName(ctx context.Context, name string) (GetResponse, error) {
	resp, err := t.db.GetByName(ctx, name)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetResponse{}, ErrNameNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	return GetResponse(resp), nil
}

func (t *Template) GetAll(ctx context.Context) (GetAllResponse, error) {
	rslt, err := t.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	response := GetAllResponse{}
	for _, i := range rslt.List {
		response.List = append(response.List, GetResponse(i))
	}

	return response, nil
}
