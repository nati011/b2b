package category

import (
	"context"
	"database/sql"
	"log"

	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/category"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_category();"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
	)

	if err != nil {
		return port.GetAllResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var resource port.GetResponse
		if err := rows.Scan(&resource.Id, &resource.Name); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, resource)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_category($1);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		id,
	)

	if err != nil {
		return port.GetResponse{}, err
	}

	rows.Scan(&response.Id, &response.Name)

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_category($1);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		req.Name,
	)

	rows.Scan(&resourceId)
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

func (p *Postgres) Remove(ctx context.Context, id int) error {
	query := "SELECT * FROM public.remove_category($1);"
	_, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}
	return nil
}
