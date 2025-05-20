package email

import (
	"context"

	port "b2b.nati011.github.com/internal/port/application/payment/db"
)

type Mock struct {
}

func NewMock() port.DB {
	return &Mock{}
}

func (p *Mock) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Mock) GetByTransactionRef(ctx context.Context, name string) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Mock) GetByOrderId(ctx context.Context, orderId int) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Mock) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}
