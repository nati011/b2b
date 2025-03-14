package transaction

import (
	"context"
	"database/sql"
	"time"

	port "b2b.nati011.github.com/internal/port/application/transaction/db"
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

func (p *Postgres) GetByDate(context.Context, time.Time) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByUserId(context.Context, int) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) GetByPartnerId(context.Context, int) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}

func (p *Postgres) Create(context.Context, *port.CreateRequest) (int, error) {
	return 0, nil
}

func (p *Postgres) Delete(context.Context, int) error {
	return nil
}
