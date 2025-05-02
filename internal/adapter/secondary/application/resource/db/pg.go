package adapter

import (
	"context"
	"database/sql"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/resource"
)

type Postgres struct {
	Pool       *sql.DB
	Pagination *config.Pagination
}

func NewPostgres(db *sql.DB, pagination *config.Pagination) port.DB {
	return &Postgres{
		Pool:       db,
		Pagination: pagination,
	}
}

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_resources_by_id($1);"

	result := []any{&response.Id, &response.Action, &response.Name}
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
	response.Action = *result[1].(*string)
	response.Name = *result[2].(*string)

	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_resources_by_name($1);"

	result := []any{&response.Id, &response.Action, &response.Name}
	args := []any{&name}

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
	response.Action = *result[1].(*string)
	response.Name = *result[2].(*string)

	return response, nil
}
func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_all_resources($1, $2);"

	dest := []any{&responseBase.Id, &responseBase.Name, &responseBase.Action}
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

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:     *res[0].(*int),
			Action: *res[1].(*string),
			Name:   *res[2].(*string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_resource($1, $2);"

	result := []any{&resourceId}
	args := []any{req.Name, req.Action}

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

func (p *Postgres) UpdateAction(ctx context.Context, req *port.UpdateActionRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_resource_action($1, $2);"

	result := []any{&resourceId}
	args := []any{req.Id, req.Action}

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
	query := "SELECT * FROM public.update_resource_name($1, $2);"

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
	query := "SELECT * FROM public.delete_resource($1);"

	// result := []any{nil}
	args := []any{&id}

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
