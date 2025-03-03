package db

import (
	"context"
	"math/rand"

	port "b2b.nati011.github.com/internal/port/distributor"
)

type MockDistributor struct {
	Id int
}

type Mock struct {
	Distributors []MockDistributor
}

func (m *Mock) GetById(ctx context.Context, id int) (port.GetResponse, error) {
	resp := port.GetResponse{}
	for _, i := range m.Distributors {
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
	distributorId := rand.Int()
	m.Distributors = append(m.Distributors, MockDistributor{
		Id: distributorId,
	})
	return distributorId, nil

}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	resp := []port.GetResponse{}
	for _, i := range m.Distributors {

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

func NewMock() port.DB {
	return &Mock{}
}
