package email

import (
	"context"
	"database/sql"

	port "b2b.nati011.github.com/internal/port/application/payment/db"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{Pool: db}
}

func (p *Postgres) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetByTransactionRef(ctx context.Context, name string) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}
