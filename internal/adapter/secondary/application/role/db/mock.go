package adapter

import (
	"context"

	port "b2b.nati011.github.com/internal/port/application/role"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

type MockRole struct {
	Id   int
	Name string
	Desc string
}

type MockResource struct {
	roleId     int
	resourceId int
	Name       string
	Action     string
}

type Mock struct {
	roles     []MockRole
	resources []MockResource
}

func NewMock() port.DB {
	return &Mock{}
}

func (p *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range p.roles {
		if i.Id == id {
			return port.GetResponse{
				Id:   i.Id,
				Name: i.Name,
				Desc: i.Desc,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (p *Mock) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	for _, i := range p.roles {
		if i.Name == name {
			return port.GetResponse{
				Id:   i.Id,
				Name: i.Name,
				Desc: i.Desc,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (p *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range p.roles {
		response = append(response, port.GetResponse{
			Id:   i.Id,
			Name: i.Name,
			Desc: i.Desc,
		})
	}
	if len(response) == 0 {
		return port.GetAllResponse{
			List: response,
		}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (p *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	newRoleId := len(p.roles) + 1
	p.roles = append(p.roles, MockRole{
		Id:   newRoleId,
		Name: req.Name,
		Desc: req.Desc,
	})
	return newRoleId, nil
}

func (p *Mock) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) (int, error) {
	updatedRoles := []MockRole{}
	var updatedResourceId int
	for _, i := range p.roles {
		if req.Id == i.Id {
			updatedResourceId = i.Id
			updatedRoles = append(updatedRoles, MockRole{
				Id:   req.Id,
				Desc: req.Desc,
				Name: i.Name,
			})
		} else {
			updatedRoles = append(updatedRoles, MockRole{
				Id:   i.Id,
				Name: i.Name,
				Desc: i.Desc,
			})
		}

	}
	p.roles = updatedRoles
	return updatedResourceId, nil
}

func (p *Mock) UpdateName(ctx context.Context, req *port.UpdateNameRequest) (int, error) {
	updatedResources := []MockRole{}
	var updatedResourceId int
	for _, i := range p.roles {
		if req.Id == i.Id {
			updatedResourceId = i.Id
			updatedResources = append(updatedResources, MockRole{
				Id:   req.Id,
				Desc: i.Desc,
				Name: req.Name,
			})
		} else {
			updatedResources = append(updatedResources, MockRole{
				Id:   i.Id,
				Name: i.Name,
				Desc: i.Desc,
			})
		}

	}
	p.roles = updatedResources
	return updatedResourceId, nil
}

func (p *Mock) Delete(ctx context.Context, id int) error {
	updatedResources := []MockRole{}
	for _, i := range p.roles {
		if i.Id != id {
			updatedResources = append(updatedResources, MockRole{
				Id:   i.Id,
				Name: i.Name,
				Desc: i.Desc,
			})
		}

	}
	p.roles = updatedResources
	return nil
}

func (p *Mock) AddResource(ctx context.Context, role_id int, resource_id int) error {
	p.resources = append(p.resources, MockResource{
		roleId:     role_id,
		resourceId: resource_id,
	})
	return nil
}

func (p *Mock) RemoveResource(ctx context.Context, role_id int, resource_id int) error {
	newResources := []MockResource{}
	for _, i := range p.resources {
		if i.resourceId != resource_id {
			newResources = append(newResources, MockResource{
				roleId:     i.roleId,
				resourceId: i.resourceId,
			})
		}
	}
	p.resources = newResources
	return nil
}

func (p *Mock) GetAllResources(ctx context.Context, role_id int) (port.GetAllResourcesResponse, error) {
	var resp []port.GetResourceResponse
	for _, i := range p.resources {
		if i.roleId == role_id {
			resp = append(resp, port.GetResourceResponse{
				Id: i.resourceId,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResourcesResponse{
			List: resp,
		}, port_commons.ErrSysNoRows
	}
	return port.GetAllResourcesResponse{
		List: resp,
	}, nil
}
