package db

import (
	"context"

	port "b2b.nati011.github.com/internal/port/application/partner/db"
)

type MockPartner struct {
	Id               int
	Name             string
	Icon             string
	Init_payment_url string
	Status           string
	Secret           string
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

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
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

func (m *Mock) GetPartnerSecret(ctx context.Context, id int) (port.GetPartnerSecret, error) {
	for _, i := range m.resources {
		if i.Id == id {
			return port.GetPartnerSecret{
				Name:             i.Name,
				Init_payment_url: i.Init_payment_url,
				Secret:           i.Secret,
			}, nil
		}
	}
	return port.GetPartnerSecret{}, port.ErrSysNoRows
}

func (m *Mock) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
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

func (m *Mock) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	response := []port.GetResponse{}
	for _, i := range m.resources {
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
		Id:               newResourceId,
		Name:             req.Name,
		Icon:             req.Icon,
		Init_payment_url: req.Init_payment_url,
		Status:           req.Status,
		Secret:           req.Secret,
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
