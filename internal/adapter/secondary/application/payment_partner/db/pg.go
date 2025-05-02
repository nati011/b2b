package db

import (
	"context"
	"database/sql"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/partner/db"
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

	query := "SELECT * FROM public.get_payment_partner_by_id($1);"

	result := []any{
		&response.Id,
		&response.Name,
		&response.Icon,
		&response.Status,
		&response.BaseURL}

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
	response.Icon = *result[2].(*string)
	response.Status = *result[3].(*string)
	response.BaseURL = *result[4].(*string)

	return response, nil
}
func (p *Postgres) GetPartnerSecret(ctx context.Context, id int) (port.GetPartnerSecret, error) {
	var response port.GetPartnerSecret

	query := "SELECT * FROM public.get_payment_partner_secrets($1);"

	result := []any{
		&response.Name,
		&response.BaseURL,
		&response.Secret,
	}

	args := []any{&id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetPartnerSecret{}, err
	}

	response.Name = *result[0].(*string)
	response.BaseURL = *result[1].(*string)
	response.Secret = *result[2].(*string)

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_all_payment_partners($1,$2);"

	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Icon,
		&responseBase.Status,
		&responseBase.BaseURL,
	}
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
			Id:      *res[0].(*int),
			Name:    *res[1].(*string),
			Icon:    *res[2].(*string),
			Status:  *res[3].(*string),
			BaseURL: *res[4].(*string),
		}
		response.List = append(response.List, responseBase)
	}
	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_payment_partner_by_status($1);"

	dest := []any{&responseBase.Id,
		&responseBase.Name,
		&responseBase.Icon,
		&responseBase.Status,
		&responseBase.BaseURL,
	}
	args := []any{&status}

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
			Id:      *res[0].(*int),
			Name:    *res[1].(*string),
			Icon:    *res[2].(*string),
			Status:  *res[3].(*string),
			BaseURL: *res[4].(*string),
		}
		response.List = append(response.List, responseBase)
	}
	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_payment_partner_by_name($1);"

	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Icon,
		&responseBase.Status,
		&responseBase.BaseURL,
	}
	args := []any{name}

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
			Id:      *res[0].(*int),
			Name:    *res[1].(*string),
			Icon:    *res[2].(*string),
			Status:  *res[3].(*string),
			BaseURL: *res[4].(*string),
		}
		response.List = append(response.List, responseBase)
	}
	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.create_payment_partner($1, $2, $3, $4, $5);"

	result := []any{&partner_id}
	args := []any{
		req.Name,
		req.Icon,
		req.Status,
		req.BaseURL,
		req.Secret}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}
	return partner_id, nil
}

func (p *Postgres) UpdateStatus(ctx context.Context, id int, status string) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.update_payment_partner_status($1, $2);"

	result := []any{&partner_id}
	args := []any{id, status}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}
	return partner_id, nil
}

func (p *Postgres) UpdateName(ctx context.Context, id int, name string) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.update_payment_partner_name($1, $2);"

	result := []any{&partner_id}
	args := []any{id, name}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	return partner_id, nil
}

func (p *Postgres) Delete(context.Context, int) error {
	return nil
}
