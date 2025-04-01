package adapter

import (
	"context"
	"database/sql"
	"log"
	"strings"

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

	err := r.Pool.QueryRowContext(ctx, query,
		req.Name,
		req.Tin,
		req.Latitude,
		req.Longitude,
		req.GeneralZone,
		req.Region,
		req.Woreda,
	).Scan(&retailerId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}
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

	_, err := r.Pool.QueryContext(ctx, query, req.Retailer_id, req.User_id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	return nil
}

func (r *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_retailer_by_id($1);"
	err := r.Pool.QueryRowContext(ctx, query, id).Scan(&response.Id, &response.Name, &response.Tin, &response.Latitude, &response.Longitude, &response.GeneralZone, &response.Region, &response.Woreda)
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

func (r *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	query := "SELECT * FROM public.update_retailer_name($1, $2);"
	_, err := r.Pool.QueryContext(ctx, query, req.Id, req.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	return nil
}

func (r *Postgres) UpdateTin(ctx context.Context, req *port.UpdateTinRequest) error {
	query := "SELECT * FROM public.update_retailer_tin($1, $2);"
	_, err := r.Pool.QueryContext(ctx, query, req.Id, req.Tin)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}
	return nil
}

func (r *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_retailers();"
	rows, err := r.Pool.QueryContext(ctx, query)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var retailer port.GetResponse
		if err := rows.Scan(&retailer.Id, &retailer.Name, &retailer.Tin, &retailer.Latitude, &retailer.Longitude, &retailer.GeneralZone, &retailer.Region, &retailer.Woreda); err != nil {
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

	query := "SELECT * FROM public.get_retailer_by_name($1);"
	rows, err := r.Pool.QueryContext(ctx, query, strings.Trim(name, `"`))
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var retailer port.GetResponse
		if err := rows.Scan(&retailer.Id, &retailer.Name, &retailer.Tin, &retailer.Latitude, &retailer.Longitude, &retailer.GeneralZone, &retailer.Region, &retailer.Woreda); err != nil {
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
	query := "SELECT * FROM public.get_retailer_by_tin($1);"
	err := r.Pool.QueryRowContext(ctx, query, strings.Trim(tin, `"`)).Scan(&response.Id, &response.Name, &response.Tin, &response.Latitude, &response.Longitude, &response.GeneralZone, &response.Region, &response.Woreda)
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

func (r *Postgres) GetAllUserAgents(ctx context.Context, id int) (port.GetAllUserResponse, error) {
	var response port.GetAllUserResponse

	query := "SELECT * FROM public.get_all_retailer_users($1);"
	rows, err := r.Pool.QueryContext(ctx, query, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllUserResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllUserResponse{}, port.ErrSysUnknown
		}

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
