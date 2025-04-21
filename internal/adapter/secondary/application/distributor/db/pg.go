package db

import (
	"context"
	"database/sql"
	"log"

	"b2b.nati011.github.com/config"
	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/distributor"
)

type Postgres struct {
	Pool       *sql.DB
	Pagination *config.Pagination
}

func NewPostgres(DB *sql.DB, pagination *config.Pagination) port.DB {
	return &Postgres{
		Pool:       DB,
		Pagination: pagination,
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
		false,
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

	rows.Row.Scan(&distributorId)
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
		false,
		id,
	)
	if err != nil {
		return port.GetResponse{}, err
	}

	rows.Row.Scan(&response.Id,
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
		false,
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
		false,
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

	limit := r.Pagination.Limit
	offset := r.Pagination.Offset

	query := "SELECT * FROM public.get_all_distributors($1,$2);"
	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		true,
		limit,
		offset,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var retailer port.GetResponse
		if err := rows.Rows.Scan(
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

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (r *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	limit := r.Pagination.Limit
	offset := r.Pagination.Offset

	query := "SELECT * FROM public.get_distributor_by_name($1,$2,$3);"
	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		true,
		name,
		limit,
		offset,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var retailer port.GetResponse
		if err := rows.Rows.Scan(
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

	if err := rows.Rows.Err(); err != nil {
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
		false,
		tin,
	)
	if err != nil {
		return port.GetResponse{}, err
	}
	rows.Row.Scan(
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
		true,
		id,
	)
	if err != nil {
		return port.GetAllUserResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var retailer_user port.GetUserResponse
		if err := rows.Rows.Scan(&retailer_user.Id); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllUserResponse{}, err
		}
		response.List = append(response.List, retailer_user)
	}

	if err := rows.Rows.Err(); err != nil {
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
		false,
		req.Distributor_Id,
		req.User_id,
	)
	if err != nil {
		return 0, err
	}
	rows.Row.Scan(&retailer_id)

	return retailer_id, nil
}
