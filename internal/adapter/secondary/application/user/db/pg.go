package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/user"
)

type Postgres struct {
	db         *sql.DB
	Pagination *config.Pagination
}

func NewPostgres(DB *sql.DB, pagination *config.Pagination) port.DB {
	return &Postgres{
		db:         DB,
		Pagination: pagination,
	}
}

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_users_by_id($1);"
	result := []any{&response.Id,
		&response.FirstName,
		&response.LastName,
		&response.Email,
		&response.Phone,
		&response.Username,
		&response.DOB,
		&response.IsActive,
		&response.ExternalId}
	args := []any{&id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetResponse{}, err
	}
	response.Id = *result[0].(*int)
	response.FirstName = *result[1].(*string)
	response.LastName = *result[2].(*string)
	response.Email = *result[3].(*string)
	response.Phone = *result[4].(*string)
	response.Username = *result[5].(*string)
	response.DOB = *result[6].(*time.Time)
	response.IsActive = *result[7].(*bool)
	response.ExternalId = *result[8].(*string)

	return response, nil
}

func (p *Postgres) GetByEmail(ctx context.Context, email string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	dest := []any{&responseBase.Id, &responseBase.FirstName, &responseBase.LastName, &responseBase.Email, &responseBase.Phone, &responseBase.Username, &responseBase.DOB, &responseBase.IsActive, &responseBase.ExternalId}
	args := []any{email}

	query := "SELECT * FROM public.get_users_by_email($1);"

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:         int(res[0].(int64)),
			FirstName:  res[1].(string),
			LastName:   res[2].(string),
			Email:      res[3].(string),
			Phone:      res[4].(string),
			Username:   res[5].(string),
			DOB:        res[6].(time.Time),
			IsActive:   res[7].(bool),
			ExternalId: res[8].(string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetByPhone(ctx context.Context, phone string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	dest := []any{&responseBase.Id, &responseBase.FirstName, &responseBase.LastName, &responseBase.Email, &responseBase.Phone, &responseBase.Username, &responseBase.DOB, &responseBase.IsActive, &responseBase.ExternalId}
	args := []any{phone}

	query := "SELECT * FROM public.get_users_by_phone($1);"

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:         int(res[0].(int64)),
			FirstName:  res[1].(string),
			LastName:   res[2].(string),
			Email:      res[3].(string),
			Phone:      res[4].(string),
			Username:   res[5].(string),
			DOB:        res[6].(time.Time),
			IsActive:   res[7].(bool),
			ExternalId: res[8].(string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetByUsername(ctx context.Context, username string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	dest := []any{&responseBase.Id, &responseBase.FirstName, &responseBase.LastName, &responseBase.Email, &responseBase.Phone, &responseBase.Username, &responseBase.DOB, &responseBase.IsActive, &responseBase.ExternalId}
	args := []any{username}

	query := "SELECT * FROM public.get_users_by_username($1);"

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:         int(res[0].(int64)),
			FirstName:  res[1].(string),
			LastName:   res[2].(string),
			Email:      res[3].(string),
			Phone:      res[4].(string),
			Username:   res[5].(string),
			DOB:        res[6].(time.Time),
			IsActive:   res[7].(bool),
			ExternalId: res[8].(string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetByActiveStatus(ctx context.Context, status bool) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_users_by_active_status($1,$2,$3);"

	dest := []any{&responseBase.Id, &responseBase.FirstName, &responseBase.LastName, &responseBase.Email, &responseBase.Phone, &responseBase.Username, &responseBase.DOB, &responseBase.IsActive, &responseBase.ExternalId}
	args := []any{status, p.Pagination.Limit, p.Pagination.Offset}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:         int(res[0].(int64)),
			FirstName:  res[1].(string),
			LastName:   res[2].(string),
			Email:      res[3].(string),
			Phone:      res[4].(string),
			Username:   res[5].(string),
			DOB:        res[6].(time.Time),
			IsActive:   res[7].(bool),
			ExternalId: res[8].(string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_users_by_external_id($1);"

	dest := []any{&responseBase.Id, &responseBase.FirstName, &responseBase.LastName, &responseBase.Email, &responseBase.Phone, &responseBase.Username, &responseBase.DOB, &responseBase.IsActive, &responseBase.ExternalId}
	args := []any{extId}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:         int(res[0].(int64)),
			FirstName:  res[1].(string),
			LastName:   res[2].(string),
			Email:      res[3].(string),
			Phone:      res[4].(string),
			Username:   res[5].(string),
			DOB:        res[6].(time.Time),
			IsActive:   res[7].(bool),
			ExternalId: res[8].(string),
		}
		response.List = append(response.List, responseBase)
	}
	return response, nil
}

func (p *Postgres) GetUserProvider(ctx context.Context, id int) (port.GetUserProviderResponse, error) {
	var response port.GetUserProviderResponse
	var responseBase port.UserProvider

	query := "SELECT * FROM public.get_user_provider($1);"
	dest := []any{&responseBase.UserId, &responseBase.ProviderId}
	args := []any{id}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetUserProviderResponse{}, err
	}
	for _, res := range result {
		response.List = append(response.List, port.UserProvider{
			UserId:     *res[0].(*int),
			ProviderId: *res[1].(*string),
		})
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_all_users($1,$2);"

	dest := []any{&responseBase.Id, &responseBase.FirstName, &responseBase.LastName, &responseBase.Email, &responseBase.Phone, &responseBase.Username, &responseBase.DOB, &responseBase.IsActive, &responseBase.ExternalId}
	args := []any{p.Pagination.Limit, p.Pagination.Offset}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	// convert
	for _, res := range result {
		responseBase := port.GetResponse{
			Id:         int(res[0].(int64)),
			FirstName:  res[1].(string),
			LastName:   res[2].(string),
			Email:      res[3].(string),
			Phone:      res[4].(string),
			Username:   res[5].(string),
			DOB:        res[6].(time.Time),
			IsActive:   res[7].(bool),
			ExternalId: res[8].(string),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	//activate by default
	var resourceId int
	query := "SELECT * FROM public.create_user($1, $2, $3, $4, $5, $6, $7);"
	args := []any{p.Pagination.Limit, p.Pagination.Offset}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	return resourceId, nil
}

func (p *Postgres) CreateAndActivate(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_and_activate_user($1, $2, $3, $4, $5, $6, $7);"

	args := []any{req.FirstName, req.LastName, req.Email, req.Phone, req.Username, req.DOB, req.ExternalId}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}
	return resourceId, nil
}

func (p *Postgres) CreateUserProvider(ctx context.Context, req *port.CreateUserProviderRequest) error {
	var resourceId any
	query := "SELECT * FROM public.create_user_provider($1, $2);"

	args := []any{req.UserId, req.ProviderId}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	log.Print(&resourceId)

	return nil
}

func (p *Postgres) Delete(ctx context.Context, id int) error {
	query := "SELECT * FROM public.delete_user($1);"
	args := []any{id}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateFirstName(ctx context.Context, req *port.UpdateFirstNameRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_FirstName($1, $2);"

	args := []any{req.Id, req.FirstName}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	return resourceId, nil
}

func (p *Postgres) UpdateEmail(ctx context.Context, req *port.UpdateEmailRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_email($1, $2);"

	args := []any{req.Id, req.Email}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	return resourceId, nil
}

func (p *Postgres) UpdateDOB(ctx context.Context, req *port.UpdateDOBRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_dob($1, $2);"

	args := []any{req.Id, req.DOB}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	return resourceId, nil
}

func (p *Postgres) UpdateIsActiveStatus(ctx context.Context, req *port.UpdateIsActiveRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_is_active_status($1, $2);"

	args := []any{req.Id, req.IsActive}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	return resourceId, nil
}

func (p *Postgres) UpdatePhone(ctx context.Context, req *port.UpdatePhoneRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_phone($1, $2);"

	args := []any{req.Id, req.Phone}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}

	return resourceId, nil
}

func (p *Postgres) UpdateUsername(ctx context.Context, req *port.UpdateUsernameRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.update_user_name($1, $2);"

	args := []any{req.Id, req.Username}
	result := []any{&resourceId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}
	return resourceId, nil
}

func (p *Postgres) AssignRole(ctx context.Context, id int, role_id int) error {
	query := "SELECT * FROM public.add_role_to_user($1, $2);"
	args := []any{id, role_id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) RemoveAssignedRole(ctx context.Context, id int, role_id int) error {
	query := "SELECT * FROM public.remove_role_from_user($1, $2);"

	args := []any{id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetAllAssignedRole(ctx context.Context, id int) (port.GetAllAssignedRoleResponse, error) {
	var response port.GetAllAssignedRoleResponse
	var responseBase port.GetAssignedRoleResponse

	query := "SELECT * FROM public.get_all_role_by_user($1);"
	dest := []any{&responseBase.Id}

	args := []any{id}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()

	if err != nil {
		return port.GetAllAssignedRoleResponse{}, err
	}

	for _, res := range result {
		response.List = append(response.List, port.GetAssignedRoleResponse{
			Id: *res[0].(*int),
		})
	}
	return response, nil
}
