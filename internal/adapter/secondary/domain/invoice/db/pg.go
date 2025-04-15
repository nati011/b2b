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

	err := p.db.QueryRowContext(ctx, query, id).Scan(&response.Id, &response.Status, &response.ExternalId, &response.OrderId, &response.SubTotal, &response.TaxAmount)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return port.GetResponse{}, port.ErrSysNoRows
		default:
			return port.GetResponse{}, port.ErrSysUnknown
		}
	}
	items, err := p.getInvoiceLineItemsByProductId(ctx, response.Id)
	if err != nil {
		return port.GetResponse{}, err
	}
	response.LineItems = items

	return response, nil
}

func (p *Postgres) getInvoiceLineItemsByProductId(ctx context.Context, invoice_Id int) ([]port.Item, error) {
	var response []port.Item

	query := "SELECT * FROM public.get_invoice_line_item_by_invoice_id($1);"
	rows, err := p.db.QueryContext(ctx, query, invoice_Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []port.Item{}, port.ErrSysNoRows
		default:
			return []port.Item{}, port.ErrSysUnknown
		}

	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var invoiceId int
		var item port.Item
		if err := rows.Scan(&id, &item.ProductName, &item.ProductQuantity, &item.ProductPrice, &item.ProductId, &invoiceId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return []port.Item{}, err
		}
		response = append(response, item)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return []port.Item{}, port.ErrSysUnknown
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
	query := "SELECT * FROM public.get_all_invoices($1,$2);"
	rows, err := p.db.QueryContext(ctx, query, limit, offset)
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
		if err := rows.Scan(&invoice.Id, &invoice.Status, &invoice.ExternalId, &invoice.OrderId, &invoice.SubTotal, &invoice.TaxAmount); err != nil {
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
		if err := rows.Scan(&invoice.Id, &invoice.Status, &invoice.ExternalId, &invoice.OrderId, &invoice.SubTotal, &invoice.TaxAmount); err != nil {
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
		if err := rows.Scan(&invoice.Id, &invoice.Status, &invoice.ExternalId, &invoice.OrderId, &invoice.SubTotal, &invoice.TaxAmount); err != nil {
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

	err := p.db.QueryRowContext(ctx, query, orderId).Scan(&response.Id, &response.Status, &response.ExternalId, &response.OrderId, &response.SubTotal, &response.TaxAmount)
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
	var invoiceId int
	query := "SELECT * FROM public.create_invoice($1, $2, $3, $4, $5);"

	err := p.db.QueryRowContext(ctx, query,
		req.Status,
		req.ExternalId,
		req.OrderId,
		req.Subtotal,
		req.TaxAmount,
	).Scan(&invoiceId)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, port.ErrSysNoRows
		default:
			return 0, port.ErrSysUnknown
		}
	}

	// invoice line items
	for _, i := range req.LineItems {
		query = "SELECT * FROM public.create_invoice_line_item($1, $2, $3, $4, $5);"

		_, err = p.db.QueryContext(ctx, query,
			i.ProductName,
			i.ProductQuantity,
			i.ProductPrice,
			i.ProductId,
			invoiceId,
		)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return 0, port.ErrSysNoRows
			default:
				return 0, port.ErrSysUnknown
			}
		}
	}

	return invoiceId, nil
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
