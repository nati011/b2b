package invoice

import (
	"context"
	"database/sql"
	"log"

	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
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

	rows, err := handler.MustQueryRow(
		p.db,
		ctx,
		query,
		false,
		id,
	)

	rows.Row.Scan(&response.Id, &response.Status, &response.ExternalId, &response.OrderId, &response.SubTotal, &response.TaxAmount)
	if err != nil {
		return port.GetResponse{}, err
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
	rows, err := handler.MustQueryRow(
		p.db,
		ctx,
		query,
		true,
		invoice_Id,
	)
	if err != nil {
		return []port.Item{}, err
	}
	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var id int
		var invoiceId int
		var item port.Item
		if err := rows.Rows.Scan(&id, &item.ProductName, &item.ProductQuantity, &item.ProductPrice, &item.ProductId, &invoiceId); err != nil {
			log.Printf("unable to scan row: %q", err)
			return []port.Item{}, err
		}
		response = append(response, item)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return []port.Item{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_invoices();"
	rows, err := handler.MustQueryRow(
		p.db,
		ctx,
		query,
		true,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}
	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var invoice port.GetResponse
		if err := rows.Rows.Scan(&invoice.Id, &invoice.Status, &invoice.ExternalId, &invoice.OrderId, &invoice.SubTotal, &invoice.TaxAmount); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, invoice)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_invoices_by_external_id($1);"
	rows, err := handler.MustQueryRow(
		p.db,
		ctx,
		query,
		true,
		extId,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}
	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var invoice port.GetResponse
		if err := rows.Rows.Scan(&invoice.Id, &invoice.Status, &invoice.ExternalId, &invoice.OrderId, &invoice.SubTotal, &invoice.TaxAmount); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, invoice)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_invoices_by_status($1);"
	rows, err := handler.MustQueryRow(
		p.db,
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
		var invoice port.GetResponse
		if err := rows.Rows.Scan(&invoice.Id, &invoice.Status, &invoice.ExternalId, &invoice.OrderId, &invoice.SubTotal, &invoice.TaxAmount); err != nil {
			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, invoice)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByOrderId(ctx context.Context, orderId int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_invoices_by_order_id($1);"
	rows, err := handler.MustQueryRow(
		p.db,
		ctx,
		query,
		false,
		orderId,
	)

	if err != nil {
		return port.GetResponse{}, err
	}

	rows.Row.Scan(&response.Id, &response.Status, &response.ExternalId, &response.OrderId, &response.SubTotal, &response.TaxAmount)

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	// invoice
	var invoiceId int
	query := "SELECT * FROM public.create_invoice($1, $2, $3, $4, $5);"

	rows, err := handler.MustQueryRow(
		p.db,
		ctx,
		query,
		false,
		req.Status,
		req.ExternalId,
		req.OrderId,
		req.Subtotal,
		req.TaxAmount,
	)
	if err != nil {
		return 0, err

	}
	rows.Row.Scan(&invoiceId)

	// invoice line items
	for _, i := range req.LineItems {
		query = "SELECT * FROM public.create_invoice_line_item($1, $2, $3, $4, $5);"

		_, err := handler.MustQueryRow(
			p.db,
			ctx,
			query,
			false,
			i.ProductName,
			i.ProductQuantity,
			i.ProductPrice,
			i.ProductId,
			invoiceId,
		)

		if err != nil {
			return 0, err
		}
	}

	return invoiceId, nil
}

func (p *Postgres) UpdateExternalId(ctx context.Context, req *port.UpdateExternalIdRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_invoice_externalId($1, $2);"
	rows, err := handler.MustQueryRow(
		p.db,
		ctx,
		query,
		false,
		req.Id,
		req.ExternalId,
	)

	rows.Row.Scan(&resourceId)
	return err
}

func (p *Postgres) UpdateStatus(ctx context.Context, req *port.UpdateStatusRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_invoice_status($1, $2);"

	rows, err := handler.MustQueryRow(
		p.db,
		ctx,
		query,
		false,
		req.Id,
		req.Status,
	)

	rows.Row.Scan(&resourceId)

	return err
}
