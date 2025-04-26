package order

import (
	"context"
	"database/sql"
	"log"

	handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/order"
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
	query := "SELECT * FROM public.get_orders_by_id($1);"
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
		&response.RetailerId,
		&response.Status,
		&response.Total,
		&response.PaymentStatus,
		&response.DeliveryStatus)

	return response, nil
}

func (p *Postgres) GetByRetailerID(ctx context.Context, retailerId int) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_orders_by_retailer_id($1);"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		true,
		retailerId,
	)

	if err != nil {
		return port.GetAllResponse{}, err
	}
	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var order port.GetResponse
		if err := rows.Rows.Scan(
			&order.Id,
			&order.RetailerId,
			&order.Status,
			&order.Total,
			&order.PaymentStatus,
			&order.DeliveryStatus); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, order)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_orders_by_status($1);"
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
		var order port.GetResponse
		if err := rows.Rows.Scan(
			&order.Id,
			&order.RetailerId,
			&order.Status,
			&order.Total,
			&order.PaymentStatus,
			&order.DeliveryStatus); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, order)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_orders();"
	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		true,
	)
	if err != nil {
		return port.GetAllResponse{}, err
	}

	defer rows.Rows.Close()

	for rows.Rows.Next() {
		var order port.GetResponse
		if err := rows.Rows.Scan(
			&order.Id,
			&order.RetailerId,
			&order.Status,
			&order.Total,
			&order.PaymentStatus,
			&order.DeliveryStatus); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, order)
	}

	if err := rows.Rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var orderId int
	query := "SELECT * FROM public.create_order($1, $2, $3, $4, $5);"

	rows, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.RetailerId,
		req.Status,
		req.Total,
		req.PaymentStatus,
		req.DeliveryStatus,
	)

	rows.Row.Scan(&orderId)
	if err != nil {
		return 0, err
	}

	//orderItem
	for _, i := range req.Items {
		query = "SELECT * FROM public.create_order_item($1, $2, $3, $4);"
		_, err := handler.MustQueryRow(
			p.Pool,
			ctx,
			query,
			false,
			orderId,
			i.ProductId,
			i.Quantity,
			i.Price,
		)

		if err != nil {
			return 0, err
		}

	}

	return orderId, nil
}

func (p *Postgres) UpdateOrderStatus(ctx context.Context, req *port.UpdateOrderStatusRequest) error {
	query := "SELECT * FROM public.update_order_status($1, $2);"
	_, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Id,
		req.Status,
	)

	return err
}

func (p *Postgres) UpdatePaymentStatus(ctx context.Context, req *port.UpdateOrderPaymentStatusRequest) error {
	query := "SELECT * FROM public.update_order_payment_status($1, $2);"
	_, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Id,
		req.PaymentStatus,
	)

	return err
}

func (p *Postgres) UpdateDeliveryStatus(ctx context.Context, req *port.UpdateOrderDeliveryStatusRequest) error {
	query := "SELECT * FROM public.update_order_delivery_status($1, $2);"
	_, err := handler.MustQueryRow(
		p.Pool,
		ctx,
		query,
		false,
		req.Id,
		req.DeliveryStatus,
	)

	return err
}
