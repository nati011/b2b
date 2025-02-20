package adapter

import (
	"context"
	"database/sql"
	"log"

	resource "b2b.nati011.github.com/internal/core/application/service/resource"
	port "b2b.nati011.github.com/internal/port/role"
)

type Postgres struct {
	Pool             *sql.DB
	resource_service resource.Provider
}

func NewPostgres(DB *sql.DB, resource_service resource.Provider) port.DB {
	return &Postgres{
		Pool:             DB,
		resource_service: resource_service,
	}
}

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	// Adjust the query to select the appropriate fields
	query := "SELECT * FROM public.get_roles_by_id($1);"

	// Use Scan to match the number of returned columns
	err := p.Pool.QueryRowContext(ctx, query, id).Scan(&response.Id, &response.Name, &response.Desc)
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
	query := "SELECT * FROM public.get_roles_by_name($1);"

	err := p.Pool.QueryRowContext(ctx, query, name).Scan(&response.Id, &response.Name, &response.Desc)
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

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_roles();"
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

	err := p.Pool.QueryRowContext(ctx, query, req.Id, req.Desc).Scan(&resourceId)
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
	query := "SELECT * FROM public.update_role_name($1, $2);"

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
	query := "SELECT * FROM public.delete_role($1);"

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

func (p *Postgres) AddResource(ctx context.Context, role_id int, resource_id int) error {
	query := "SELECT * FROM public.add_resource_to_role($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, role_id, resource_id).Err()
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

func (p *Postgres) RemoveResource(ctx context.Context, role_id int, resource_id int) error {
	query := "SELECT * FROM public.remove_resource_from_role($1, $2);"

	err := p.Pool.QueryRowContext(ctx, query, role_id, resource_id).Err()
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

func (p *Postgres) GetAllResources(ctx context.Context, role_id int) (port.GetAllResourcesResponse, error) {
	query := "SELECT * FROM public.remove_resource_from_role($1);"

	err := p.Pool.QueryRowContext(ctx, query, role_id).Err()
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResourcesResponse{}, nil
		default:
			return port.GetAllResourcesResponse{}, port.ErrSysUnknown
		}
	}
	return port.GetAllResourcesResponse{}, nil
}
