package db

import (
	"context"
	"database/sql"
	"strconv"

	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/distributor_subscription"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (m *Postgres) GetAllPlan(ctx context.Context) (port.GetAllPlanResponse, error) {
	var response port.GetAllPlanResponse
	var responseBase port.GetPlanResponse

	query := "SELECT * FROM public.get_all_subscription_plan();"
	dest := []any{
		&responseBase.Id,
		&responseBase.Name,
		&responseBase.Price,
		&responseBase.TermInMonth,
		responseBase.Description,
	}
	args := []any{}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllPlanResponse{}, err
	}
	for _, res := range result {
		price, _ := strconv.ParseFloat(res[2].(string), 64)
		val := port.GetPlanResponse{
			Id:          int(res[0].(int64)),
			Name:        res[1].(string),
			Price:       price,
			TermInMonth: int(res[3].(int64)),
			Description: res[4].(string),
		}
		response.List = append(response.List, val)
	}
	return response, nil

}

func (m *Postgres) GetPlan(ctx context.Context, id int) (port.GetPlanResponse, error) {
	var response port.GetPlanResponse
	query := "SELECT * FROM public.get_subscription_plan_by_id($1);"
	result := []any{
		&response.Id,
		&response.Name,
		&response.Price,
		&response.TermInMonth,
		&response.Description,
	}
	args := []any{id}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetPlanResponse{}, err
	}
	response.Id = *result[0].(*int)
	response.Name = *result[1].(*string)
	response.Price = *result[2].(*float64)
	response.TermInMonth = *result[3].(*int)
	response.Description = *result[4].(*string)

	return response, nil
}

func (m *Postgres) CreatePlan(ctx context.Context, req *port.CreatePlanRequest) (int, error) {
	var planId int
	query := "SELECT * FROM public.create_subscription_plan($1, $2, $3, $4);"

	result := []any{&planId}
	args := []any{
		&req.Name,
		&req.Price,
		&req.TermInMonth,
		&req.Description}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}
	return planId, nil
}

func (m *Postgres) Place(ctx context.Context, req *port.PlaceRequest) (int, error) {
	var subId int
	query := "SELECT * FROM public.create_distributor_subscription($1, $2, $3, $4);"

	result := []any{&subId}
	args := []any{
		&req.SubscriptionPlanId,
		&req.DistributorId,
		&req.PaymentPartnerId,
		&req.Status,
	}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return 0, err
	}
	return subId, nil
}

func (m *Postgres) GetSubscription(ctx context.Context, subId int) (port.GetSubscriptionResponse, error) {
	var response port.GetSubscriptionResponse
	query := "SELECT * FROM public.get_subscription_by_id($1);"
	result := []any{
		&response.Id,
		&response.SubscriptionPlanId,
		&response.DistributorId,
		&response.Status}
	args := []any{&subId}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetSubscriptionResponse{}, err
	}
	response.Id = *result[0].(*int)
	response.SubscriptionPlanId = *result[1].(*int)
	response.DistributorId = *result[2].(*int)
	response.Status = *result[3].(*string)

	return response, nil
}

func (m *Postgres) GetAllSubscriptions(ctx context.Context) (port.GetAllSubscriptionResponse, error) {
	var response port.GetAllSubscriptionResponse
	var responseBase port.GetSubscriptionResponse

	query := "SELECT * FROM public.get_all_subscriptions();"
	dest := []any{
		&responseBase.Id,
		&responseBase.SubscriptionPlanId,
		&responseBase.DistributorId,
		&responseBase.Status}
	args := []any{}
	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllSubscriptionResponse{}, err
	}
	for _, res := range result {
		val := port.GetSubscriptionResponse{
			Id:                 int(res[0].(int64)),
			SubscriptionPlanId: int(res[1].(int64)),
			DistributorId:      int(res[2].(int64)),
			Status:             res[3].(string),
		}
		response.List = append(response.List, val)
	}
	return response, nil
}

func (m *Postgres) GetSubscriptionByDistributorId(ctx context.Context, distributorId int) (port.GetSubscriptionResponse, error) {
	var response port.GetSubscriptionResponse
	query := "SELECT * FROM public.get_subscription_by_distributor_id($1);"
	result := []any{
		&response.Id,
		&response.SubscriptionPlanId,
		&response.DistributorId,
		&response.Status}
	args := []any{&distributorId}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetSubscriptionResponse{}, err
	}
	response.Id = *result[0].(*int)
	response.SubscriptionPlanId = *result[1].(*int)
	response.DistributorId = *result[2].(*int)
	response.Status = *result[3].(*string)

	return response, nil
}

func (m *Postgres) UpdateStatus(ctx context.Context, req *port.UpdateSubscriptionRequest) error {
	query := "SELECT * FROM public.update_distributor_subscription_status($1, $2);"

	result := []any{}
	args := []any{
		&req.Id,
		&req.Status,
	}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(m.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}
