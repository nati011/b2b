package adapter

import (
	"context"
	"database/sql"
	"log"

	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/retailer"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{
		Pool: db,
	}
}

func (r *Postgres) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	var retailerId int
	query := "SELECT * FROM public.create_retailer($1, $2, $3, $4, $5, $6, $7);"

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

	rows.Row.Scan(&retailerId)

	err = r.CreateRetailerUser(ctx, &port.CreateUserAgentRequest{
		User_id:     req.UserId,
		Retailer_id: retailerId,
	})
	if err != nil {
		switch err {
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return retailerId, nil
}

func (r *Postgres) CreateRetailerUser(ctx context.Context, req *port.CreateUserAgentRequest) error {
	query := "SELECT * FROM public.create_retailer_user($1, $2);"
	_, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		false,
		req.Retailer_id,
		req.User_id,
	)

	return err
}

func (r *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_retailer_by_id($1);"
	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		false,
		id,
	)
	rows.Row.Scan(
		&response.Id,
		&response.Name,
		&response.Tin,
		&response.Latitude,
		&response.Longitude,
		&response.GeneralZone,
		&response.Region,
		&response.Woreda)
	if err != nil {
		return port.GetResponse{}, err
	}
	return response, nil
}

func (r *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	query := "SELECT * FROM public.update_retailer_name($1, $2);"
	_, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		false,
		req.Id,
		req.Name,
	)
	return err
}

func (r *Postgres) UpdateTin(ctx context.Context, req *port.UpdateTinRequest) error {
	query := "SELECT * FROM public.update_retailer_tin($1, $2);"
	_, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		false,
		req.Id,
		req.Tin,
	)
	return err
}

func (r *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_retailers();"

	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		true,
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

	query := "SELECT * FROM public.get_retailer_by_name($1);"
	rows, err := handler.MustQueryRow(
		r.Pool,
		ctx,
		query,
		true,
		name,
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
	query := "SELECT * FROM public.get_retailer_by_tin($1);"
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

	query := "SELECT * FROM public.get_all_retailer_users($1);"
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
	if len(response.List) == 0 {
		return response, port.ErrSysNoRows
	}
	return response, nil
}
