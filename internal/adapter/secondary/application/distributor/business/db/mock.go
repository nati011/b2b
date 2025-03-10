package db

import (
	"context"
	"math/rand"

	port "b2b.nati011.github.com/internal/port/distributor/business"
)

type MockBusiness struct {
	Id            int
	Name          string
	Tin           string
	DistributorId int
}

type Mock struct {
	Businesses []MockBusiness
}

func (m *Mock) GetById(ctx context.Context, id int) (port.GetResponse, error) {
	resp := port.GetResponse{}
	for _, i := range m.Businesses {
		resp = port.GetResponse{
			Id: i.Id,
		}
	}
	if (resp == port.GetResponse{}) {
		return port.GetResponse{}, port.ErrSysNoRows
	}

	return resp, nil
}

func (m *Mock) GetByDistributorId(ctx context.Context, distributorId int) (port.GetResponse, error) {
	resp := port.GetResponse{}
	for _, i := range m.Businesses {
		resp = port.GetResponse{
			DistributorId: i.DistributorId,
		}
	}
	if (resp == port.GetResponse{}) {
		return port.GetResponse{}, port.ErrSysNoRows
	}

	return resp, nil
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.Businesses {

		resp = append(resp, port.GetResponse{
			Id: i.Id,
		})
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port.ErrSysNoRows
	}

	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) Create(ctx context.Context, req *port.CreateBusinessInformation) (port.CreateBusinessResponse, error) {
	businessId := rand.Int()
	m.Businesses = append(m.Businesses, MockBusiness{
		Id:            businessId,
		Name:          req.Name,
		Tin:           req.Tin,
		DistributorId: req.DistributorId,
	})

	response := port.CreateBusinessResponse{
		BusinessId: businessId,
	}
	return response, nil
}

func (m *Mock) Update(ctx context.Context, req *port.CreateBusinessInformation) (port.CreateBusinessResponse, error) {
	businessId := rand.Int()
	m.Businesses = append(m.Businesses, MockBusiness{
		Id:            businessId,
		Name:          req.Name,
		Tin:           req.Tin,
		DistributorId: req.DistributorId,
	})

	response := port.CreateBusinessResponse{
		BusinessId: businessId,
	}
	return response, nil

}

func NewMock() port.DB {
	return &Mock{}
}
