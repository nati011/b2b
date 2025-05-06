package email

import (
	"context"
	"database/sql"

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
	query := "SELECT * FROM public.create_payment($1, $2, $3);"

	result := []any{&resourceId}
	args := []any{req.OrderId, req.PartnerId, req.TransactionRef}

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

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetByTransactionRef(ctx context.Context, name string) (port.GetResponse, error) {
	return port.GetResponse{}, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	return port.GetAllResponse{}, nil
}
