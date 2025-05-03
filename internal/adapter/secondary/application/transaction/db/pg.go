package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"b2b.nati011.github.com/config"
	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/transaction/db"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
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
	var amountStr string
	query := "SELECT * FROM public.get_transaction_by_id($1);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		id,
	)

	if err != nil {
		return port.GetResponse{}, err
	}

	rows.Row.Scan(
		&response.Id,
		&amountStr,
		&response.PartnerId,
		&response.TxRef,
		&response.Status,
		&response.Date,
	)

	// Convert amountStr to int64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return port.GetResponse{}, fmt.Errorf("failed to parse amount: %w", err)
	}
	response.Amount = amount
	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_transactions($1,$2);"

	limit := p.Pagination.Limit
	offset := p.Pagination.Offset

	rows, err := handler.MustQueryRow(p.Pool, ctx, query, true, limit, offset)

	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var transaction port.GetResponse
		var amountStr string
		if err := rows.Rows.Scan(
			&transaction.Id,
			&amountStr,
			&transaction.PartnerId,
			&transaction.TxRef,
			&transaction.Status,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		// Convert amountStr to int64
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return port.GetAllResponse{}, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = amount
		response.List = append(response.List, transaction)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port_commons.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByDate(ctx context.Context, date time.Time) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_transactions_by_date($1,$2,$3);"

	limit := p.Pagination.Limit
	offset := p.Pagination.Offset

	rows, err := handler.MustQueryRow(p.Pool, ctx, query, true, date, limit, offset)

	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var transaction port.GetResponse
		var amountStr string
		if err := rows.Rows.Scan(
			&transaction.Id,
			&amountStr,
			&transaction.PartnerId,
			&transaction.TxRef,
			&transaction.Status,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		// Convert amountStr to int64
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return port.GetAllResponse{}, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = amount
		response.List = append(response.List, transaction)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port_commons.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByPartnerId(ctx context.Context, partner_id int) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_transactions_by_partner_id($1,$2,$3);"

	limit := p.Pagination.Limit
	offset := p.Pagination.Offset
	rows, err := handler.MustQueryRow(p.Pool, ctx, query, true, partner_id, limit, offset)

	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var transaction port.GetResponse
		var amountStr string
		if err := rows.Rows.Scan(
			&transaction.Id,
			&amountStr,
			&transaction.PartnerId,
			&transaction.TxRef,
			&transaction.Status,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		// Convert amountStr to int64
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return port.GetAllResponse{}, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = amount
		response.List = append(response.List, transaction)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port_commons.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByTxRef(ctx context.Context, tx_ref string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_transactions_by_tx_ref($1,$2,$3);"

	limit := p.Pagination.Limit
	offset := p.Pagination.Offset
	rows, err := handler.MustQueryRow(p.Pool, ctx, query, true, tx_ref, limit, offset)

	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var transaction port.GetResponse
		var amountStr string
		if err := rows.Rows.Scan(
			&transaction.Id,
			&amountStr,
			&transaction.PartnerId,
			&transaction.TxRef,
			&transaction.Status,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		// Convert amountStr to int64
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return port.GetAllResponse{}, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = amount
		response.List = append(response.List, transaction)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port_commons.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_transactions_by_status($1,$2,$3);"

	limit := p.Pagination.Limit
	offset := p.Pagination.Offset
	rows, err := handler.MustQueryRow(p.Pool, ctx, query, true, status, limit, offset)

	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var transaction port.GetResponse
		var amountStr string
		if err := rows.Rows.Scan(
			&transaction.Id,
			&amountStr,
			&transaction.PartnerId,
			&transaction.TxRef,
			&transaction.Status,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		// Convert amountStr to int64
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return port.GetAllResponse{}, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = amount
		response.List = append(response.List, transaction)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port_commons.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var id int
	query := "SELECT * FROM public.record_transaction($1, $2, $3, $4);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Amount,
		req.PartnerId,
		req.Status,
		req.TxRef,
	)
	if err != nil {
		return 0, err
	}
	rows.Row.Scan(&id)
	return id, nil
}

func (p *Postgres) UpdateStatus(ctx context.Context, req *port.UpdateRequest) error {
	var id int
	query := "SELECT * FROM public.update_transaction_status($1, $2);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Id,
		req.Status,
	)
	if err != nil {
		return err
	}
	rows.Row.Scan(&id)
	return nil
}
