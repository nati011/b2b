package adapter

import (
	"context"
	"database/sql"
	"log"

	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/role"
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

	query := "SELECT * FROM public.get_roles_by_id($1);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		id,
	)

	if err != nil {
		return port.GetResponse{}, err
	}

	log.Printf("Rows %v", rows)

	rows.Scan(&response.Id, &response.Name, &response.Desc)
	log.Printf(
		"Resouce id %v", response.Id,
	)
	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_roles_by_name($1);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		name,
	)

	if err != nil {
		return port.GetResponse{}, err
	}

	rows.Scan(&response.Id, &response.Name, &response.Desc)

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_roles();"

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
		if err := rows.Scan(&resource.Id, &resource.Desc, &resource.Name); err != nil {
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
	query := "SELECT * FROM public.create_role($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, req.Name, req.Desc).Scan(&resourceId)
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

func (p *Postgres) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_role_desc($1, $2);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		req.Id,
		req.Desc,
	)

	if err != nil {
		return 0, err
	}

	rows.Scan(&resourceId)
	return resourceId, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_role_name($1, $2);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		req.Id,
		req.Name,
	)

	if err != nil {
		return 0, err
	}

	rows.Scan(&resourceId)
	return resourceId, nil
}

func (p *Postgres) Delete(ctx context.Context, id int) error {
	query := "SELECT * FROM public.delete_role($1);"
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

func (p *Postgres) AddResource(ctx context.Context, role_id int, resource_id int) error {
	query := "SELECT * FROM public.add_resource_to_role($1, $2);"

	_, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		role_id,
		resource_id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) RemoveResource(ctx context.Context, role_id int, resource_id int) error {
	query := "SELECT * FROM public.remove_resource_from_role($1, $2);"

	_, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		role_id,
		resource_id,
	)

	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetAllResources(ctx context.Context, role_id int) (port.GetAllResourcesResponse, error) {
	var response port.GetAllResourcesResponse

	query := "SELECT * FROM public.get_all_resource_by_role($1);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		role_id,
	)

	if err != nil {
		return port.GetAllResourcesResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var resourceID int
		if err := rows.Scan(&resourceID); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResourcesResponse{}, err
		}
		response.List = append(response.List, resourceID)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResourcesResponse{}, port.ErrSysUnknown
	}

	return response, nil
}
