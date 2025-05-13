package role

import (
	"context"
	"errors"
	"log"

	resource "b2b.nati011.github.com/internal/core/application/resource"
	port "b2b.nati011.github.com/internal/port/application/role"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

var (
	ErrDuplicateName               = errors.New("oopsy, duplicate name")
	ErrEmptyName                   = errors.New("oopsy, empty name")
	ErrIdNotFound                  = errors.New("oopsy, id not found")
	ErrUnknown                     = errors.New("oopsy, unknown error has occured")
	ErrEmptyUpdateContent          = errors.New("oopsy, update content empty")
	ErrEmptyGetContent             = errors.New("oopsy, get content empty")
	ErrResourceAlreadyExistsInRole = errors.New("oopsy, resource already exists in role")
	ErrResourceNotFound            = errors.New("oopsy, resource not found")
	ErrResourceNotFoundInRole      = errors.New("oopsy, resource not found in role")
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

type GetResourceResponse struct {
	Id     int
	Name   string
	Action string
}

type GetAllResourcesResponse struct {
	List []GetResourceResponse
}

type AddResourceRequest struct {
	ResourceId int
	RoleId     int
}

type RemoveResourceRequest struct {
	ResourceId int
	RoleId     int
}

type HasResourceRequest struct {
	ResourceId int
	RoleId     int
}

type Provider interface {
	Get(context.Context, *GetRequest) (GetResponse, error)
	GetAll(context.Context) (GetAllResponse, error)

	Create(context.Context, *CreateRequest) (int, error)
	Update(context.Context, *UpdateRequest) (int, error)
	Delete(context.Context, int) error

	AddResource(context.Context, *AddResourceRequest) error
	RemoveResource(context.Context, *RemoveResourceRequest) error
	GetAllResources(context.Context, int) (GetAllResourcesResponse, error)
	HasResource(context.Context, *HasResourceRequest) (bool, error)
}

type RoleProvider struct {
	db               port.DB
	resource_service resource.Provider
}

func NewRole(DB port.DB, resource_service resource.Provider) Provider {
	return &RoleProvider{
		db:               DB,
		resource_service: resource_service,
	}
}

func (r *RoleProvider) Get(ctx context.Context, req *GetRequest) (GetResponse, error) {
	log.Printf("Get Request Id %v", req.Id)
	//either id or name need tobe provided
	if req.Id == 0 && req.Name == "" {
		return GetResponse{}, ErrEmptyGetContent
	}

	log.Printf("Get Request Name %v", req.Name)

	//Get by Id
	if req.Id != 0 && req.Name == "" {
		resp, err := r.db.GetByID(ctx, req.Id)
		if err != nil {
			switch err {
			default:
				return GetResponse{}, ErrUnknown
			}
		}
		if resp.Id == req.Id {
			return GetResponse{
				Id:   resp.Id,
				Desc: resp.Desc,
				Name: resp.Name,
			}, nil
		}
	}

	//Get by Name
	if req.Name != "" && req.Id == 0 {
		resp, err := r.db.GetByName(ctx, req.Name)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
				return GetResponse{}, ErrEmptyGetContent
			default:
				return GetResponse{}, ErrUnknown
			}
		}
		if resp.Name == req.Name {
			return GetResponse{
				Id:   resp.Id,
				Desc: resp.Desc,
				Name: resp.Name,
			}, nil
		}
	}

	//Get by Name and Id
	if req.Name != "" && req.Id != 0 {
		resultByName, err := r.db.GetByName(ctx, req.Name)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
				return GetResponse{}, ErrEmptyGetContent
			default:
				return GetResponse{}, ErrUnknown
			}
		}

		resultById, err := r.db.GetByID(ctx, req.Id)
		if err != nil {
			switch err {
			case port_commons.ErrSysNoRows:
				return GetResponse{}, ErrEmptyGetContent
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
	return GetResponse{}, ErrEmptyGetContent
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
		switch err {
		case ErrEmptyGetContent:
		default:
			return 0, err
		}

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
	log.Printf("Deleting id %v", id)
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

func (r *RoleProvider) HasResource(ctx context.Context, req *HasResourceRequest) (bool, error) {
	err := r.validateId(ctx, req.RoleId)
	if err != nil {
		return false, err
	}
	resp, err := r.db.GetAllResources(ctx, req.RoleId)
	if err != nil {
		switch err {
		default:
			return false, ErrUnknown
		}
	}

	for _, i := range resp.List {
		if i.Id == req.ResourceId {
			return true, nil
		}
	}
	return false, nil
}

func (r *RoleProvider) AddResource(ctx context.Context, req *AddResourceRequest) error {
	//check if role exists
	err := r.validateId(ctx, req.RoleId)
	if err != nil {
		return err
	}

	//check if resource exists
	_, err = r.resource_service.Get(ctx, req.ResourceId)
	if err != nil {
		switch err {
		case resource.ErrIdNotFound:
			return ErrResourceNotFound
		default:
			return ErrUnknown
		}
	}

	//check if resource already exists in role
	resp, err := r.db.GetAllResources(ctx, req.RoleId)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	for _, i := range resp.List {
		if i.Id == req.ResourceId {
			return ErrResourceAlreadyExistsInRole
		}
	}

	//add role to resource
	err = r.db.AddResource(ctx, req.RoleId, req.ResourceId)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (r *RoleProvider) GetAllResources(ctx context.Context, id int) (GetAllResourcesResponse, error) {
	var response GetAllResourcesResponse
	//check if role exists
	err := r.validateId(ctx, id)
	if err != nil {
		return GetAllResourcesResponse{}, err
	}

	//check if resource already exists in role
	resp, err := r.db.GetAllResources(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
		default:
			return GetAllResourcesResponse{}, ErrUnknown
		}
	}
	for _, i := range resp.List {
		response.List = append(response.List, GetResourceResponse{
			Id:     i.Id,
			Name:   i.Name,
			Action: i.Action,
		})
	}
	return response, nil
}

func (r *RoleProvider) RemoveResource(ctx context.Context, req *RemoveResourceRequest) error {
	//check if role exists
	err := r.validateId(ctx, req.RoleId)
	if err != nil {
		return err
	}

	//check if resource exists
	resResp, err := r.resource_service.Get(ctx, req.ResourceId)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	if resResp.Id == 0 {
		return ErrResourceNotFound
	}

	//check if resource exists in role
	resp, err := r.db.GetAllResources(ctx, req.RoleId)
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	hasResource := false
	for _, i := range resp.List {
		if i.Id == req.ResourceId {
			hasResource = true
		}
	}
	if !hasResource {
		return ErrResourceNotFoundInRole
	}

	//remove role
	err = r.db.RemoveResource(ctx, req.RoleId, req.ResourceId)
	if err != nil {
		return ErrUnknown
	}
	return nil
}
