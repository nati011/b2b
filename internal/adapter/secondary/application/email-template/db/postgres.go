package email

import (
	"context"
	"database/sql"

	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/email-template"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{Pool: db}
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_email_template($1, $2);"
	result := []any{&resourceId}
	args := []any{req.Name, req.HtmlTemplate}
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

func (p *Postgres) Update(ctx context.Context, req *port.UpdateRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_template($1, $2);"
	result := []any{&resourceId}
	args := []any{req.Id, req.Name, req.HtmlTemplate}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_template_by_id($1);"

	result := []any{&response.Id, &response.Name, &response.HtmlTemplate}
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
	response.HtmlTemplate = *result[2].(*string)

	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_template_by_name($1);"

	result := []any{&response.Id, &response.Name, &response.HtmlTemplate}
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
	response.Name = *result[1].(*string)
	response.HtmlTemplate = *result[2].(*string)

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_all_templates();"

	dest := []any{&responseBase.Id, &responseBase.Name, &responseBase.HtmlTemplate}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(nil, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:           int(res[0].(int64)),
			Name:         res[1].(string),
			HtmlTemplate: res[2].(string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}
