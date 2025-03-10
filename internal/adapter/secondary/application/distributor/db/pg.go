package db

import (
	"context"

	"database/sql"

	port "b2b.nati011.github.com/internal/port/distributor"
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_distributor_user($1, $2, $3, $4, $5, $6);"

	err := p.db.QueryRowContext(ctx, query,
		req.FirstName,
		req.Email,
		req.Phone,
		req.Username,
		req.DOB,
		req.ExternalId,
	).Scan(&resourceId)

	if err != nil {
		print(err.Error())
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_distributors();"

	err := p.db.QueryRowContext(ctx, query).Scan()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}
	}

	return response, nil
}

func (p *Postgres) GetById(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_distributor_by_id($1);"

	err := p.db.QueryRowContext(ctx, query, id).Scan(&response.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}

	return response, nil
}

func (p *Postgres) Update(ctx context.Context, req *port.CreateBusinessInformation) (port.CreateBusinessResponse, error) {
	panic("unimplemented")
}

func (p *Postgres) CreateBusinessInformation(ctx context.Context, req *port.CreateBusinessInformation) (port.CreateBusinessResponse, error) {
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

func (p *Postgres) GetBusinessAll(ctx context.Context) (port.GetAllResponse, error) {
	panic("Unimplemented")
}

func (p *Postgres) GetBusinessById(ctx context.Context, id int) (port.GetResponse, error) {
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
