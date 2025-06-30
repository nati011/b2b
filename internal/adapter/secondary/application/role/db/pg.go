package adapter

import (
	"context"
	"database/sql"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/role"
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

	query := "SELECT * FROM public.get_roles_by_id($1);"

	result := []any{&response.Id, &response.Name, &response.Desc}
	args := []any{&id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetResponse{}, err
	}

	response.Id = *result[0].(*int)
	response.Name = *result[1].(*string)
	response.Desc = *result[2].(*string)
	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_roles_by_name($1);"

	result := []any{&response.Id, &response.Name, &response.Desc}
	args := []any{name}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetResponse{}, err
	}

	response.Id = *result[0].(*int)
	response.Name = *result[1].(*string)
	response.Desc = *result[2].(*string)

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_all_roles($1,$2);"

	dest := []any{&responseBase.Id, &responseBase.Name, &responseBase.Desc}
	args := []any{p.Pagination.Limit, p.Pagination.Offset}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		responseBase := port.GetResponse{
			Id:   int(res[0].(int64)),
			Name: res[1].(string),
			Desc: res[2].(string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_role($1, $2);"

	result := []any{&resourceId}
	args := []any{req.Name, req.Desc}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	return resourceId, nil
}

func (p *Postgres) UpdateDesc(ctx context.Context, req *port.UpdateDescRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_role_desc($1, $2);"

	result := []any{&resourceId}
	args := []any{req.Id, req.Desc}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}
	return resourceId, nil
}

func (p *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_role_name($1, $2);"

	result := []any{&resourceId}
	args := []any{req.Id, req.Name}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}
	return resourceId, nil
}

func (p *Postgres) Delete(ctx context.Context, id int) error {
	query := "SELECT * FROM public.delete_role($1);"

	args := []any{id}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) AddResource(ctx context.Context, role_id int, resource_id int) error {
	query := "SELECT * FROM public.add_resource_to_role($1, $2);"

	args := []any{role_id, resource_id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) RemoveResource(ctx context.Context, role_id int, resource_id int) error {
	query := "SELECT * FROM public.remove_resource_from_role($1, $2);"

	args := []any{role_id, resource_id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetAllResources(ctx context.Context, role_id int) (port.GetAllResourcesResponse, error) {
	var response port.GetAllResourcesResponse
	var responseBase port.GetResourceResponse

	query := "SELECT * FROM public.get_all_resource_by_role($1);"
	dest := []any{&responseBase.Id, &responseBase.Name, &responseBase.Action}
	args := []any{role_id}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResourcesResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResourceResponse{
			Id:     int(res[0].(int64)),
			Name:   res[1].(string),
			Action: res[2].(string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}
