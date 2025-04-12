package adapter

import (
	"context"
	"database/sql"
	"log"

	port "b2b.nati011.github.com/internal/port/application/resource"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_resources_by_id($1);"

	// Use Scan to match the number of returned columns
	err := p.Pool.QueryRowContext(ctx, query, id).Scan(&response.Id, &response.Action, &response.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, nil
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}

	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_resources_by_name($1);"

	err := p.Pool.QueryRowContext(ctx, query, name).Scan(&response.Id, &response.Action, &response.Name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, nil
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context, pagination *port.Pagination) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_resources();"
	rows, err := p.Pool.QueryContext(ctx, query)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, nil
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var resource port.GetResponse
		if err := rows.Scan(&resource.Id, &resource.Action, &resource.Name); err != nil {
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

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_resource($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, req.Name, req.Action).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, nil
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) UpdateAction(ctx context.Context, req *port.UpdateActionRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_resource_action($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, req.Id, req.Action).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return resourceId, nil
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_resource_name($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, req.Id, req.Name).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, nil
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return resourceId, nil
}

func (p *Postgres) Delete(ctx context.Context, id int) error {
	query := "SELECT * FROM public.delete_resource($1);"

	err := p.Pool.QueryRowContext(ctx, query, id).Err()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil
		default:
			return port.ErrSysUnknown
		}
	}
	return nil
}
