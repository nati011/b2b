package adapter

import (
	"context"

	port "b2b.nati011.github.com/internal/port/resource"
)

type Mock struct {
}

func NewMock() port.DB {
	return &Mock{}
}

func (p *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Mock) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Mock) UpdateAction(ctx context.Context, req *port.UpdateActionRequest) (int, error) {
	return 0, nil
}

func (p *Mock) UpdateName(ctx context.Context, req *port.UpdateNameRequest) (int, error) {
	return 0, nil
}

func (p *Mock) Delete(ctx context.Context, id int) error {
	return nil
}
