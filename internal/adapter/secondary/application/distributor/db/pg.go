package db

import (
	"context"
	"database/sql"
	"log"

	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/distributor"
)

type Postgres struct {
	Pool *sql.DB
}

<<<<<<< HEAD
func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{
		Pool: db,
	}
}

func (r *Postgres) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	var distributorId int
	query := "SELECT * FROM public.create_distributor($1, $2, $3, $4, $5, $6, $7);"

	err := r.Pool.QueryRowContext(ctx, query,
		req.Name,
		req.Tin,
		req.Latitude,
		req.Longitude,
		req.GeneralZone,
		req.Region,
		req.Woreda,
	).Scan(&distributorId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		req.Name,
		req.Tin,
		req.Latitude,
		req.Longitude,
		req.GeneralZone,
		req.Region,
		req.Woreda,
	)
	if err != nil {
		return 0, err
	}

	rows.Scan(&distributorId)
	_, err = r.CreateDistributorUser(ctx, &port.CreateUserAgentRequest{
		User_id:        req.UserId,
		Distributor_Id: distributorId,
	})
	if err != nil {
		switch err {
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return distributorId, nil
}

func (r *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_distributor_by_id($1);"

	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		id,
	)
	if err != nil {
		return port.GetResponse{}, err
	}

	rows.Scan(&response.Id,
		&response.Name,
		&response.Tin,
		&response.Latitude,
		&response.Longitude,
		&response.GeneralZone,
		&response.Region,
		&response.Woreda)

	return response, nil
}

func (r *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	query := "SELECT * FROM public.update_distributor_name($1, $2);"
	_, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		req.Id,
		req.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Postgres) UpdateTin(ctx context.Context, req *port.UpdateTinRequest) error {
	query := "SELECT * FROM public.update_distributor_tin($1, $2);"
	_, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		req.Id,
		req.Tin,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_distributors();"
	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var retailer port.GetResponse
		if err := rows.Scan(
			&retailer.Id,
			&retailer.Name,
			&retailer.Tin,
			&retailer.Latitude,
			&retailer.Longitude,
			&retailer.GeneralZone,
			&retailer.Region,
			&retailer.Woreda); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, retailer)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (r *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_distributor_by_name($1);"
	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		name,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var retailer port.GetResponse
		if err := rows.Scan(
			&retailer.Id,
			&retailer.Name,
			&retailer.Tin,
			&retailer.Latitude,
			&retailer.Longitude,
			&retailer.GeneralZone,
			&retailer.Region,
			&retailer.Woreda); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, retailer)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (r *Postgres) GetByTin(ctx context.Context, tin string) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_distributor_by_tin($1);"

	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		tin,
	)
	if err != nil {
		return port.GetResponse{}, err
	}
	rows.Scan(
		&response.Id,
		&response.Name,
		&response.Tin,
		&response.Latitude,
		&response.Longitude,
		&response.GeneralZone,
		&response.Region,
		&response.Woreda)

	return response, nil
}

func (r *Postgres) GetAllUserAgents(ctx context.Context, id int) (port.GetAllUserResponse, error) {
	var response port.GetAllUserResponse

	query := "SELECT * FROM public.get_all_distributor_users($1);"

	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		id,
	)
	if err != nil {
		return port.GetAllUserResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var retailer_user port.GetUserResponse
		if err := rows.Scan(&retailer_user.Id); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllUserResponse{}, err
		}
		response.List = append(response.List, retailer_user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllUserResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (r *Postgres) CreateDistributorUser(ctx context.Context, req *port.CreateUserAgentRequest) (int, error) {
	query := "SELECT * FROM public.create_distributor_user($1, $2);"
	var retailer_id int

	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		req.Distributor_Id,
		req.User_id).Scan(&retailer_id)
=======
func (p *Postgres) Create(ctx context.Context, req *port.RegisterDistributorRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_distributor_user($1, $2, $3, $4, $5, $6, $7);"

	err := p.db.QueryRowContext(ctx, query,
		req.FirstName,
		req.LastName,
		req.Email,
		req.PhoneNumber,
		req.Username,
		req.DOB,
		req.ExternalId,
	).Scan(&resourceId)
>>>>>>> 8f0b9404 (init distributor refactor)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}
<<<<<<< HEAD
	return retailer_id, nil
=======

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
func (p *Postgres) GetBusinessByDistributorId(ctx context.Context, distributor_id int) (port.GetBusinessResponse, error) {
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
>>>>>>> 8f0b9404 (init distributor refactor)
}
