package resource

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/application/resource"
)

var (
	ErrDuplicateName      = errors.New("oopsy, duplicate name")
	ErrEmptyAction        = errors.New("oopsy, empty action")
	ErrEmptyName          = errors.New("oopsy, empty name")
	ErrIdNotFound         = errors.New("oopsy, id not found")
	ErrUnknown            = errors.New("oopsy, unknown error has occured")
	ErrEmptyUpdateContent = errors.New("oopsy, update content empty")
	ErrEmptyGetContent    = errors.New("oopsy, get content empty")
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

type Pagination struct {
	Limit  int
	Offset int
}

type GetAllResponse struct {
	List []GetResponse
}

type Provider interface {
	Get(context.Context, *GetRequest) (GetResponse, error)
	GetAll(context.Context, *Pagination) (GetAllResponse, error)

	Create(context.Context, *CreateRequest) (int, error)
	Update(context.Context, *UpdateRequest) (int, error)
	Delete(context.Context, int) error
}

type ResourceProvider struct {
	db port.DB
}

func NewResource(DB port.DB) Provider {
	return &ResourceProvider{
		db: DB,
	}
}

func (r *ResourceProvider) Create(ctx context.Context, req *CreateRequest) (int, error) {
	err := r.validateName(ctx, req.Name)
	if err != nil {
		return 0, err
	}
	err = r.validateAction(ctx, req.Action)
	if err != nil {
		return 0, err
	}
	id, err := r.db.Create(ctx, &port.CreateRequest{
		Action: req.Action,
		Name:   req.Name,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}
	return id, nil
}

func (r *ResourceProvider) Update(ctx context.Context, req *UpdateRequest) (int, error) {
	err := r.validateId(ctx, req.Id)
	if err != nil {
		return 0, err
	}

	//either id or action need tobe updated
	if req.Action == "" && req.Name == "" {
		return 0, ErrEmptyUpdateContent
	}

	//update action
	if req.Action != "" {
		err = r.validateAction(ctx, req.Action)
		if err != nil {
			return 0, err
		}
		_, err = r.db.UpdateAction(ctx, &port.UpdateActionRequest{
			Id:     req.Id,
			Action: req.Action,
		})
		if err != nil {
			switch err {
			default:
				return 0, ErrUnknown
			}
		}
	}

	//update name
	if req.Name != "" {
		err = r.validateName(ctx, req.Name)
		if err != nil {
			return 0, err
		}
		_, err = r.db.UpdateName(ctx, &port.UpdateNameRequest{
			Id:   req.Id,
			Name: req.Name,
		})
		if err != nil {
			switch err {
			default:
				return 0, ErrUnknown
			}
		}
	}

	return req.Id, nil
}

func (r *ResourceProvider) Delete(ctx context.Context, id int) error {
	//validate Id
	err := r.validateId(ctx, id)
	if err != nil {
		return err
	}
	err = r.db.Delete(ctx, id)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (r *ResourceProvider) Get(ctx context.Context, req *GetRequest) (GetResponse, error) {
	// TODO: split to get by param and get by Id
	//either id or name need tobe provided
	if req.Id == 0 && req.Name == "" {
		return GetResponse{}, ErrEmptyGetContent
	}
	//Get by Id
	if req.Id != 0 && req.Name == "" {
		resp, err := r.db.GetByID(ctx, req.Id)
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
		return GetResponse{
			Id:     resp.Id,
			Action: resp.Action,
			Name:   resp.Name,
		}, nil
	}
	//Get by Name
	if req.Name != "" && req.Id == 0 {
		resp, err := r.db.GetByName(ctx, req.Name)
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
		return GetResponse{
			Id:     resp.Id,
			Action: resp.Action,
			Name:   resp.Name,
		}, nil
	}
	//Get by Name and Id
	if req.Name != "" && req.Id != 0 {
		resultByName, err := r.db.GetByName(ctx, req.Name)
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}

		resultById, err := r.db.GetByID(ctx, req.Id)
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
		if resultByName.Id == resultById.Id {
			return GetResponse{
				Id:     resultById.Id,
				Action: resultById.Action,
				Name:   resultById.Name,
			}, nil
		}
	}
	return GetResponse{}, nil
}

func (r *ResourceProvider) GetAll(ctx context.Context, pagination *Pagination) (GetAllResponse, error) {

	allResources, err := r.db.GetAll(ctx, &port.Pagination{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	})
	if err != nil {
		switch err {
		default:
			return GetAllResponse{}, nil
		}
	}
	var response GetAllResponse
	for _, i := range allResources.List {
		response.List = append(response.List, GetResponse{
			Id:     i.Id,
			Action: i.Action,
			Name:   i.Name,
		})
	}
	return response, nil
}
