package order

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"b2b.nati011.github.com/config"
	query_handler "b2b.nati011.github.com/internal/adapter/secondary/sql"
	port "b2b.nati011.github.com/internal/port/domain/order"
)

type GetOrderItem struct {
	Id          int
	ProductId   int
	ProductName string
	Quantity    int
	Price       float64
}

type GetAllOrderItems struct {
	Items []GetOrderItem
}

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
		&response.RetailerName,
		&response.Status,
		&response.Total,
		&response.PaymentStatus,
		&response.DeliveryStatus,
		&response.ConfirmationStatus,
	}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
	if err != nil {
		return port.GetResponse{}, err
	}

	response.Id = *result[0].(*int)
	response.RetailerId = *result[1].(*int)
	response.RetailerName = *result[2].(*string)
	response.Status = *result[3].(*string)
	response.Total = *result[4].(*float64)
	response.PaymentStatus = *result[5].(*string)
	response.DeliveryStatus = *result[6].(*string)
	response.ConfirmationStatus = *result[7].(*string)

	allOrderItems, err := p.GetAllOrderItems(ctx, response.Id)
	if err != nil {
		return port.GetResponse{}, nil
	}

	for _, s := range allOrderItems.Items {
		response.Items = append(response.Items, port.Item{
			ProductId:   s.ProductId,
			ProductName: s.ProductName,
			Quantity:    s.Quantity,
			Price:       s.Price,
		})
	}

	return response, nil
}

func (p *Postgres) GetAllOrderItems(ctx context.Context, orderId int) (GetAllOrderItems, error) {
	var response GetAllOrderItems
	var responseBase GetOrderItem

	query := "SELECT * FROM public.get_order_items_by_order_id($1);"
	args := []any{orderId}

	dest := []any{
		&responseBase.Id,
		&responseBase.ProductId,
		&responseBase.ProductName,
		&responseBase.Quantity,
		&responseBase.Price,
	}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()

	if err != nil {
		return GetAllOrderItems{}, err
	}

	for _, res := range result {
		v, _ := strconv.ParseFloat(res[4].(string), 64)
		response.Items = append(response.Items, GetOrderItem{
			Id:          int(res[0].(int64)),
			ProductId:   int(res[1].(int64)),
			ProductName: res[2].(string),
			Quantity:    int(res[3].(int64)),
			Price:       v,
		})
	}

	return response, nil
}

func (p *Postgres) GetByRetailerID(ctx context.Context, retailerId int) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_orders_by_retailer_id($1);"

	args := []any{&retailerId}
	dest := []any{
		&responseBase.Id,
		&responseBase.RetailerId,
		&responseBase.RetailerName,
		&responseBase.Status,
		&responseBase.Total,
		&responseBase.PaymentStatus,
		&responseBase.DeliveryStatus,
		&responseBase.CreatedAt,
		&responseBase.ConfirmationStatus,
		&responseBase.PaymentMethod,
	}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()

	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		v, _ := strconv.ParseFloat(res[4].(string), 64)
		val := port.GetResponse{
			Id:                 int(res[0].(int64)),
			RetailerId:         int(res[1].(int64)),
			RetailerName:       res[2].(string),
			Status:             res[3].(string),
			Total:              v,
			PaymentStatus:      res[5].(string),
			DeliveryStatus:     res[6].(string),
			CreatedAt:          res[7].(time.Time),
			ConfirmationStatus: res[8].(string),
			PaymentMethod:      res[9].(string),
		}
		allOrderItems, err := p.GetAllOrderItems(ctx, val.Id)
		if err != nil {
			return port.GetAllResponse{}, err
		}

		for _, s := range allOrderItems.Items {
			val.Items = append(val.Items, port.Item{
				ProductId:   s.ProductId,
				ProductName: s.ProductName,
				Quantity:    s.Quantity,
				Price:       s.Price,
			})
		}

		response.List = append(response.List, val)
	}

	return response, nil
}

func (p *Postgres) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse

	query := "SELECT * FROM public.get_orders_by_status($1);"
	args := []any{&status}
	dest := []any{
		&responseBase.Id,
		&responseBase.RetailerId,
		&responseBase.RetailerName,
		&responseBase.Status,
		&responseBase.Total,
		&responseBase.PaymentStatus,
		&responseBase.DeliveryStatus,
		&responseBase.ConfirmationStatus,
	}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()

	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		v, _ := strconv.ParseFloat(res[3].(string), 64)
		val := port.GetResponse{
			Id:                 int(res[0].(int64)),
			RetailerId:         int(res[1].(int64)),
			RetailerName:       res[2].(string),
			Status:             res[3].(string),
			Total:              v,
			PaymentStatus:      res[5].(string),
			DeliveryStatus:     res[6].(string),
			ConfirmationStatus: res[7].(string),
		}
		allOrderItems, err := p.GetAllOrderItems(ctx, val.Id)
		if err != nil {
			return port.GetAllResponse{}, err
		}

		for _, s := range allOrderItems.Items {
			val.Items = append(val.Items, port.Item{
				ProductId: s.ProductId,
				Quantity:  s.Quantity,
				Price:     s.Price,
			})
		}

		response.List = append(response.List, val)
	}

	return response, nil
}

func (p *Postgres) GetAll(ctx context.Context) (port.GetAllResponse, error) {
	var response port.GetAllResponse
	var responseBase port.GetResponse
	var totalCount int64

	query := "SELECT * FROM public.get_all_orders($1,$2);"
	args := []any{p.Pagination.Limit, p.Pagination.Offset}

	dest := []any{
		&responseBase.Id,
		&responseBase.RetailerId,
		&responseBase.RetailerName,
		&responseBase.Status,
		&responseBase.Total,
		&responseBase.PaymentStatus,
		&responseBase.DeliveryStatus,
		&responseBase.CreatedAt,
		&totalCount,
		&responseBase.ConfirmationStatus,
	}

	result, err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithMultiRowResultSet(args, dest),
	).DoMultiQuery()
	if err != nil {
		return port.GetAllResponse{}, err
	}

	for _, res := range result {
		v, _ := strconv.ParseFloat(res[4].(string), 64)
		val := port.GetResponse{
			Id:                 int(res[0].(int64)),
			RetailerId:         int(res[1].(int64)),
			RetailerName:       res[2].(string),
			Status:             res[3].(string),
			Total:              v,
			PaymentStatus:      res[5].(string),
			DeliveryStatus:     res[6].(string),
			CreatedAt:          res[7].(time.Time),
			ConfirmationStatus: res[9].(string),
		}
		totalCount = res[8].(int64)
		allOrderItems, err := p.GetAllOrderItems(ctx, val.Id)
		if err != nil {
			return port.GetAllResponse{}, err
		}

		for _, s := range allOrderItems.Items {
			val.Items = append(val.Items, port.Item{
				ProductId:   s.ProductId,
				ProductName: s.ProductName,
				Quantity:    s.Quantity,
				Price:       s.Price,
			})
		}

		response.List = append(response.List, val)
		response.TotalCount = totalCount
	}

	return response, nil
}

func (p *Postgres) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	var orderId int
	query := "SELECT * FROM public.create_order($1, $2, $3, $4, $5, $6);"
	args := []any{
		req.RetailerId,
		req.Status,
		req.Total,
		req.PaymentStatus,
		req.DeliveryStatus,
		req.ConfirmationStatus}

	result := []any{&orderId}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, result),
	).DoSingleQuery()
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
		).DoSingleQuery()

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
	).DoSingleQuery()

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
	).DoSingleQuery()

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
	).DoSingleQuery()

	return err
}

func (p *Postgres) UpdateConfirmationStatus(ctx context.Context, req *port.UpdateOrderConfirmationStatusRequest) error {
	query := "SELECT * FROM public.update_order_confirmation_status($1, $2);"
	args := []any{req.Id, req.ConfirmationStatus}

	err := query_handler.NewQuery(
		query_handler.WithCtx(ctx),
		query_handler.WithDB(p.Pool),
		query_handler.WithQuery(query),
		query_handler.WithSingleRowResultSet(args, nil),
	).DoSingleQuery()

	return err
}
