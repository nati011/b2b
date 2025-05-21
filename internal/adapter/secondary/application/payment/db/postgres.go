package email

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/payment/db"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(db *sql.DB) port.DB {
	return &Postgres{Pool: db}
}

func (p *Postgres) Create(ctx context.Context, req port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.create_payment($1, $2, $3, $4);"

	result := []any{&resourceId}
	args := []any{
		req.OrderId,
		req.PartnerId,
		req.TransactionRef,
		req.Amount}

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

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_payment_by_id($1);"
	result := []any{
		&response.Id,
		&response.OrderId,
		&response.PartnerId,
		&response.TransactionRef,
		&response.Amount,
		&response.Date}
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
	response.OrderId = *result[1].(*int)
	response.PartnerId = *result[2].(*int)
	response.TransactionRef = *result[3].(*string)
	response.Amount = *result[4].(*float64)
	response.Date = *result[5].(*time.Time)

	return response, nil
}

func (p *Postgres) GetByTransactionRef(ctx context.Context, txRef string) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_payment_by_tx_ref($1);"
	result := []any{
		&response.Id,
		&response.OrderId,
		&response.PartnerId,
		&response.TransactionRef,
		&response.Amount,
		&response.Date}
	args := []any{txRef}

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
	response.OrderId = *result[1].(*int)
	response.PartnerId = *result[2].(*int)
	response.TransactionRef = *result[3].(*string)
	response.Amount = *result[4].(*float64)
	response.Date = *result[5].(*time.Time)

	return response, nil
}

func (p *Postgres) GetByOrderId(ctx context.Context, orderId int) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_payment_by_order_id($1);"

	dest := []any{&responseBase.Id,
		&responseBase.OrderId,
		&responseBase.PartnerId,
		&responseBase.TransactionRef,
		&responseBase.Amount,
		&responseBase.Date}
	args := []any{orderId}

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
		v, _ := strconv.ParseFloat(res[4].(string), 64)
		responseBase := port.GetResponse{
			Id:             int(res[0].(int64)),
			OrderId:        responseBase.OrderId,
			PartnerId:      responseBase.PartnerId,
			TransactionRef: responseBase.TransactionRef,
			Amount:         v,
			Date:           responseBase.Date,
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_all_payment();"

	dest := []any{&responseBase.Id,
		&responseBase.OrderId,
		&responseBase.PartnerId,
		&responseBase.TransactionRef,
		&responseBase.Amount,
		&responseBase.Date}
	args := []any{}

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
		v, _ := strconv.ParseFloat(res[4].(string), 64)
		responseBase := port.GetResponse{
			Id:             int(res[0].(int64)),
			OrderId:        responseBase.OrderId,
			PartnerId:      responseBase.PartnerId,
			TransactionRef: responseBase.TransactionRef,
			Amount:         v,
			Date:           responseBase.Date,
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}
