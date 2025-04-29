package invoice

import (
	"context"
	"database/sql"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/invoice/db"
)

type Postgres struct {
	db         *sql.DB
	Pagination *config.Pagination
}

func NewPostgres(DB *sql.DB, pagination *config.Pagination) port.DB {
	return &Postgres{
		db:         DB,
		Pagination: pagination,
	}
}

func (p *Postgres) Get(ctx context.Context, id int) (port.GetResponse, error) {
	var response port.GetResponse
	query := "SELECT * FROM public.get_invoices_by_id($1);"
	result := []any{
		&response.Id,
		&response.Status,
		&response.ExternalId,
		&response.OrderId,
		&response.SubTotal,
		&response.TaxAmount}
	args := []any{id}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return port.GetResponse{}, err
	}

	items, err := p.getInvoiceLineItemsByProductId(ctx, response.Id)
	if err != nil {
		return port.GetResponse{}, err
	}
	response.LineItems = items

	response.Id = *result[0].(*int)
	response.Status = *result[1].(*string)
	response.ExternalId = *result[2].(*string)
	response.OrderId = *result[3].(*int)
	response.SubTotal = *result[4].(*float64)
	response.TaxAmount = *result[5].(*float64)

	return response, nil
}

func (p *Postgres) getInvoiceLineItemsByProductId(ctx context.Context, invoice_Id int) ([]port.Item, error) {
	var response []port.Item
	var responseBase port.Item
	var id int
	var invoiceId int
	query := "SELECT * FROM public.get_invoice_line_item_by_invoice_id($1);"
	result := [][]any{
		{
			&id,
			&responseBase.ProductName,
			&responseBase.ProductQuantity,
			&responseBase.ProductPrice,
			&responseBase.ProductId,
			&invoiceId},
	}

	args := []any{invoice_Id}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return []port.Item{}, err
	}

	//convert
	for _, res := range result {
		responseBase := port.Item{
			ProductName:     *res[1].(*string),
			ProductQuantity: *res[2].(*int),
			ProductPrice:    *res[3].(*float64),
			ProductId:       *res[4].(*int),
		}
		response = append(response, responseBase)
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_all_invoices($1, $2);"
	result := [][]any{
		{&responseBase.Id,
			&responseBase.Status,
			&responseBase.ExternalId,
			&responseBase.OrderId,
			&responseBase.SubTotal,
			&responseBase.TaxAmount,
		},
	}
	args := []any{p.Pagination.Limit, p.Pagination.Offset}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return port.GetAllResponse{}, err
	}
	for _, res := range result {
		resp := port.GetResponse{
			Id:         *res[0].(*int),
			ExternalId: *res[1].(*string),
			Status:     *res[2].(*string),
			OrderId:    *res[3].(*int),
			SubTotal:   *res[4].(*float64),
			TaxAmount:  *res[5].(*float64),
		}
		items, err := p.getInvoiceLineItemsByProductId(ctx, resp.Id)
		if err != nil {
			return port.GetAllResponse{}, err
		}
		resp.LineItems = items
		response.List = append(response.List, resp)
	}
	return response, nil
}

func (p *Postgres) GetByExternalId(ctx context.Context, extId string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_invoices_by_external_id($1);"

	result := [][]any{
		{&responseBase.Id,
			&responseBase.Status,
			&responseBase.ExternalId,
			&responseBase.OrderId,
			&responseBase.SubTotal,
			&responseBase.TaxAmount,
		},
	}
	args := []any{extId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return port.GetAllResponse{}, err
	}
	for _, res := range result {
		resp := port.GetResponse{
			Id:         *res[0].(*int),
			ExternalId: *res[1].(*string),
			Status:     *res[2].(*string),
			OrderId:    *res[3].(*int),
			SubTotal:   *res[4].(*float64),
			TaxAmount:  *res[5].(*float64),
		}
		items, err := p.getInvoiceLineItemsByProductId(ctx, resp.Id)
		if err != nil {
			return port.GetAllResponse{}, err
		}
		resp.LineItems = items
		response.List = append(response.List, resp)
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_invoices_by_status($1, $2, $3);"

	result := [][]any{
		{&responseBase.Id,
			&responseBase.Status,
			&responseBase.ExternalId,
			&responseBase.OrderId,
			&responseBase.SubTotal,
			&responseBase.TaxAmount,
		},
	}
	args := []any{status, p.Pagination.Limit, p.Pagination.Offset}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		resp := port.GetResponse{
			Id:         *res[0].(*int),
			ExternalId: *res[1].(*string),
			Status:     *res[2].(*string),
			OrderId:    *res[3].(*int),
			SubTotal:   *res[4].(*float64),
			TaxAmount:  *res[5].(*float64),
		}
		items, err := p.getInvoiceLineItemsByProductId(ctx, resp.Id)
		if err != nil {
			return port.GetAllResponse{}, err
		}
		resp.LineItems = items
		response.List = append(response.List, resp)
	}

	return response, nil
}

func (p *Postgres) GetByOrderId(ctx context.Context, orderId int) (port.GetResponse, error) {
	var response port.GetResponse

	query := "SELECT * FROM public.get_invoices_by_order_id($1);"
	result := []any{
		&response.Id,
		&response.Status,
		&response.ExternalId,
		&response.OrderId,
		&response.SubTotal,
		&response.TaxAmount}
	args := []any{&orderId}
	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return port.GetResponse{}, err
	}

	items, err := p.getInvoiceLineItemsByProductId(ctx, response.Id)
	if err != nil {
		return port.GetResponse{}, err
	}
	response.LineItems = items

	response.Id = *result[0].(*int)
	response.Status = *result[1].(*string)
	response.ExternalId = *result[2].(*string)
	response.OrderId = *result[3].(*int)
	response.SubTotal = *result[4].(*float64)
	response.TaxAmount = *result[5].(*float64)

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	// invoice
	var invoiceId int
	query := "SELECT * FROM public.create_invoice($1, $2, $3, $4, $5);"

	result := []any{&invoiceId}
	args := []any{
		req.Status,
		req.ExternalId,
		req.OrderId,
		req.Subtotal,
		req.TaxAmount}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return 0, err
	}

	return invoiceId, nil
}

func (p *Postgres) UpdateExternalId(ctx context.Context, req *port.UpdateExternalIdRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_invoice_externalId($1, $2);"

	result := []any{&resourceId}
	args := []any{
		req.Id,
		req.ExternalId,
	}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return err
	}
	return err
}

func (p *Postgres) UpdateStatus(ctx context.Context, req *port.UpdateStatusRequest) error {
	var resourceId int
	query := "SELECT * FROM public.update_invoice_status($1, $2);"

	result := []any{&resourceId}
	args := []any{
		req.Id,
		req.Status,
	}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.db),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return err
	}

	return err
}
