package order

import (
	"context"
	"errors"
	"log"
	"time"

	"b2b.nati011.github.com/internal/core/application/checkout"
	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/product"
	"b2b.nati011.github.com/internal/core/domain/retailer"
	port_commons "b2b.nati011.github.com/internal/port/commons/db"
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
	ErrDuplicateOrderNotAllowed           = errors.New("oopsy, duplicate order not allowed")
)

// order status
const (
	CANCELED_STATUS  = "CANCELED"
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

type Item struct {
	ProductId    int
	ProductName  string
	ProductPrice float64
	Quantity     int
}

type PlaceRequest struct {
	RetailerId       int
	Items            []Item
	PaymentPartnerId int
}

type GetResponse struct {
	Id             int
	RetailerId     int
	RetailerName   string
	Items          []Item
	Total          float32
	Status         string
	DeliveryStatus string
	PaymentStatus  string
	CreatedAt      time.Time
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
}

type OrderPlaceResponse struct {
	Id          int    `json:"id"`
	CheckoutUrl string `json:"checkout_url"`
}

type Provider interface {
	Place(ctx context.Context, req *PlaceRequest) (OrderPlaceResponse, error)
	InitPayment(ctx context.Context, id int) (OrderPlaceResponse, error)
	Cancel(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
	UpdateStatus(ctx context.Context, req *UpdateRequest) (int, error)
	UpdatePaymentStatus(ctx context.Context, req *UpdateRequest) (int, error)
	UpdateDeliveryStatus(ctx context.Context, req *UpdateRequest) (int, error)
	GetDistributorOrders(ctx context.Context, id int) (GetAllResponse, error)
	GetRetailerOrders(ctx context.Context, id int) (GetAllResponse, error)
}

type OrderService struct {
	DB              port.DB
	InvoiceService  invoice.Provider
	ProductService  product.Provider
	RetailerService retailer.Provider
	CheckoutService checkout.Provider
}

func NewOrderService(
	db port.DB,
	is invoice.Provider,
	ps product.Provider,
	rs retailer.Provider,
	pays checkout.Provider,

) Provider {

	return &OrderService{
		DB:              db,
		InvoiceService:  is,
		ProductService:  ps,
		RetailerService: rs,
		CheckoutService: pays,
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

// order cannot placed if retailer has ongoing order
func (o *OrderService) checkOrderDuplicacyEligibility(ctx context.Context, retailerId int) (bool, error) {
	// check if there is an incomplete(PENDING, ...) order with the same retailer
	retailer_orders, err := o.GetRetailerOrders(ctx, retailerId)
	if err != nil {
		switch err {
		case ErrEmptyGetResponse:
		default:
			return false, ErrUnknown
		}
	}

	var isEligible bool = true
	for _, o := range retailer_orders.List {
		if o.Status == PENDING_STATUS {
			isEligible = false
			break
		}
	}
	return isEligible, nil
}

func (o *OrderService) Place(ctx context.Context, req *PlaceRequest) (OrderPlaceResponse, error) {
	err := o.validate_placement(ctx, req)
	if err != nil {
		return OrderPlaceResponse{}, err
	}

	isEligible, err := o.checkOrderDuplicacyEligibility(ctx, req.RetailerId)
	if err != nil {
		return OrderPlaceResponse{}, ErrUnknown
	}
	if !isEligible {
		return OrderPlaceResponse{}, ErrDuplicateOrderNotAllowed
	}

	items := make([]port.Item, 0, len(req.Items))
	var itemsTotal float64
	for _, i := range req.Items {
		prod_resp, err := o.ProductService.Get(ctx, i.ProductId)
		if err != nil {
			switch err {
			default:
				return OrderPlaceResponse{}, ErrUnknown
			}
		}
		items = append(items, port.Item{
			ProductId:   i.ProductId,
			Quantity:    i.Quantity,
			Price:       prod_resp.Price,
			ProductName: i.ProductName,
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
	})
	if err != nil {
		return OrderPlaceResponse{}, ErrUnknown
	}

	checkout_resp, err := o.CheckoutService.Checkout(ctx, &checkout.CheckoutRequest{
		OrderId:          order_id,
		Amount:           itemsTotal,
		PaymentPartnerId: req.PaymentPartnerId,
	})
	if err != nil {
		o.Cancel(ctx, order_id)
		return OrderPlaceResponse{}, err
	}

	// create invoice
	lineItems := make([]invoice.Item, 0, len(req.Items))
	for _, i := range req.Items {
		prod_resp, err := o.ProductService.Get(ctx, i.ProductId)
		if err != nil {
			switch err {
			default:
				return OrderPlaceResponse{}, ErrUnknown
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
		return OrderPlaceResponse{}, ErrUnknown
	}

	//reserve stock
	for _, i := range req.Items {
		if err = o.ProductService.Reserve(ctx, i.ProductId, i.Quantity); err != nil {
			o.Cancel(ctx, order_id)
			log.Printf("order placement failed due to inability to reserve stock qty: %v for productId: %v", i.Quantity, i.ProductId)
			return OrderPlaceResponse{}, ErrUnknown
		}
	}

	// o.InvoiceService.Cancel(ctx, order_id)
	// o.Cancel(ctx, order_id)
	// o.ReserveStpc

	//TODO
	/*
		send sms
	*/

	//TODO
	/*
		send email
	*/
	return OrderPlaceResponse{
		Id:          order_id,
		CheckoutUrl: checkout_resp.CheckoutUrl,
	}, nil
}

func (o *OrderService) InitPayment(ctx context.Context, id int) (OrderPlaceResponse, error) {
	resp, err := o.CheckoutService.ReinitiateCheckout(ctx, &checkout.ReinitiateCheckoutRequest{
		OrderId: id,
	})
	if err != nil {
		log.Printf("failed to init payment: %v", err)
		return OrderPlaceResponse{}, nil
	}
	return OrderPlaceResponse{
		Id:          id,
		CheckoutUrl: resp.CheckoutUrl,
	}, nil
}

func (o *OrderService) Cancel(ctx context.Context, id int) error {
	//validate
	got, err := o.DB.GetByID(ctx, id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return ErrIdNotFound
		default:
			return ErrUnknown
		}
	}

	//check if already given status
	if got.Status == CANCELED_STATUS {
		return ErrAlreadyCanceled
	}

	err = o.DB.UpdateOrderStatus(ctx, &port.UpdateOrderStatusRequest{
		Id:     id,
		Status: CANCELED_STATUS,
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
		case port_commons.ErrSysNoRows:
			return GetResponse{}, ErrIdNotFound
		default:
			return GetResponse{}, ErrUnknown
		}
	}
	items := []Item{}
	for _, i := range resp.Items {
		items = append(items, Item{
			ProductId:    i.ProductId,
			ProductName:  i.ProductName,
			ProductPrice: i.ProductPrice,
			Quantity:     i.Quantity,
		})
	}
	return GetResponse{
		Id:             resp.Id,
		RetailerId:     resp.RetailerId,
		RetailerName:   resp.RetailerName,
		Total:          float32(resp.Total),
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
		case port_commons.ErrSysNoRows:
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
				ProductId:    i.ProductId,
				ProductName:  i.ProductName,
				ProductPrice: i.Price,
				Quantity:     i.Quantity,
			})

		}
		return_response.List = append(return_response.List, GetResponse{
			Id:             i.Id,
			RetailerId:     i.RetailerId,
			RetailerName:   i.RetailerName,
			Items:          items,
			Total:          float32(i.Total),
			Status:         i.Status,
			DeliveryStatus: i.DeliveryStatus,
			PaymentStatus:  i.PaymentStatus,
			CreatedAt:      i.CreatedAt,
		})
	}
	return return_response, nil
}

func (o *OrderService) GetDistributorOrders(ctx context.Context, distributor_id int) (GetAllResponse, error) {
	resp, err := o.DB.GetAll(ctx)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
			return GetAllResponse{}, ErrEmptyGetResponse
		default:
			return GetAllResponse{}, ErrUnknown
		}
	}

	return_response := GetAllResponse{}
	for _, i := range resp.List {
		distributors := []int{}
		items := []Item{}
		for _, j := range i.Items {
			items = append(items, Item{
				ProductId:   j.ProductId,
				Quantity:    j.Quantity,
				ProductName: j.ProductName,
			})
			product, err := o.ProductService.Get(ctx, j.ProductId)
			if err != nil {
				log.Fatalf("failed to fetch product %v", j.ProductId)
				return GetAllResponse{}, ErrUnknown
			}
			distributors = append(distributors, product.DistributorId)
		}
		order_belongs_to_distributor := false
		for _, j := range distributors {
			if j == distributor_id {
				order_belongs_to_distributor = true
			}
		}
		if order_belongs_to_distributor {
			return_response.List = append(return_response.List, GetResponse{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				RetailerName:   i.RetailerName,
				Items:          items,
				Total:          float32(i.Total),
				Status:         i.Status,
				DeliveryStatus: i.DeliveryStatus,
				PaymentStatus:  i.PaymentStatus,
			})
		}
	}
	return return_response, nil
}

func (o *OrderService) GetRetailerOrders(ctx context.Context, retailer_id int) (GetAllResponse, error) {
	if err := o.validate_retailerId(ctx, retailer_id); err != nil {
		return GetAllResponse{}, err
	}

	resp, err := o.DB.GetByRetailerID(ctx, retailer_id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
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
				ProductId:    i.ProductId,
				ProductName:  i.ProductName,
				ProductPrice: i.Price,
				Quantity:     i.Quantity,
			})
		}
		return_response.List = append(return_response.List, GetResponse{
			Id:             i.Id,
			RetailerId:     i.RetailerId,
			RetailerName:   i.RetailerName,
			Items:          items,
			Total:          float32(i.Total),
			Status:         i.Status,
			DeliveryStatus: i.DeliveryStatus,
			CreatedAt:      i.CreatedAt,
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
			case port_commons.ErrSysNoRows:
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		for _, i := range resp.List {
			items := []Item{}
			for _, i := range i.Items {
				items = append(items, Item{
					ProductId:   i.ProductId,
					Quantity:    i.Quantity,
					ProductName: i.ProductName,
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
					RetailerName:   i.RetailerName,
					Items:          items,
					Total:          float32(i.Total),
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
			case port_commons.ErrSysNoRows:
				return GetAllResponse{}, ErrEmptyGetResponse
			default:
				return GetAllResponse{}, ErrUnknown
			}
		}

		for _, i := range resp.List {
			items := []Item{}
			for _, i := range i.Items {
				items = append(items, Item{
					ProductId:   i.ProductId,
					Quantity:    i.Quantity,
					ProductName: i.ProductName,
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
					RetailerName:   i.RetailerName,
					Items:          items,
					Total:          float32(i.Total),
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
		case port_commons.ErrSysNoRows:
			return 0, ErrIdNotFound
		default:
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

	switch req.Status {
	case CANCELED_STATUS:
		//free reserved stock
		for _, item := range resp.Items {
			err = o.ProductService.FreeReservation(ctx, item.ProductId, item.Quantity)
			if err != nil {
				log.Printf("failed to free reserved stock for productId: %v", item.ProductId)
			}
		}
	}

	return resp.Id, nil
}

func (o *OrderService) UpdatePaymentStatus(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate
	resp, err := o.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
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

	return resp.Id, nil
}

func (o *OrderService) UpdateDeliveryStatus(ctx context.Context, req *UpdateRequest) (int, error) {
	//validate
	resp, err := o.Get(ctx, req.Id)
	if err != nil {
		switch err {
		case port_commons.ErrSysNoRows:
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

	switch req.DeliveryStatus {
	case DELIVERY_COMPLETED_STATUS:
		//free reserved stock
		for _, item := range resp.Items {
			err = o.ProductService.FreeReservation(ctx, item.ProductId, item.Quantity)
			if err != nil {
				log.Printf("failed to free reserved stock for productId: %v", item.ProductId)
			}
		}
	}

	return resp.Id, nil
}
