package db

import (
	"context"

	port "b2b.nati011.github.com/internal/port/application/partner/db"
)

type MockPartner struct {
	Id     int
	Name   string
	Icon   string
	Status string
}

type Mock struct {
	resources []MockPartner
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.resources {
		if i.Id == id {
			return port.GetResponse{
				Id:     i.Id,
				Name:   i.Name,
				Status: i.Status,
				Icon:   i.Icon,
			}, nil
		}
	}
	return port.GetResponse{}, port.ErrSysNoRows
}

func (m *Mock) GetAll(ctx context.Context, pagination *port.Pagination) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	count := 0
	for _, i := range m.resources {
		count++
		if count > pagination.Limit {
			break
		}
		response = append(response, port.GetResponse{
			Id:     i.Id,
			Name:   i.Name,
			Status: i.Status,
			Icon:   i.Icon,
		})
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) GetByStatus(ctx context.Context, status string, pagination *port.Pagination) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	count := 0
	for _, i := range m.resources {
		count++
		if count > pagination.Limit {
			break
		}
		if i.Status == status {
			response = append(response, port.GetResponse{
				Id:     i.Id,
				Name:   i.Name,
				Status: i.Status,
				Icon:   i.Icon,
			})
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) GetByName(ctx context.Context, name string, pagination *port.Pagination) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	count := 0
	for _, i := range m.resources {
		count++
		if count > pagination.Limit {
			break
		}
		if i.Name == name {
			response = append(response, port.GetResponse{
				Id:     i.Id,
				Name:   i.Name,
				Status: i.Status,
				Icon:   i.Icon,
			})
		}
	}
	if len(response) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: response,
	}, nil
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	newResourceId := len(m.resources) + 1
	m.resources = append(m.resources, MockPartner{
		Id:   newResourceId,
		Name: req.Name,
		Icon: req.Icon,
	})
	return newResourceId, nil
}

func (m *Mock) UpdateStatus(ctx context.Context, id int, status string) (int, error) {
	updatedResources := []MockPartner{}
	var updatedResourceId int
	for _, i := range m.resources {
		if i.Id == id {
			updatedResourceId = i.Id
			updatedResources = append(updatedResources, MockPartner{
				Id:     i.Id,
				Name:   i.Name,
				Status: status,
				Icon:   i.Icon,
			})
		} else {
			updatedResources = append(updatedResources, MockPartner{
				Id:     i.Id,
				Name:   i.Name,
				Status: i.Status,
				Icon:   i.Icon,
			})
		}

	}
	m.resources = updatedResources
	return updatedResourceId, nil
}

func (m *Mock) UpdateName(ctx context.Context, id int, name string) (int, error) {
	updatedResources := []MockPartner{}
	var updatedResourceId int
	for _, i := range m.resources {
		if i.Id == id {
			updatedResourceId = i.Id
			updatedResources = append(updatedResources, MockPartner{
				Id:     i.Id,
				Name:   i.Name,
				Status: i.Status,
				Icon:   i.Icon,
			})
		} else {
			updatedResources = append(updatedResources, MockPartner{
				Id:     i.Id,
				Name:   i.Name,
				Status: i.Status,
				Icon:   i.Icon,
			})
		}

	}
	m.resources = updatedResources
	return updatedResourceId, nil
}

func (m *Mock) Delete(ctx context.Context, id int) error {
	updatedResources := []MockPartner{}
	for _, i := range m.resources {
		if i.Id != id {
			updatedResources = append(updatedResources, MockPartner{
				Id:     i.Id,
				Name:   i.Name,
				Status: i.Status,
				Icon:   i.Icon,
			})
		}

	}
	m.resources = updatedResources
	return nil
}
