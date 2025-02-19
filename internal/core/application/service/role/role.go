package role

import (
	"context"
	"errors"

	port "b2b.nati011.github.com/internal/port/role"
)

var (
	ErrDuplicateName      = errors.New("oopsy, duplicate name")
	ErrEmptyName          = errors.New("oopsy, empty name")
	ErrIdNotFound         = errors.New("oopsy, id not found")
	ErrUnknown            = errors.New("oopsy, unknown error has occured")
	ErrEmptyUpdateContent = errors.New("oopsy, update content empty")
	ErrEmptyGetContent    = errors.New("oopsy, get content empty")
)

type CreateRequest struct {
	Name string
	Desc string
}

type UpdateRequest struct {
	Id   int
	Name string
	Desc string
}

type GetRequest struct {
	Id   int
	Name string
}

type GetResponse struct {
	Id   int
	Name string
	Desc string
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

	HasResource(context.Context, int) (bool, error)
}

type RoleProvider struct {
	db port.DB
}

func NewRole(DB port.DB) Provider {
	return &RoleProvider{
		db: DB,
	}
}

func (r *RoleProvider) Get(ctx context.Context, req *GetRequest) (GetResponse, error) {
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
			Id:   resp.Id,
			Desc: resp.Desc,
			Name: resp.Name,
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
			Id:   resp.Id,
			Desc: resp.Desc,
			Name: resp.Name,
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
				Id:   resultById.Id,
				Desc: resultById.Desc,
				Name: resultById.Name,
			}, nil
		}
	}
	return GetResponse{}, nil
}

func (r *RoleProvider) GetAll(ctx context.Context) (GetAllResponse, error) {
	allResources, err := r.db.GetAll(ctx)
	if err != nil {
		switch err {
		default:
			return GetAllResponse{}, nil
		}
	}
	var response GetAllResponse
	for _, i := range allResources.List {
		response.List = append(response.List, GetResponse{
			Id:   i.Id,
			Desc: i.Desc,
			Name: i.Name,
		})
	}
	return response, nil
}

func (r *RoleProvider) Create(ctx context.Context, req *CreateRequest) (int, error) {
	err := r.validateName(ctx, req.Name)
	if err != nil {
		return 0, err
	}
	id, err := r.db.Create(ctx, &port.CreateRequest{
		Desc: req.Desc,
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

func (r *RoleProvider) Update(ctx context.Context, req *UpdateRequest) (int, error) {
	err := r.validateId(ctx, req.Id)
	if err != nil {
		return 0, err
	}

	//either id or desc need tobe updated
	if req.Desc == "" && req.Name == "" {
		return 0, ErrEmptyUpdateContent
	}

	//update desc
	if req.Desc != "" {
		_, err = r.db.UpdateDesc(ctx, &port.UpdateDescRequest{
			Id:   req.Id,
			Desc: req.Desc,
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

func (r *RoleProvider) Delete(ctx context.Context, id int) error {
	// validate Id
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

func (r *RoleProvider) HasResource(ctx context.Context, id int) (bool, error) {
	response, err := r.Get(ctx, &GetRequest{
		Id: id,
	})
	if err != nil {
		switch err {
		default:
			return false, ErrUnknown
		}
	}
	if response.Id == 0 {
		return false, nil
	}
	return true, nil
}
