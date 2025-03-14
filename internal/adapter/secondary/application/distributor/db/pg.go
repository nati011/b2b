package db

import (
	"context"

	"database/sql"

	port "b2b.nati011.github.com/internal/port/application/distributor"
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

	err := p.db.QueryRowContext(ctx, query).Scan(response)
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

	err := p.db.QueryRowContext(ctx, query, id).Scan(&response)
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

func (p *Postgres) UpdateBusiness(ctx context.Context, req *port.UpdateBusinessRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_distributor_business_location($1, $2, $3, $4, $5, $6);"

	err := p.db.QueryRowContext(ctx, query,
		req.Id,
		req.Name,
		req.Tin,
		req.DistributorId,
		req.Location.GeneralZone,
		req.Location.Woreda,
		req.Location.Region,
	).Scan(&resourceId)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return resourceId, port.ErrSysNoRows
		default:
			return resourceId, port.ErrSysUnknown
		}
	}

	return resourceId, nil

}

func (p *Postgres) CreateBusiness(ctx context.Context, req *port.CreateBusinessInformation) (port.CreateBusinessResponse, error) {
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

func (p *Postgres) GetBusinessAll(ctx context.Context) (resp port.GetAllResponse, err error) {

	query := "SELECT * FROM public.get_all_distributors();"

	err = p.db.QueryRowContext(ctx, query).Scan(resp)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}
	}

	return resp, nil
}

func (p *Postgres) GetBusinessById(ctx context.Context, id int) (port.GetBusinessResponse, error) {
	var response port.GetBusinessResponse

	query := "SELECT * FROM public.get_distributor_by_id($1);"

	err := p.db.QueryRowContext(ctx, query, id).Scan(&response)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetBusinessResponse{}, port.ErrSysNoRows
		default:
			return port.GetBusinessResponse{}, port.ErrSysUnknown
		}
	}

	return response, nil
}
func (p *Postgres) GetByDistributorId(ctx context.Context, distributor_id int) (port.GetBusinessResponse, error) {
	var response port.GetBusinessResponse

	query := "SELECT * FROM public.get_distributor_business($1);"

	err := p.db.QueryRowContext(ctx, query, distributor_id).Scan(&response)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetBusinessResponse{}, port.ErrSysNoRows
		default:
			return port.GetBusinessResponse{}, port.ErrSysUnknown
		}
	}

	return response, nil
}

func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{
		db: db,
	}
}
