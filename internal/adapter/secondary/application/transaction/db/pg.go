package transaction

import (
	"context"
	"database/sql"
	"log"
	"time"

	port "b2b.nati011.github.com/internal/port/application/transaction/db"
)

type Postgres struct {
	Pool *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		Pool: DB,
	}
}

func (p *Postgres) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_transaction_by_id($1);"
	err := p.Pool.QueryRowContext(ctx, query, id).Scan(
		&response.Id,
		&response.User_Id,
		&response.Amount,
		&response.Partner_Id,
		&response.Date)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}
	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_transactions();"
	rows, err := p.Pool.QueryContext(ctx, query)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var transaction port.GetResponse
		if err := rows.Scan(
			&transaction.Id,
			&transaction.User_Id,
			&transaction.Amount,
			&transaction.Partner_Id,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, transaction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByDate(ctx context.Context, date time.Time) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_transactions_by_date($1);"
	rows, err := p.Pool.QueryContext(ctx, query, date)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var transaction port.GetResponse
		if err := rows.Scan(
			&transaction.Id,
			&transaction.User_Id,
			&transaction.Amount,
			&transaction.Partner_Id,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, transaction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByUserId(ctx context.Context, user_id int) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_transactions_by_user_id($1);"
	rows, err := p.Pool.QueryContext(ctx, query, user_id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var transaction port.GetResponse
		if err := rows.Scan(
			&transaction.Id,
			&transaction.User_Id,
			&transaction.Amount,
			&transaction.Partner_Id,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, transaction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByPartnerId(ctx context.Context, partner_id int) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_transactions_by_partner_id($1);"
	rows, err := p.Pool.QueryContext(ctx, query, partner_id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetAllResponse{}, port.ErrSysNoRows
		default:
			return port.GetAllResponse{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var transaction port.GetResponse
		if err := rows.Scan(
			&transaction.Id,
			&transaction.User_Id,
			&transaction.Amount,
			&transaction.Partner_Id,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, transaction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var id int
	query := "SELECT * FROM public.record_transaction($1, $2, $3);"
	err := p.Pool.QueryRowContext(ctx, query,
		req.User_Id,
		req.Partner_Id,
		req.Amount).Scan(
		&id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}
	return id, nil
}
