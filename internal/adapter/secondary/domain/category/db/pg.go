package category

import (
	"context"
	"database/sql"
	"log"

	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
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
	// var responseBase port.GetResponse

	query := "SELECT * FROM public.get_all_category();"
	// result := [][]any{
	// 	{
	// 		&responseBase.Id,
	// 		&responseBase.Name,
	// 	}}
	rows, err := p.Pool.QueryContext(ctx, query)
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
		var category port.GetResponse
		if err := rows.Scan(&category.Id, &category.Name); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, category)
	}

	// err := query_handler.NewQuery(
	// 	query_handler.WithCtx(ctx),
	// 	query_handler.WithDB(p.Pool),
	// 	query_handler.WithQuery(query),
	// 	query_handler.WithMultiRowResultSet(nil, result),
	// ).DoStuff()
	// if err != nil {
	// 	return port.GetAllResponse{}, err
	// }
	// for _, res := range result {
	// 	responseBase := port.GetResponse{
	// 		Id:   *res[0].(*int),
	// 		Name: *res[1].(*string),
	// 	}
	// 	response.List = append(response.List, responseBase)
	// }

	return response, nil
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_category($1);"
	args := []any{id}
	result := []any{
		&response.Id,
		&response.Name,
	}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return port.GetResponse{}, err
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_category($1);"
	args := []any{req.Name}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return 0, err
	}

	return resourceId, nil
}

func (p *Postgres) Remove(ctx context.Context, id int) error {
	query := "SELECT * FROM public.remove_category($1);"
	args := []any{id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoStuff()
	if err != nil {
		return err
	}

	return nil
}
