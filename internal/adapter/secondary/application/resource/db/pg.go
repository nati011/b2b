package adapter

import (
	"context"
	"database/sql"
	"log"

	"b2b.nati011.github.com/config"
	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/resource"
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

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_resources_by_id($1);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		id,
	)
	if err != nil {
		return port.GetResponse{}, err
	}

	rows.Row.Scan(
		&response.Id,
		&response.Action,
		&response.Name,
	)
	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_resources_by_name($1);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query, false,
		name,
	)
	if err != nil {
		return port.GetResponse{}, err
	}

	rows.Row.Scan(
		&response.Id,
		&response.Action,
		&response.Name,
	)

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	limit := p.Pagination.Limit
	offset := p.Pagination.Offset

	query := "SELECT * FROM public.get_all_resources($1,$2);"

	rows, err := handler.MustQueryRow(
		p.Pool,
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
		var resource port.GetResponse
		if err := rows.Rows.Scan(&resource.Id, &resource.Action, &resource.Name); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, resource)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_resource($1, $2);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Name,
		req.Action,
	)
	if err != nil {
		return 0, err
	}

	rows.Row.Scan(&resourceId)
	return resourceId, nil
}

func (p *Postgres) UpdateAction(ctx context.Context, req *port.UpdateActionRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_resource_action($1, $2);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Id,
		req.Action,
	)

	if err != nil {
		return 0, err
	}

	rows.Row.Scan(&resourceId)

	return resourceId, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_resource_name($1, $2);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Id,
		req.Name,
	)

	if err != nil {
		return 0, err
	}

	rows.Row.Scan(&resourceId)
	return resourceId, nil
}

func (p *Postgres) Delete(ctx context.Context, id int) error {
	query := "SELECT * FROM public.delete_resource($1);"

	_, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
