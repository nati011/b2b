package transaction

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/transaction/db"
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
	query := "SELECT * FROM public.get_transaction_by_id($1);"

	result := []any{
		&response.Id,
		&response.Amount,
		&response.PartnerId,
		&response.TxRef,
		&response.Status,
		&response.Date}

	args := []any{id}

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
	response.Amount = *result[1].(*float64)
	response.PartnerId = *result[2].(*int)
	response.TxRef = *result[3].(*string)
	response.Status = *result[4].(*string)
	response.Date = *result[5].(*time.Time)

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_all_transactions($1,$2);"
	dest := []any{
		&responseBase.Id,
		&responseBase.Amount,
		&responseBase.PartnerId,
		&responseBase.TxRef,
		&responseBase.Status,
		&responseBase.Date}
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
		amount, _ := strconv.ParseFloat(res[1].(string), 64)
		responseBase := port.GetResponse{
			Id:        int(res[0].(int64)),
			Amount:    amount,
			PartnerId: int(res[2].(int64)),
			TxRef:     res[3].(string),
			Status:    res[4].(string),
			Date:      res[5].(time.Time),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetByDate(ctx context.Context, date time.Time) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_transactions_by_date($1,$2,$3);"

	dest := []any{
		&responseBase.Id,
		&responseBase.Amount,
		&responseBase.PartnerId,
		&responseBase.TxRef,
		&responseBase.Status,
		&responseBase.Date}
	args := []any{date, p.Pagination.Limit, p.Pagination.Offset}

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
		amount, _ := strconv.ParseFloat(res[1].(string), 64)
		responseBase := port.GetResponse{
			Id:        int(res[0].(int64)),
			Amount:    amount,
			PartnerId: int(res[2].(int64)),
			TxRef:     res[3].(string),
			Status:    res[4].(string),
			Date:      res[5].(time.Time),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetByPartnerId(ctx context.Context, partner_id int) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_transactions_by_partner_id($1,$2,$3);"

	dest := []any{
		&responseBase.Id,
		&responseBase.Amount,
		&responseBase.PartnerId,
		&responseBase.TxRef,
		&responseBase.Status,
		&responseBase.Date}
	args := []any{partner_id, p.Pagination.Limit, p.Pagination.Offset}

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
		amount, _ := strconv.ParseFloat(res[1].(string), 64)
		responseBase := port.GetResponse{
			Id:        int(res[0].(int64)),
			Amount:    amount,
			PartnerId: int(res[2].(int64)),
			TxRef:     res[3].(string),
			Status:    res[4].(string),
			Date:      res[5].(time.Time),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetByTxRef(ctx context.Context, tx_ref string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_transactions_by_tx_ref($1,$2,$3);"

	dest := []any{
		&responseBase.Id,
		&responseBase.Amount,
		&responseBase.PartnerId,
		&responseBase.TxRef,
		&responseBase.Status,
		&responseBase.Date}

	log.Printf("Args, %v", tx_ref)
	args := []any{tx_ref, p.Pagination.Limit, p.Pagination.Offset}

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
		amount, _ := strconv.ParseFloat(res[1].(string), 64)
		responseBase := port.GetResponse{
			Id:        int(res[0].(int64)),
			Amount:    amount,
			PartnerId: int(res[2].(int64)),
			TxRef:     res[3].(string),
			Status:    res[4].(string),
			Date:      res[5].(time.Time),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_transactions_by_status($1,$2,$3);"

	dest := []any{
		&responseBase.Id,
		&responseBase.Amount,
		&responseBase.PartnerId,
		&responseBase.TxRef,
		&responseBase.Status,
		&responseBase.Date}
	args := []any{status, p.Pagination.Limit, p.Pagination.Offset}

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
		amount, _ := strconv.ParseFloat(res[1].(string), 64)
		responseBase := port.GetResponse{
			Id:        int(res[0].(int64)),
			Amount:    amount,
			PartnerId: int(res[2].(int64)),
			TxRef:     res[3].(string),
			Status:    res[4].(string),
			Date:      res[5].(time.Time),
		}
		response.List = append(response.List, responseBase)
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var resourceId int
	query := "SELECT * FROM public.record_transaction($1, $2, $3, $4);"

	result := []any{&resourceId}
	args := []any{req.Amount, req.PartnerId, req.Status, req.TxRef}

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

func (p *Postgres) UpdateByTransactionRef(ctx context.Context, req *port.UpdateByTransactionRefRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_transaction_by_transaction_ref($1, $2);"

	result := []any{&resourceId}
	args := []any{req.TransactionRef, req.Status}

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

func (p *Postgres) UpdateStatus(ctx context.Context, req *port.UpdateRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_transaction_status($1, $2);"

	result := []any{&resourceId}
	args := []any{req.Id, req.Status}

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
