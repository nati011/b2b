package adapter

import (
	"context"

	port "b2b.nati011.github.com/internal/port/application/resource"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
)

type MockResource struct {
	Id       int
	Name     string
	Resource string
	Action   string
}

type Mock struct {
	resources []MockResource
}

func NewMock() port.DB {
	return &Mock{}
}

func (p *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range p.resources {
		if i.Id == id {
			return port.GetResponse{
				Id:       i.Id,
				Name:     i.Name,
				Action:   i.Action,
				Resource: i.Resource,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (p *Mock) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	for _, i := range p.resources {
		if i.Name == name {
			return port.GetResponse{
				Id:       i.Id,
				Name:     i.Name,
				Action:   i.Action,
				Resource: i.Resource,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (p *Mock) GetByResource(ctx context.Context, resource string) (port.GetResponse, error) {
	for _, i := range p.resources {
		if i.Resource == resource {
			return port.GetResponse{
				Id:       i.Id,
				Name:     i.Name,
				Action:   i.Action,
				Resource: i.Resource,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (p *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range p.resources {
		response = append(response, port.GetResponse{
			Id:       i.Id,
			Name:     i.Name,
			Action:   i.Action,
			Resource: i.Resource,
		})
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (p *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	newResourceId := len(p.resources) + 1
	p.resources = append(p.resources, MockResource{
		Id:       newResourceId,
		Name:     req.Name,
		Action:   req.Action,
		Resource: req.Resource,
	})
	return newResourceId, nil
}

func (p *Mock) UpdateAction(ctx context.Context, req *port.UpdateActionRequest) (int, error) {
	updatedResources := []MockResource{}
	var updatedResourceId int
	for _, i := range p.resources {
		if req.Id == i.Id {
			updatedResourceId = i.Id
			updatedResources = append(updatedResources, MockResource{
				Id:       req.Id,
				Action:   req.Action,
				Resource: i.Resource,
				Name:     i.Name,
			})
		} else {
			updatedResources = append(updatedResources, MockResource{
				Id:       i.Id,
				Name:     i.Name,
				Resource: i.Resource,
				Action:   i.Action,
			})
		}

	}
	p.resources = updatedResources
	return updatedResourceId, nil
}

func (p *Mock) UpdateName(ctx context.Context, req *port.UpdateNameRequest) (int, error) {
	updatedResources := []MockResource{}
	var updatedResourceId int
	for _, i := range p.resources {
		if req.Id == i.Id {
			updatedResourceId = i.Id
			updatedResources = append(updatedResources, MockResource{
				Id:       req.Id,
				Action:   i.Action,
				Resource: i.Resource,
				Name:     req.Name,
			})
		} else {
			updatedResources = append(updatedResources, MockResource{
				Id:       i.Id,
				Name:     i.Name,
				Resource: i.Resource,
				Action:   i.Action,
			})
		}

	}
	p.resources = updatedResources
	return updatedResourceId, nil
}

func (p *Mock) UpdateResource(ctx context.Context, req *port.UpdateResouceRequest) (int, error) {
	updatedResources := []MockResource{}
	var updatedResourceId int
	for _, i := range p.resources {
		if req.Id == i.Id {
			updatedResourceId = i.Id
			updatedResources = append(updatedResources, MockResource{
				Id:       req.Id,
				Action:   i.Action,
				Name:     i.Name,
				Resource: req.Resource,
			})
		} else {
			updatedResources = append(updatedResources, MockResource{
				Id:       i.Id,
				Name:     i.Name,
				Action:   i.Action,
				Resource: req.Resource,
			})
		}

	}
	p.resources = updatedResources
	return updatedResourceId, nil
}

func (p *Mock) Delete(ctx context.Context, id int) error {
	updatedResources := []MockResource{}
	for _, i := range p.resources {
		if i.Id != id {
			updatedResources = append(updatedResources, MockResource{
				Id:       i.Id,
				Name:     i.Name,
				Action:   i.Action,
				Resource: i.Resource,
			})
		}

	}
	p.resources = updatedResources
	return nil
}
