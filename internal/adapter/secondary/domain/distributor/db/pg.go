package db

import (
	"context"
	"database/sql"
	"log"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/distributor"
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
	var distributorId int
	query := "SELECT * FROM public.create_distributor($1, $2, $3, $4, $5, $6, $7, $8);"

	result := []any{&distributorId}
	args := []any{
		req.Name,
		req.Tin,
		req.Latitude,
		req.Longitude,
		req.GeneralZone,
		req.Region,
		req.Woreda,
		req.LicenceURL}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	_, err = r.CreateDistributorUser(ctx, &port.CreateUserAgentRequest{
		User_id:        req.UserId,
		Distributor_Id: distributorId,
	})
	if err != nil {
		switch err {
		default:
			return 0, port_commons.ErrSysUnknown
		}
	}

	return distributorId, nil
}

func (r *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_distributor_by_id($1);"
	result := []any{&response.Id,
		&response.Name,
		&response.Tin,
		&response.LicenceURL,
		&response.Latitude,
		&response.Longitude,
		&response.GeneralZone,
		&response.Region,
		&response.Woreda,
		&response.IsActive,
		&response.Verdict,
	}

	args := []any{&id}

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
	response.LicenceURL = *result[3].(*string)
	response.Latitude = *result[4].(*string)
	response.Longitude = *result[5].(*string)
	response.GeneralZone = *result[6].(*string)
	response.Region = *result[7].(*string)
	response.Woreda = *result[8].(*string)
	response.IsActive = *result[9].(*bool)
	response.Verdict = *result[10].(*string)

	return response, nil
}

func (r *Postgres) UpdateName(ctx context.Context, req *port.UpdateNameRequest) error {
	query := "SELECT * FROM public.update_distributor_name($1, $2);"
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

	return nil
}

func (r *Postgres) UpdateTin(ctx context.Context, req *port.UpdateTinRequest) error {
	query := "SELECT * FROM public.update_distributor_tin($1, $2);"
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
	return nil
}

func (r *Postgres) Activate(ctx context.Context, id int) error {
	query := "SELECT * FROM public.activate_distributors($1);"
	args := []any{&id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (r *Postgres) Dectivate(ctx context.Context, id int) error {
	query := "SELECT * FROM public.deactivate_distributors($1);"
	args := []any{&id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (r *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	query := "SELECT * FROM public.get_all_distributors($1, $2);"

	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Tin,
		&responseBase.Latitude,
		&responseBase.Longitude,
		&responseBase.GeneralZone,
		&responseBase.Region,
		&responseBase.Woreda,
		&responseBase.IsActive,
		&responseBase.Verdict,
		&response.TotalCount,
	}
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
			Id:          int(res[0].(int64)),
			Name:        res[1].(string),
			Tin:         res[2].(string),
			Latitude:    res[3].(string),
			Longitude:   res[4].(string),
			GeneralZone: res[5].(string),
			Region:      res[6].(string),
			Woreda:      res[7].(string),
			IsActive:    res[8].(bool),
			Verdict:     res[9].(string),
		}
		response.List = append(response.List, responseBase)
		response.TotalCount = res[10].(int64)
	}

	return response, nil
}

func (r *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_distributor_by_name($1,$2,$3);"

	dest := []any{&responseBase.Id,
		&responseBase.Name,
		&responseBase.Tin,
		&responseBase.Latitude,
		&responseBase.Longitude,
		&responseBase.GeneralZone,
		&responseBase.Region,
		&responseBase.Woreda,
		&responseBase.IsActive}
	args := []any{name, r.Pagination.Limit, r.Pagination.Offset}

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
			Id:          int(res[0].(int64)),
			Name:        res[1].(string),
			Tin:         res[2].(string),
			Latitude:    res[3].(string),
			Longitude:   res[4].(string),
			GeneralZone: res[5].(string),
			Region:      res[6].(string),
			Woreda:      res[7].(string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (r *Postgres) GetByStatus(ctx context.Context, status bool) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_distributor_by_status($1,$2,$3);"

	dest := []any{&responseBase.Id,
		&responseBase.Name,
		&responseBase.Tin,
		&responseBase.Latitude,
		&responseBase.Longitude,
		&responseBase.GeneralZone,
		&responseBase.Region,
		&responseBase.Woreda,
		&responseBase.IsActive,
		&responseBase.Verdict,
		&response.TotalCount,
	}
	args := []any{status, r.Pagination.Limit, r.Pagination.Offset}

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
			Id:          int(res[0].(int64)),
			Name:        res[1].(string),
			Tin:         res[2].(string),
			Latitude:    res[3].(string),
			Longitude:   res[4].(string),
			GeneralZone: res[5].(string),
			Region:      res[6].(string),
			Woreda:      res[7].(string),
			IsActive:    res[8].(bool),
			Verdict:     res[9].(string),
		}
		response.List = append(response.List, responseBase)
		response.TotalCount = res[10].(int64)
	}

	return response, nil
}

func (r *Postgres) GetByApprovalStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_distributor_by_approval_status($1,$2,$3);"

	dest := []any{&responseBase.Id,
		&responseBase.Name,
		&responseBase.Tin,
		&responseBase.Latitude,
		&responseBase.Longitude,
		&responseBase.GeneralZone,
		&responseBase.Region,
		&responseBase.Woreda,
		&responseBase.IsActive,
		&responseBase.Verdict,
		&response.TotalCount,
	}
	args := []any{status, r.Pagination.Limit, r.Pagination.Offset}

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
			Id:          int(res[0].(int64)),
			Name:        res[1].(string),
			Tin:         res[2].(string),
			Latitude:    res[3].(string),
			Longitude:   res[4].(string),
			GeneralZone: res[5].(string),
			Region:      res[6].(string),
			Woreda:      res[7].(string),
			IsActive:    res[8].(bool),
			Verdict:     res[9].(string),
		}
		response.List = append(response.List, responseBase)
		response.TotalCount = res[10].(int64)
	}

	return response, nil
}

func (r *Postgres) GetByTin(ctx context.Context, tin string) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_distributor_by_tin($1);"

	result := []any{&response.Id,
		&response.Name,
		&response.Tin,
		&response.Latitude,
		&response.Longitude,
		&response.GeneralZone,
		&response.Region,
		&response.Woreda,
		&response.IsActive}

	args := []any{&tin}

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

func (r *Postgres) GetByUserId(ctx context.Context, user_id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_distributor_by_user_id($1);"

	result := []any{&response.Id,
		&response.Name,
		&response.Tin,
		&response.Latitude,
		&response.Longitude,
		&response.GeneralZone,
		&response.Region,
		&response.Woreda,
		&response.IsActive,
		&response.Verdict,
	}

	args := []any{&user_id}

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

	query := "SELECT * FROM public.get_all_distributor_users($1, $2, $3);"
	dest := []any{&responseBase.Id}
	args := []any{id, r.Pagination.Limit, r.Pagination.Offset}

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
			Id: int(res[0].(int64)),
		}
		response.List = append(response.List, responseBase)
	}
	return response, nil
}

func (m *Postgres) GetAllUserDetail(ctx context.Context, distributor_id int) (port.GetAllUserDetailResponse, error) {
	var response port.GetAllUserDetailResponse
	var responseBase port.GetUserDetailResponse

	query := "SELECT * FROM public.get_all_distributor_user_details($1, $2, $3);"

	dest := []any{
		&responseBase.Id,
		&responseBase.FirstName,
		&responseBase.LastName,
		&responseBase.Email,
		&responseBase.Phone,
		&responseBase.Username,
		&response.TotalCount,
	}
	args := []any{distributor_id, m.Pagination.Limit, m.Pagination.Offset}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()

	if err != nil {
		return port.GetAllUserDetailResponse{}, err
	}
	log.Print(result[0]...)
	for _, res := range result {
		responseBase := port.GetUserDetailResponse{
			Id:        int(res[0].(int64)),
			FirstName: res[1].(string),
			LastName:  res[2].(string),
			Email:     res[3].(string),
			Phone:     res[4].(string),
			Username:  res[5].(string),
		}
		response.List = append(response.List, responseBase)
		response.TotalCount = int(res[0].(int64))
	}
	return response, nil
}

func (r *Postgres) CreateDistributorUser(ctx context.Context, req *port.CreateUserAgentRequest) (int, error) {
	query := "SELECT * FROM public.create_distributor_user($1, $2);"
	var distributorUserId int

	result := []any{&distributorUserId}
	args := []any{req.Distributor_Id, req.User_id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(r.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	return distributorUserId, nil
}
