package category

import (
	"context"
	"database/sql"

	port "b2b.nati011.github.com/internal/port/domain/category"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) Remove(ctx context.Context, id int) error {
	return nil
}
