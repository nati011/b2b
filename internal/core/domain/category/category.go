package category

import (
	"context"
	"errors"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/category"
)

var (
	ErrNameIsNotSupplied = errors.New("oopsy, name is mandatory")
	ErrDescIsNotSupplied = errors.New("oopsy, desc is mandatory")
	ErrIdNotFound        = errors.New("oopsy, id not found")
	ErrEmptyGetContent   = errors.New("oopsy, empty content")
	ErrDuplicateName     = errors.New("oopsy, name already taken")
	ErrUnknown           = errors.New("oopsy, unknown error")
)

type CreateRequest struct {
	Name string
}

type GetResponse struct {
	Id   int
	Name string
}

type GetAllResponse struct {
	List []GetResponse
}

type Provider interface {
	Create(ctx context.Context, req *CreateRequest) (int, error)
	Remove(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
}

type CategoryService struct {
	db port.DB
}

func NewCategory(db port.DB) Provider {
	return &CategoryService{
		db: db,
	}
}

func (c *CategoryService) Create(ctx context.Context, req *CreateRequest) (int, error) {
	//validate
	err := c.validateName(ctx, req.Name)
	if err != nil {
		return 0, err
	}

	id, err := c.db.Create(ctx, &port.CreateRequest{
		Name: req.Name,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}
	return id, nil
}

func (c *CategoryService) Remove(ctx context.Context, id int) error {
	//validate
	resp, err := c.Get(ctx, id)
	if err != nil {
		switch err {
		case ErrIdNotFound:
			return err
		default:
			return ErrUnknown
		}
	}

	err = c.db.Remove(ctx, resp.Id)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (c *CategoryService) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := c.db.Get(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrNoRows:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	return GetResponse(resp), nil
}

func (c *CategoryService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp, err := c.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	if len(resp.List) == 0 {
		return GetAllResponse{}, ErrEmptyGetContent
	}
	getResp := GetAllResponse{}
	for _, i := range resp.List {
		getResp.List = append(getResp.List, GetResponse(i))
	}
	return getResp, nil
}
