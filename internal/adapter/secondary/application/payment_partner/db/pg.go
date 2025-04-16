package db

import (
	"context"
	"database/sql"
	"log"

	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/partner/db"
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

	query := "SELECT * FROM public.get_payment_partner_by_id($1);"

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
		&response.Name,
		&response.Icon,
		&response.Status,
		&response.Init_payment_url)

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_payment_partners();"
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
		var partner port.GetResponse
		if err := rows.Scan(
			&partner.Id,
			&partner.Name,
			&partner.Icon,
			&partner.Status,
			&partner.Init_payment_url); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, partner)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_payment_partner_by_status($1);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		status,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var partner port.GetResponse
		if err := rows.Scan(
			&partner.Id,
			&partner.Name,
			&partner.Icon,
			&partner.Status,
			&partner.Init_payment_url); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, partner)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByName(ctx context.Context, name string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_payment_partner_by_name($1);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		name,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var partner port.GetResponse
		if err := rows.Scan(
			&partner.Id,
			&partner.Name,
			&partner.Icon,
			&partner.Status,
			&partner.Init_payment_url); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, partner)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.create_payment_partner($1, $2, $3, $4);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		req.Name,
		req.Icon,
		req.Status,
		req.Init_payment_url,
	)
	if err != nil {
		return 0, err
	}

	log.Printf("Rows: %v", rows)
	rows.Scan(&partner_id)
	return partner_id, nil
}

func (p *Postgres) UpdateStatus(ctx context.Context, id int, status string) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.update_payment_partner_status($1, $2);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		id,
		status,
	)
	if err != nil {
		return 0, err
	}
	rows.Scan(&partner_id)
	return partner_id, nil
}

func (p *Postgres) UpdateName(ctx context.Context, id int, name string) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.update_payment_partner_name($1, $2);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		id,
		name,
	)
	if err != nil {
		return 0, err
	}
	rows.Scan(&partner_id)

	return partner_id, nil
}

func (p *Postgres) Delete(context.Context, int) error {
	return nil
}
