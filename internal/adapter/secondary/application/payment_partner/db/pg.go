package db

import (
	"context"
	"database/sql"
	"log"

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

	err := p.Pool.QueryRowContext(ctx, query, id).Scan(
		&response.Id,
		&response.Name,
		&response.Icon,
		&response.Status,
		&response.Init_payment_url)
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

func (p *Postgres) GetAll(ctx context.Context, pagination *port.Pagination) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	// Default limit and offset
	limit := 10
	offset := pagination.Offset

	if pagination.Limit != 0 {
		limit = pagination.Limit
	}

	query := "SELECT * FROM public.get_all_payment_partners($1, $2);"
	rows, err := p.Pool.QueryContext(ctx, query, limit, offset)
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

func (p *Postgres) GetByStatus(ctx context.Context, status string, pagination *port.Pagination) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	// Default limit and offset
	limit := 10
	offset := pagination.Offset

	if pagination.Limit != 0 {
		limit = pagination.Limit
	}

	query := "SELECT * FROM public.get_payment_partner_by_status($1, $2, $3);"
	rows, err := p.Pool.QueryContext(ctx, query, status, limit, offset)
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

func (p *Postgres) GetByName(ctx context.Context, name string, pagination *port.Pagination) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	// Default limit and offset
	limit := 10
	offset := pagination.Offset

	if pagination.Limit != 0 {
		limit = pagination.Limit
	}

	query := "SELECT * FROM public.get_payment_partner_by_name($1, $2, $3);"
	rows, err := p.Pool.QueryContext(ctx, query, name, limit, offset)
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

	err := p.Pool.QueryRowContext(ctx, query, req.Name, req.Icon, req.Status, req.Init_payment_url).Scan(
		&partner_id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return partner_id, nil
}

func (p *Postgres) UpdateStatus(ctx context.Context, id int, status string) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.update_payment_partner_status($1, $2);"

	_, err := p.Pool.QueryContext(ctx, query, id, status)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return partner_id, nil
}

func (p *Postgres) UpdateName(ctx context.Context, id int, name string) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.update_payment_partner_name($1, $2);"

	_, err := p.Pool.QueryContext(ctx, query, id, name)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	return partner_id, nil
}

func (p *Postgres) Delete(context.Context, int) error {
	return nil
}
