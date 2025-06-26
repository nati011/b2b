package resource

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/application/resource"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

var (
	ErrDuplicateName      = errors.New(" duplicate name")
	ErrEmptyAction        = errors.New(" empty action")
	ErrEmptyName          = errors.New(" empty name")
	ErrIdNotFound         = errors.New(" id not found")
	ErrNameNotFound       = errors.New(" name not found")
	ErrUnknown            = errors.New(" unknown error has occured")
	ErrEmptyUpdateContent = errors.New(" update content empty")
	ErrEmptyGetContent    = errors.New(" get content empty")
)

const (
	ANY   = "ANY"
	READ  = "READ"
	WRITE = "WRITE"
)

type CreateRequest struct {
	Action   string
	Resource string
	Name     string
}

type UpdateRequest struct {
	Id       int
	Action   string
	Resource string
	Name     string
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

type Provider interface {
	GetByName(context.Context, string) (GetResponse, error)
	Get(context.Context, int) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)

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

	//update resource
	if req.Resource != "" {
		err = r.validateResource(ctx, req.Resource)
		if err != nil {
			return 0, err
		}
		_, err = r.db.UpdateResource(ctx, &port.UpdateResouceRequest{
			Id:       req.Id,
			Resource: req.Resource,
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

func (r *ResourceProvider) GetByName(ctx context.Context, name string) (GetResponse, error) {
	resp, err := r.db.GetByName(ctx, name)
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

func (r *ResourceProvider) GetByResource(ctx context.Context, resource string) (GetResponse, error) {
	resp, err := r.db.GetByResource(ctx, resource)
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

func (r *ResourceProvider) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := r.db.GetByID(ctx, id)
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

func (r *ResourceProvider) GetAll(ctx context.Context) (GetAllResponse, error) {
	allResources, err := r.db.GetAll(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetContent
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}
	var response GetAllResponse
	for _, i := range allResources.List {
		response.List = append(response.List, GetResponse(i))
	}
	return response, nil
}
