package order

import (
	"context"
	"database/sql"
	"log"

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
	err := p.Pool.QueryRowContext(ctx, query, id).Scan(
		&response.Id,
		&response.RetailerId,
		&response.Status,
		&response.Total)
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

func (p *Postgres) GetByRetailerID(ctx context.Context, retailerId int) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_orders_by_retailer_id($1);"
	rows, err := p.Pool.QueryContext(ctx, query, retailerId)
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
		var order port.GetResponse
		if err := rows.Scan(
			&order.Id,
			&order.RetailerId,
			&order.Status,
			&order.Total); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, order)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_orders_by_status($1);"
	rows, err := p.Pool.QueryContext(ctx, query, status)
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
		var order port.GetResponse
		if err := rows.Scan(
			&order.Id,
			&order.RetailerId,
			&order.Status,
			&order.Total); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, order)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse

	query := "SELECT * FROM public.get_all_orders();"
	rows, err := p.Pool.QueryContext(ctx, query)
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
		var order port.GetResponse
		if err := rows.Scan(
			&order.Id,
			&order.RetailerId,
			&order.Status,
			&order.Total); err != nil {

			log.Printf("unable to scan row: %q", err)
			return port.GetAllResponse{}, err
		}
		response.List = append(response.List, order)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred during rows iteration: %q", err)
		return port.GetAllResponse{}, port.ErrSysUnknown
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var orderId int
	query := "SELECT * FROM public.create_order($1, $2, $3);"

	err := p.Pool.QueryRowContext(ctx, query,
		req.RetailerId,
		req.Status,
		req.Total,
	).Scan(&orderId)
	if err != nil {
		switch err {
		default:
			return 0, port.ErrSysUnknown
		}
	}

	//orderItem
	for _, i := range req.Items {
		query = "SELECT * FROM public.create_order_item($1, $2, $3, $4);"
		_, err := p.Pool.QueryContext(ctx, query,
			orderId,
			i.ProductId,
			i.Quantity,
			i.Price,
		)
		if err != nil {
			switch err {
			default:
				return 0, port.ErrSysUnknown
			}
		}
	}

	return orderId, nil
}

func (p *Postgres) UpdateOrderStatus(ctx context.Context, req *port.UpdateOrderStatusRequest) error {
	query := "SELECT * FROM public.update_order_status($1, $2);"
	_, err := p.Pool.QueryContext(ctx, query, req.Id, req.Status)
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
