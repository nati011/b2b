package db

import (
	"context"
	"database/sql"
	"log"

	"b2b.nati011.github.com/config"
	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/application/partner/db"
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

	query := "SELECT * FROM public.get_payment_partner_by_id($1);"

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
		&response.Name,
		&response.Icon,
		&response.Status,
		&response.Init_payment_url)

	return response, nil
}
func (p *Postgres) GetPartnerSecret(ctx context.Context, id int) (port.GetPartnerSecret, error) {
	var response port.GetPartnerSecret

	query := "SELECT * FROM public.get_payment_partner_secret($1);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		id,
	)

	if err != nil {
		return port.GetPartnerSecret{}, err
	}

	rows.Row.Scan(
		&response.Name,
		&response.Secret,
		&response.Init_payment_url,
	)

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	limit := p.Pagination.Limit
	offset := p.Pagination.Offset

	query := "SELECT * FROM public.get_all_payment_partners($1,$2);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		true,
		limit,
		offset,
	)
	if err != nil {
		return port.GetAllResponse{}, err

	}
	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var partner port.GetResponse
		if err := rows.Rows.Scan(
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

	if err := rows.Rows.Err(); err != nil {
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
		true,
		status,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}
	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var partner port.GetResponse
		if err := rows.Rows.Scan(
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

	if err := rows.Rows.Err(); err != nil {
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
		true,
		name,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var partner port.GetResponse
		if err := rows.Rows.Scan(
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

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.create_payment_partner($1, $2, $3, $4, $5);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Name,
		req.Icon,
		req.Status,
		req.Init_payment_url,
		req.Secret,
	)
	if err != nil {
		return 0, err
	}

	log.Printf("Rows: %v", rows)
	rows.Row.Scan(&partner_id)
	return partner_id, nil
}

func (p *Postgres) UpdateStatus(ctx context.Context, id int, status string) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.update_payment_partner_status($1, $2);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		id,
		status,
	)
	if err != nil {
		return 0, err
	}
	rows.Row.Scan(&partner_id)
	return partner_id, nil
}

func (p *Postgres) UpdateName(ctx context.Context, id int, name string) (int, error) {
	var partner_id int
	query := "SELECT * FROM public.update_payment_partner_name($1, $2);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		id,
		name,
	)
	if err != nil {
		return 0, err
	}
	rows.Row.Scan(&partner_id)

	return partner_id, nil
}

func (p *Postgres) Delete(context.Context, int) error {
	return nil
}
