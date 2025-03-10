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
	var resourceId int
	var resp port.CreateBusinessResponse
	query := "SELECT * FROM public.create_distributor_business_location($1, $2, $3, $4, $5, $6);"

	err := p.db.QueryRowContext(ctx, query,
		req.Name,
		req.Tin,
		req.DistributorId,
		req.GeneralZone,
		req.Region,
		req.Woreda,
	).Scan(&resourceId)

	if err != nil {
		print(err.Error())
		switch err {
		case sql.ErrNoRows:
			return resp, port.ErrSysNoRows
		default:
			return resp, port.ErrSysUnknown
		}
	}
	resp = port.CreateBusinessResponse{
		BusinessId: resourceId,
	}
	return resp, nil
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
