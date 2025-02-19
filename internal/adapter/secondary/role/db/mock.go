package adapter

import (
	"context"

	port "b2b.nati011.github.com/internal/port/role"
)

type MockResource struct {
	Id   int
	Name string
	Desc string
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
				Id:   i.Id,
				Name: i.Name,
				Desc: i.Desc,
			}, nil
		}
	}
	return port.GetResponse{}, nil
}

func (p *Mock) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	for _, i := range p.resources {
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

func (p *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range p.resources {
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
	newResourceId := len(p.resources) + 1
	p.resources = append(p.resources, MockResource{
		Id:   newResourceId,
		Name: req.Name,
		Desc: req.Desc,
	})
	return newResourceId, nil
}

func (p *Mock) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) (int, error) {
	updatedResources := []MockResource{}
	var updatedResourceId int
	for _, i := range p.resources {
		if req.Id == i.Id {
			updatedResourceId = i.Id
			updatedResources = append(updatedResources, MockResource{
				Id:   req.Id,
				Desc: req.Desc,
				Name: i.Name,
			})
		} else {
			updatedResources = append(updatedResources, MockResource{
				Id:   i.Id,
				Name: i.Name,
				Desc: i.Desc,
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
				Id:   req.Id,
				Desc: i.Desc,
				Name: req.Name,
			})
		} else {
			updatedResources = append(updatedResources, MockResource{
				Id:   i.Id,
				Name: i.Name,
				Desc: i.Desc,
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
				Id:   i.Id,
				Name: i.Name,
				Desc: i.Desc,
			})
		}

	}
	p.resources = updatedResources
	return nil
}
