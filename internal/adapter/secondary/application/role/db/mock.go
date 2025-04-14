package adapter

import (
	"context"

	port "b2b.nati011.github.com/internal/port/application/role"
)

type MockRole struct {
	Id   int
	Name string
	Desc string
}

type MockResource struct {
	roleId     int
	resourceId int
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
	return port.GetResponse{}, nil
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
	return port.GetResponse{}, nil
}

func (p *Mock) GetAll(ctx context.Context, pagination *port.Pagination) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	count := 0
	for _, i := range p.roles {
		count++
		if count > pagination.Limit {
			break
		}
		response = append(response, port.GetResponse{
			Id:   i.Id,
			Name: i.Name,
			Desc: i.Desc,
		})
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
	var resp []int
	for _, i := range p.resources {
		if i.roleId == role_id {
			resp = append(resp, i.resourceId)
		}
	}
	return port.GetAllResourcesResponse{
		List: resp,
	}, nil
}
