package order

import (
	"context"
	"database/sql"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/order"
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
	query := "SELECT * FROM public.get_orders_by_id($1);"
	args := []any{&id}
	result := []any{
		&response.Id,
		&response.RetailerId,
		&response.Status,
		&response.Total,
		&response.PaymentStatus,
		&response.DeliveryStatus,
	}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoStuff()
	if err != nil {
		return port.GetResponse{}, nil
	}

	response.Id = *result[0].(*int)
	response.RetailerId = *result[1].(*int)
	response.Status = *result[2].(*string)
	response.Total = *result[3].(*float64)
	response.PaymentStatus = *result[4].(*string)
	response.DeliveryStatus = *result[5].(*string)

	return response, nil
}

func (p *Postgres) GetByRetailerID(ctx context.Context, retailerId int) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_orders_by_retailer_id($1);"

	args := []any{&retailerId}
	result := [][]any{{
		&responseBase.Id,
		&responseBase.RetailerId,
		&responseBase.Status,
		&responseBase.Total,
		&responseBase.PaymentStatus,
		&responseBase.DeliveryStatus,
	}}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, result),
	).DoStuff()

	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		response.List = append(response.List, port.GetResponse{
			Id:             *res[0].(*int),
			RetailerId:     *res[1].(*int),
			Status:         *res[2].(*string),
			Total:          *res[3].(*float64),
			PaymentStatus:  *res[4].(*string),
			DeliveryStatus: *res[5].(*string),
		})
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_orders_by_status($1);"
	args := []any{&status}
	result := [][]any{{
		&responseBase.Id,
		&responseBase.RetailerId,
		&responseBase.Status,
		&responseBase.Total,
		&responseBase.PaymentStatus,
		&responseBase.DeliveryStatus,
	}}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, result),
	).DoStuff()

	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		response.List = append(response.List, port.GetResponse{
			Id:             *res[0].(*int),
			RetailerId:     *res[1].(*int),
			Status:         *res[2].(*string),
			Total:          *res[3].(*float64),
			PaymentStatus:  *res[4].(*string),
			DeliveryStatus: *res[5].(*string),
		})
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_all_orders($1,$2);"
	args := []any{p.Pagination.Limit, p.Pagination.Offset}
	result := [][]any{{
		&responseBase.Id,
		&responseBase.RetailerId,
		&responseBase.Status,
		&responseBase.Total,
		&responseBase.PaymentStatus,
		&responseBase.DeliveryStatus,
	}}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, result),
	).DoStuff()

	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		response.List = append(response.List, port.GetResponse{
			Id:             *res[0].(*int),
			RetailerId:     *res[1].(*int),
			Status:         *res[2].(*string),
			Total:          *res[3].(*float64),
			PaymentStatus:  *res[4].(*string),
			DeliveryStatus: *res[5].(*string),
		})
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var orderId int
	query := "SELECT * FROM public.create_order($1, $2, $3, $4, $5);"
	args := []any{req.RetailerId, req.Status, req.Total, req.PaymentStatus, req.DeliveryStatus}
	result := []any{&orderId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoStuff()

	if err != nil {
		return 0, err
	}

	//orderItem
	for _, i := range req.Items {
		query = "SELECT * FROM public.create_order_item($1, $2, $3, $4);"
		itemArgs := []any{orderId, i.ProductId, i.Quantity, i.Price}

		err := query_handler.NewQuery(
			query_handler.WithCtx(ctx),
			query_handler.WithDB(p.Pool),
			query_handler.WithQuery(query),
			query_handler.WithSingleRowResultSet(itemArgs, nil),
		).DoStuff()

		if err != nil {
			return 0, err
		}

	}

	return orderId, nil
}

func (p *Postgres) UpdateOrderStatus(ctx context.Context, req *port.UpdateOrderStatusRequest) error {
	query := "SELECT * FROM public.update_order_status($1, $2);"
	args := []any{req.Id, req.Status}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoStuff()

	return err
}

func (p *Postgres) UpdatePaymentStatus(ctx context.Context, req *port.UpdateOrderPaymentStatusRequest) error {
	query := "SELECT * FROM public.update_order_payment_status($1, $2);"
	args := []any{req.Id, req.PaymentStatus}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoStuff()

	return err
}

func (p *Postgres) UpdateDeliveryStatus(ctx context.Context, req *port.UpdateOrderDeliveryStatusRequest) error {
	query := "SELECT * FROM public.update_order_delivery_status($1, $2);"
	args := []any{req.Id, req.DeliveryStatus}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoStuff()

	return err
}
