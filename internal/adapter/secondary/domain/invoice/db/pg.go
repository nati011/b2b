package invoice

import (
	"context"
	"database/sql"
	"log"

	port "b2b.nati011.github.com/internal/port/domain/invoice/db"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(DB *sql.DB) port.DB {
	return &Postgres{
		db: DB,
	}
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_invoices_by_id($1);"

	err := p.db.QueryRowContext(ctx, query, id).Scan(&response.Id, &response.Status, &response.ExternalId, &response.OrderId, &response.SubTotal)
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

	query := "SELECT * FROM public.get_all_invoices();"
	rows, err := p.db.QueryContext(ctx, query)
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
		var invoice port.GetResponse
		if err := rows.Scan(&invoice.Id, &invoice.Status, &invoice.ExternalId, &invoice.OrderId, &invoice.SubTotal); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, invoice)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_invoices_by_external_id($1);"
	rows, err := p.db.QueryContext(ctx, query, extId)
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
		var invoice port.GetResponse
		if err := rows.Scan(&invoice.Id, &invoice.Status, &invoice.ExternalId, &invoice.OrderId, &invoice.SubTotal); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, invoice)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_invoices_by_status($1);"
	rows, err := p.db.QueryContext(ctx, query, status)
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
		var invoice port.GetResponse
		if err := rows.Scan(&invoice.Id, &invoice.Status, &invoice.ExternalId, &invoice.OrderId, &invoice.SubTotal); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, invoice)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByOrderId(ctx context.Context, orderId int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_invoices_by_order_id($1);"

	err := p.db.QueryRowContext(ctx, query, orderId).Scan(&response.Id, &response.Status, &response.ExternalId, &response.OrderId, &response.SubTotal)
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

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	// invoice
	var resourceId int
	query := "SELECT * FROM public.create_invoice($1, $2, $3, $4);"

	err := p.db.QueryRowContext(ctx, query,
		req.Status,
		req.ExternalId,
		req.OrderId,
		req.Subtotal,
	).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}
	// create invoice line items
	return resourceId, nil

}

func (p *Postgres) UpdateExternalId(ctx context.Context, req *port.UpdateExternalIdRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_invoice_externalId($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.ExternalId).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	return nil
}

func (p *Postgres) UpdateStatus(ctx context.Context, req *port.UpdateStatusRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_invoice_status($1, $2);"

	err := p.db.QueryRowContext(ctx, query, req.Id, req.Status).Scan(&resourceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.ErrSysNoRows
		default:
			return port.ErrSysUnknown
		}
	}

	return nil
}
