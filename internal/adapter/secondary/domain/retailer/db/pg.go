package adapter

import (
	"context"
	"database/sql"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/retailer"
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

func (r *Postgres) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	var retailerId int
	query := "SELECT * FROM public.create_retailer($1, $2, $3, $4, $5, $6, $7);"

	result := []any{&retailerId}
	args := []any{
		req.Name,
		req.Tin,
		req.Latitude,
		req.Longitude,
		req.GeneralZone,
		req.Region,
		req.Woreda}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	err = r.CreateRetailerUser(ctx, &port.CreateUserAgentRequest{
		User_id:     req.UserId,
		Retailer_id: retailerId,
	})
	if err != nil {
		switch err {
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return retailerId, nil
}

func (r *Postgres) CreateRetailerUser(ctx context.Context, req *port.CreateUserAgentRequest) error {
	query := "SELECT * FROM public.create_retailer_user($1, $2);"
	var distributorUserId int

	result := []any{&distributorUserId}
	args := []any{req.Retailer_id, req.User_id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	return err
}

func (r *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_retailer_by_id($1);"
	result := []any{&response.Id,
		&response.Name,
		&response.Tin,
		&response.Latitude,
		&response.Longitude,
		&response.GeneralZone,
		&response.Region,
		&response.Woreda}

	args := []any{id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetResponse{}, err
	}

	response.Id = *result[0].(*int)
	response.Name = *result[1].(*string)
	response.Tin = *result[2].(*string)
	response.Latitude = *result[3].(*string)
	response.Longitude = *result[4].(*string)
	response.GeneralZone = *result[5].(*string)
	response.Region = *result[6].(*string)
	response.Woreda = *result[7].(*string)

	return response, nil
}

func (r *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	query := "SELECT * FROM public.update_retailer_name($1, $2);"
	args := []any{&req.Id, &req.Name}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return err
}

func (r *Postgres) UpdateTin(ctx context.Context, req *port.UpdateTinRequest) error {
	query := "SELECT * FROM public.update_retailer_tin($1, $2);"
	args := []any{&req.Id, &req.Tin}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return err
}

func (r *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_all_retailers($1, $2);"

	dest := []any{&responseBase.Id,
		&responseBase.Name,
		&responseBase.Tin,
		&responseBase.Latitude,
		&responseBase.Longitude,
		&responseBase.GeneralZone,
		&responseBase.Region,
		&responseBase.Woreda}
	args := []any{r.Pagination.Limit, r.Pagination.Offset}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:          *res[0].(*int),
			Name:        *res[1].(*string),
			Tin:         *res[2].(*string),
			Latitude:    *res[3].(*string),
			Longitude:   *res[4].(*string),
			GeneralZone: *res[5].(*string),
			Region:      *res[6].(*string),
			Woreda:      *res[7].(*string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (r *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_retailer_by_name($1);"
	dest := []any{&responseBase.Id,
		&responseBase.Name,
		&responseBase.Tin,
		&responseBase.Latitude,
		&responseBase.Longitude,
		&responseBase.GeneralZone,
		&responseBase.Region,
		&responseBase.Woreda}
	args := []any{name}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:          *res[0].(*int),
			Name:        *res[1].(*string),
			Tin:         *res[2].(*string),
			Latitude:    *res[3].(*string),
			Longitude:   *res[4].(*string),
			GeneralZone: *res[5].(*string),
			Region:      *res[6].(*string),
			Woreda:      *res[7].(*string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (r *Postgres) GetByTin(ctx context.Context, tin string) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_retailer_by_tin($1);"
	result := []any{&response.Id,
		&response.Name,
		&response.Tin,
		&response.Latitude,
		&response.Longitude,
		&response.GeneralZone,
		&response.Region,
		&response.Woreda}

	args := []any{tin}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetResponse{}, err
	}
	response.Id = *result[0].(*int)
	response.Name = *result[1].(*string)
	response.Tin = *result[2].(*string)
	response.Latitude = *result[3].(*string)
	response.Longitude = *result[4].(*string)
	response.GeneralZone = *result[5].(*string)
	response.Region = *result[6].(*string)
	response.Woreda = *result[7].(*string)

	return response, nil
}

func (r *Postgres) GetAllUserAgents(ctx context.Context, id int) (port.GetAllUserResponse, error) {
	var response port.GetAllUserResponse
	var responseBase port.GetUserResponse

	query := "SELECT * FROM public.get_all_retailer_users($1);"
	dest := []any{&responseBase.Id}

	args := []any{id}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllUserResponse{}, err
	}

	for _, res := range result {
		responseBase := port.GetUserResponse{
			Id: *res[0].(*int),
		}
		response.List = append(response.List, responseBase)
	}
	return response, nil
}
