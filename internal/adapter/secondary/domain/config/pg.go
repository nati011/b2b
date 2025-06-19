package config

import (
	"context"
	"database/sql"

	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/config"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (p *Postgres) GetOrderExpiry(ctx context.Context) (port.GetOrderExpiryResponse, error) {
	var response int
	query := "SELECT * FROM public.get_order_expiry_duration_config();"
	args := []any{}
	result := []any{
		&response,
	}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetOrderExpiryResponse{}, err
	}
	return port.GetOrderExpiryResponse{
		ExpiryDurationInMinues: response,
	}, nil
}

func (p *Postgres) SetOrderExpiryConfig(ctx context.Context, req *port.SetOrderExpiryRequest) error {
	query := "SELECT * FROM public.set_order_expiry_duration_config($1);"
	args := []any{
		req.ExpiryDurationInMinues,
	}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()
	if err != nil {
		return err
	}
	return nil
}
