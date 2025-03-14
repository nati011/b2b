package db

import (
	"context"
	"log"
	"math/rand"
	"time"

	port "b2b.nati011.github.com/internal/port/application/retailer"
)

type MockRetailer struct {
	Id         int
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Username   string
	DOB        time.Time
	ExternalId string
}

type MockBusiness struct {
	Id         int
	Name       string
	Tin        int
	RetailerId int
}

type Mock struct {
	Retailers  []MockRetailer
	Businesses []MockBusiness
}

func (m *Mock) GetById(ctx context.Context, id int) (port.GetResponse, error) {
	resp := port.GetResponse{}
	for _, i := range m.Retailers {
		resp = port.GetResponse{
			Id: i.Id,
		}
	}
	if (resp == port.GetResponse{}) {
		return port.GetResponse{}, port.ErrSysNoRows
	}

	return resp, nil
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	retailerId := rand.Int()
	m.Retailers = append(m.Retailers, MockRetailer{
		Id: retailerId,
	})
	return retailerId, nil

}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.Retailers {

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

func (m *Mock) GetBusinessById(ctx context.Context, id int) (resp port.GetBusinessResponse, err error) {
	for _, i := range m.Businesses {
		log.Print(i)
		resp = port.GetBusinessResponse{
			Id: i.Id,
		}
	}

	if (resp == port.GetBusinessResponse{}) {
		return port.GetBusinessResponse{}, port.ErrSysNoRows
	}

	return resp, nil
}

func (m *Mock) GetByRetailerId(ctx context.Context, retailerId int) (port.GetBusinessResponse, error) {
	resp := port.GetBusinessResponse{}
	for _, i := range m.Businesses {
		resp = port.GetBusinessResponse{
			Id: i.RetailerId,
		}
	}
	if (resp == port.GetBusinessResponse{}) {
		return port.GetBusinessResponse{}, port.ErrSysNoRows
	}

	return resp, nil
}

func (m *Mock) GetBusinessAll(ctx context.Context) (port.GetAllResponse, error) {
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

func (m *Mock) CreateBusiness(ctx context.Context, req *port.CreateBusinessInformation) (port.CreateBusinessResponse, error) {
	businessId := rand.Int()
	m.Businesses = append(m.Businesses, MockBusiness{
		Id:         businessId,
		Name:       req.Name,
		Tin:        req.Tin,
		RetailerId: req.RetailerId,
	})

	response := port.CreateBusinessResponse{
		BusinessId: businessId,
	}
	return response, nil
}

func (m *Mock) UpdateBusiness(ctx context.Context, req *port.UpdateBusinessRequest) (int, error) {
	updatedRetailer := []MockBusiness{}
	var updatedResourceId int
	for _, i := range m.Businesses {
		if req.Id == i.Id {
			updatedResourceId = i.Id
			updatedRetailer = append(updatedRetailer, MockBusiness{
				Id:   req.Id,
				Name: req.Name,
				Tin:  req.Tin,
			})
		} else {
			updatedRetailer = append(updatedRetailer, MockBusiness{
				Id:   i.Id,
				Name: i.Name,
				Tin:  i.Tin,
			})
		}

	}
	m.Businesses = updatedRetailer
	return updatedResourceId, nil

}

func NewMock() port.DB {
	return &Mock{}
}
