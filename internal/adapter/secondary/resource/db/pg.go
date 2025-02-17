package adapter

import (
	"context"

	port "b2b.nati011.github.com/internal/port/resource"
)

type Postgres struct {
}

func NewPostgres() port.DB {
	return &Postgres{}
}

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateAction(ctx context.Context, req *port.UpdateActionRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) Delete(ctx context.Context, id int) error {
	return nil
}
