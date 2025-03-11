package db

import (
	"context"
	"database/sql"

	port "b2b.nati011.github.com/internal/port/application/partner/db"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (p *Postgres) GetByID(context.Context, int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetAll(context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByStatus(context.Context, string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByName(context.Context, string) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) Create(context.Context, *port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateStatus(context.Context, int, string) (int, error) {
	return 0, nil
}

func (p *Postgres) UpdateName(context.Context, int, string) (int, error) {
	return 0, nil
}

func (p *Postgres) Delete(context.Context, int) error {
	return nil
}
