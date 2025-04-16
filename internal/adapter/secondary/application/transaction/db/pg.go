package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
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
	var amountStr string
	query := "SELECT * FROM public.get_transaction_by_id($1);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		id,
	)

	if err != nil {
		return port.GetResponse{}, err
	}

	rows.Scan(
		&response.Id,
		&response.User_Id,
		&amountStr,
		&response.Partner_Id,
		&response.Date)

	// Convert amountStr to int64
	amountFloat, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return port.GetResponse{}, fmt.Errorf("failed to parse amount: %w", err)
	}
	response.Amount = int64(amountFloat)
	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_transactions();"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
	)

	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction port.GetResponse
		var amountStr string
		if err := rows.Scan(
			&transaction.Id,
			&transaction.User_Id,
			&amountStr,
			&transaction.Partner_Id,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		// Convert amountStr to int64
		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return port.GetAllResponse{}, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = int64(amountFloat)
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

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		date,
	)

	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction port.GetResponse
		var amountStr string
		if err := rows.Scan(
			&transaction.Id,
			&transaction.User_Id,
			&amountStr,
			&transaction.Partner_Id,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		// Convert amountStr to int64
		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return port.GetAllResponse{}, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = int64(amountFloat)
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

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		user_id,
	)

	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction port.GetResponse
		var amountStr string
		if err := rows.Scan(
			&transaction.Id,
			&transaction.User_Id,
			&amountStr,
			&transaction.Partner_Id,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		// Convert amountStr to int64
		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return port.GetAllResponse{}, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = int64(amountFloat)
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

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		partner_id,
	)

	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction port.GetResponse
		var amountStr string
		if err := rows.Scan(
			&transaction.Id,
			&transaction.User_Id,
			&amountStr,
			&transaction.Partner_Id,
			&transaction.Date); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		// Convert amountStr to int64
		amountFloat, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return port.GetAllResponse{}, fmt.Errorf("failed to parse amount: %w", err)
		}
		transaction.Amount = int64(amountFloat)
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
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		req.User_Id,
		req.Partner_Id,
		req.Amount,
	)
	if err != nil {
		return 0, err
	}
	rows.Scan(&id)
	return id, nil
}
