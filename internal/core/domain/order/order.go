package order

import (
	"context"
	"errors"

	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
	port "b2b.nati011.github.com/internal/port/domain/order"
)

var (
	ErrIdNotFound                         = errors.New("oopsy, id not found")
	ErrRetailerIdNotSupplied              = errors.New("oopsy, retailer id mandatory")
	ErrRetailerIdNotFound                 = errors.New("oopsy, retailer does not exist")
	ErrAtleastOneOrderItemNeeded          = errors.New("oopsy, order items cannot be empty")
	ErrItemMemberProductIdOrQuantityEmpty = errors.New("oopsy, either order item member productId or quantity missing")
	ErrUnknown                            = errors.New("oopsy, unknown error")
	ErrEmptyGetResponse                   = errors.New("oopsy, empty get response")
	ErrAlreadyCanceled                    = errors.New("oopsy, order already canceled")
	ErrItemMemberProductNotFound          = errors.New("oopsy, product not found")
	ErrItemMemberProductQuantityNotFound  = errors.New("oopsy, product quantity not found")
)

// order status
// order status
const (
	CANCELD_STATUS   = "CANCELED"
	PENDING_STATUS   = "PENDING"
	COMPLETED_STATUS = "COMPLETED"
)

// payment
const (
	PAYMENT_PENDING_STATUS  = "PENDING"
	PAYMENT_ACCEPTED_STATUS = "ACCEPTED"
)

// delivery
const (
	DELIVERY_PENDING_STATUS    = "PENDING"
	DELIVERY_DISPATCHED_STATUS = "DISPATCHED"
	DELIVERY_COMPLETED_STATUS  = "COMPLETED"
)

// payment
const (
	PAYMENT_PENDING_STATUS  = "PENDING"
	PAYMENT_ACCEPTED_STATUS = "ACCEPTED"
)

// delivery
const (
	DELIVERY_PENDING_STATUS    = "PENDING"
	DELIVERY_DISPATCHED_STATUS = "DISPATCHED"
	DELIVERY_COMPLETED_STATUS  = "COMPLETED"
)

type Item struct {
	ProductId int
	Quantity  int
}

type PlaceRequest struct {
	RetailerId int
	Items      []Item
}

type GetResponse struct {
	Id             int
	RetailerId     int
	Items          []Item
	Total          float32
	Status         string
	DeliveryStatus string
	PaymentStatus  string
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByParamRequest struct {
	RetailerId int
	Status     string
}

type UpdateRequest struct {
	Id             int
	Status         string
	PaymentStatus  string
	DeliveryStatus string
	Id             int
	Status         string
	PaymentStatus  string
	DeliveryStatus string
}

type Provider interface {
	Place(ctx context.Context, req *PlaceRequest) (int, error)
	Cancel(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
	UpdateStatus(ctx context.Context, req *UpdateRequest) (int, error)
	UpdatePaymentStatus(ctx context.Context, req *UpdateRequest) (int, error)
	UpdateDeliveryStatus(ctx context.Context, req *UpdateRequest) (int, error)
	UpdateStatus(ctx context.Context, req *UpdateRequest) (int, error)
	UpdatePaymentStatus(ctx context.Context, req *UpdateRequest) (int, error)
	UpdateDeliveryStatus(ctx context.Context, req *UpdateRequest) (int, error)
}

type OrderService struct {
	DB              port.DB
	InvoiceService  invoice.Provider
	ProductService  product.Provider
	RetailerService retailer.Provider
}

func NewOrderService(
	db port.DB,
	is invoice.Provider,
	ps product.Provider,
	rs retailer.Provider) Provider {

	return &OrderService{
		DB:              db,
		InvoiceService:  is,
		ProductService:  ps,
		RetailerService: rs,
	}
}

func (o *OrderService) validate_placement(ctx context.Context, req *PlaceRequest) error {
	if err := o.validate_retailerId(ctx, req.RetailerId); err != nil {
		return err
	}

	if err := o.validate_items(ctx, req.Items); err != nil {
		return err
	}
	return nil
}

func (o *OrderService) Place(ctx context.Context, req *PlaceRequest) (int, error) {
	err := o.validate_placement(ctx, req)
	if err != nil {
		return 0, err
	}

	items := make([]port.Item, 0, len(req.Items))
	var itemsTotal float64
	for _, i := range req.Items {
		//fetch price from product
		prod_resp, err := o.ProductService.Get(ctx, i.ProductId)
		if err != nil {
			switch err {
			default:
				return 0, ErrUnknown
			}
		}
		items = append(items, port.Item{
			ProductId: i.ProductId,
			Quantity:  i.Quantity,
			Price:     prod_resp.Price,
		})
		itemsTotal += prod_resp.Price * float64(i.Quantity)
	}

	order_id, err := o.DB.Create(ctx, &port.CreateRequest{
		RetailerId:     req.RetailerId,
		Items:          items,
		Status:         PENDING_STATUS,
		PaymentStatus:  PAYMENT_PENDING_STATUS,
		DeliveryStatus: DELIVERY_PENDING_STATUS,
		Total:          itemsTotal,
		RetailerId:     req.RetailerId,
		Items:          items,
		Status:         PENDING_STATUS,
		PaymentStatus:  PAYMENT_PENDING_STATUS,
		DeliveryStatus: DELIVERY_PENDING_STATUS,
		Total:          itemsTotal,
	})
	if err != nil {
		return 0, ErrUnknown
	}

	// create invoice
	lineItems := make([]invoice.Item, 0, len(req.Items))
	for _, i := range req.Items {
		prod_resp, err := o.ProductService.Get(ctx, i.ProductId)
		if err != nil {
			switch err {
			default:
				return 0, ErrUnknown
			}
		}
		lineItems = append(lineItems, invoice.Item{
			ProductId:       i.ProductId,
			ProductName:     prod_resp.Name,
			ProductQuantity: i.Quantity,
			ProductPrice:    prod_resp.Price,
		})
	}

	if _, err = o.InvoiceService.Create(ctx, &invoice.CreateRequest{
		Status:    invoice.DRAFT_STATUS,
		OrderId:   order_id,
		SubTotal:  itemsTotal,
		LineItems: lineItems,
		TaxAmount: 0, //default
	}); err != nil {
		o.Cancel(ctx, order_id)
		return 0, ErrUnknown
	}

	//TODO
	/*
		reserve stock
	*/

	// o.InvoiceService.Cancel(ctxm order_id)
	// o.Cancel(ctx, order_id)

	//TODO
	/*
		send sms
	*/

	//TODO
	/*
		send email
	*/
	return order_id, nil
}

func (o *OrderService) Cancel(ctx context.Context, id int) error {
	//validate
	got, err := o.DB.GetByID(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//check if already given status
	if got.Status == CANCELD_STATUS {
		return ErrAlreadyCanceled
	}

	err = o.DB.UpdateOrderStatus(ctx, &port.UpdateOrderStatusRequest{
		Id:     id,
		Status: CANCELD_STATUS,
	})
	if err != nil {
		switch err {
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (o *OrderService) Get(ctx context.Context, id int) (GetResponse, error) {
	resp, err := o.DB.GetByID(ctx, id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	items := []Item{}
	for _, i := range resp.Items {
		items = append(items, Item{
			ProductId: i.ProductId,
			Quantity:  i.Quantity,
		})
	}
	return GetResponse{
		Id:             resp.Id,
		RetailerId:     resp.RetailerId,
		Items:          items,
		Status:         resp.Status,
		DeliveryStatus: resp.DeliveryStatus,
		PaymentStatus:  resp.PaymentStatus,
	}, nil
}

func (o *OrderService) GetAll(ctx context.Context) (GetAllResponse, error) {
	resp, err := o.DB.GetAll(ctx)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetResponse
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}

	return_response := GetAllResponse{}
	for _, i := range resp.List {
		items := []Item{}
		for _, i := range i.Items {
			items = append(items, Item{
				ProductId: i.ProductId,
				Quantity:  i.Quantity,
			})
		}
		return_response.List = append(return_response.List, GetResponse{
			Id:             i.Id,
			RetailerId:     i.RetailerId,
			Items:          items,
			Status:         i.Status,
			DeliveryStatus: i.DeliveryStatus,
			PaymentStatus:  i.PaymentStatus,
		})
	}
	return return_response, nil
}

func (o *OrderService) GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error) {
	return_response := GetAllResponse{}
	if req.RetailerId != 0 {
		resp, err := o.DB.GetByRetailerID(ctx, req.RetailerId)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		for _, i := range resp.List {
			items := []Item{}
			for _, i := range i.Items {
				items = append(items, Item{
					ProductId: i.ProductId,
					Quantity:  i.Quantity,
				})
			}
			alreadyPresent := false
			for _, j := range return_response.List {
				if j.Id == i.Id {
					alreadyPresent = true
				}
			}
			if !alreadyPresent {
				return_response.List = append(return_response.List, GetResponse{
					Id:             i.Id,
					RetailerId:     i.RetailerId,
					Items:          items,
					Status:         i.Status,
					DeliveryStatus: i.DeliveryStatus,
					PaymentStatus:  i.PaymentStatus,
				})
			}
		}
	}

	if req.Status != "" {
		resp, err := o.DB.GetByStatus(ctx, req.Status)
		if err != nil {
			switch err {
			case port.ErrSysNoRows:
				return GetAllResponse{}, ErrEmptyGetResponse
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		for _, i := range resp.List {
			items := []Item{}
			for _, i := range i.Items {
				items = append(items, Item{
					ProductId: i.ProductId,
					Quantity:  i.Quantity,
				})
			}
			alreadyPresent := false
			for _, j := range return_response.List {
				if j.Id == i.Id {
					alreadyPresent = true
				}
			}
			if !alreadyPresent {
				return_response.List = append(return_response.List, GetResponse{
					Id:             i.Id,
					RetailerId:     i.RetailerId,
					Items:          items,
					Status:         i.Status,
					DeliveryStatus: i.DeliveryStatus,
					PaymentStatus:  i.PaymentStatus,
				})
			}
		}
	}
	if len(return_response.List) == 0 {
		return return_response, ErrEmptyGetResponse
	}

	return return_response, nil
}

func (o *OrderService) UpdateStatus(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate
	resp, err := o.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return 0, ErrIdNotFound
			return 0, ErrIdNotFound
		default:
			return 0, ErrUnknown
			return 0, ErrUnknown
		}
	}

	//update
	err = o.DB.UpdateOrderStatus(ctx, &port.UpdateOrderStatusRequest{
		Id:     req.Id,
		Status: req.Status,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	//TODO: reserve stock
<<<<<<< HEAD

	return resp.Id, nil
}

func (o *OrderService) UpdatePaymentStatus(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate
	resp, err := o.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return 0, ErrIdNotFound
		default:
			return 0, ErrUnknown
		}
	}

	//update
	err = o.DB.UpdatePaymentStatus(ctx, &port.UpdateOrderPaymentStatusRequest{
		Id:            req.Id,
		PaymentStatus: req.PaymentStatus,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	// deplete stock if order status is COMPELETED

	return resp.Id, nil
}

func (o *OrderService) UpdateDeliveryStatus(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate
	resp, err := o.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return 0, ErrIdNotFound
		default:
			return 0, ErrUnknown
		}
	}

	//update
	err = o.DB.UpdateDeliveryStatus(ctx, &port.UpdateOrderDeliveryStatusRequest{
		Id:             req.Id,
		DeliveryStatus: req.DeliveryStatus,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
			return 0, ErrUnknown
		}
	}

	// deplete stock if order status is COMPELETED
=======
>>>>>>> f7d0812b (+ passing integration tests order)

	return resp.Id, nil
}

func (o *OrderService) UpdatePaymentStatus(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate
	resp, err := o.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return 0, ErrIdNotFound
		default:
			return 0, ErrUnknown
		}
	}

	//update
	err = o.DB.UpdatePaymentStatus(ctx, &port.UpdateOrderPaymentStatusRequest{
		Id:            req.Id,
		PaymentStatus: req.PaymentStatus,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	// deplete stock if order status is COMPELETED

	return resp.Id, nil
}

func (o *OrderService) UpdateDeliveryStatus(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate
	resp, err := o.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case port.ErrSysNoRows:
			return 0, ErrIdNotFound
		default:
			return 0, ErrUnknown
		}
	}

	//update
	err = o.DB.UpdateDeliveryStatus(ctx, &port.UpdateOrderDeliveryStatusRequest{
		Id:             req.Id,
		DeliveryStatus: req.DeliveryStatus,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	// deplete stock if order status is COMPELETED

	return resp.Id, nil

	return resp.Id, nil
}
