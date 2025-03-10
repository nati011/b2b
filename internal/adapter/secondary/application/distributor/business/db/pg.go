package db

import (
	"context"

	"database/sql"

	port "b2b.nati011.github.com/internal/port/distributor/business"
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) Update(ctx context.Context, req *port.CreateBusinessInformation) (port.CreateBusinessResponse, error) {
	panic("unimplemented")
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateBusinessInformation) (port.CreateBusinessResponse, error) {
	panic("Unimplemented")
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	panic("Unimplemented")
}

func (p *Postgres) GetById(ctx context.Context, id int) (port.GetResponse, error) {
	panic("Unimplemented")
}
func (p *Postgres) GetByDistributorId(ctx context.Context, id int) (port.GetResponse, error) {
	panic("Unimplemented")
}
func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{
		db: db,
	}
}
