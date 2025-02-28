package distributor

import (
	"context"

	port "b2b.nati011.github.com/internal/port/distributor"
)

type MockDistributor struct {
	Id int
}

type Mock struct {
	Distributors []MockDistributor
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	panic("unimplemented")
}

func (m *Mock) Get(ctx context.Context, id int) (port.GetResponse, error) {
	panic("unimplemented")
}

func (m *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	panic("unimplemented")
}

func (m *Mock) GetById(ctx context.Context, req *port.GetByIdRequest) (port.GetAllResponse, error) {
	panic("unimplemented")
}

func NewMock() port.DB {
	return &Mock{}
}
