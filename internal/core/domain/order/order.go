package order

import (
	"context"
	"errors"

	"b2b.nati011.github.com/internal/core/domain/invoice"
	"b2b.nati011.github.com/internal/core/domain/product"
	transaction_control "b2b.nati011.github.com/internal/core/util/transaction_control"
	port "b2b.nati011.github.com/internal/port/domain/order"
)

var (
	ErrIdNotFound                         = errors.New("oopsy, id not found")
	ErrRetailerIdNotSupplied              = errors.New("oopsy, retailer id mandatory")
	ErrAtleastOneOrderItemNeeded          = errors.New("oopsy, order items cannot be empty")
	ErrItemMemberProductIdOrQuantityEmpty = errors.New("oopsy, either order item member productId or quantity missing")
	ErrUnknown                            = errors.New("oopsy, unknown error")
	ErrEmptyGetResponse                   = errors.New("oopsy, empty get response")
	ErrAlreadyCanceled                    = errors.New("oopsy, order already canceled")
	ErrItemMemberProductNotFound          = errors.New("oopsy, product not found")
	ErrItemMemberProductQuantityNotFound  = errors.New("oopsy, product quantity not found")
)

const (
	CANCELD_STATUS = "CANCELED"
	PENDING_STATUS = "PENDING"
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
	Id         int
	RetailerId int
	Items      []Item
	Total      float32
	Status     string
}

type GetAllResponse struct {
	List []GetResponse
}

type GetByParamRequest struct {
	RetailerId int
	Status     string
}

type Provider interface {
	Place(ctx context.Context, req *PlaceRequest) (int, error)
	Cancel(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (GetResponse, error)
	GetAll(ctx context.Context) (GetAllResponse, error)
	GetByParam(ctx context.Context, req *GetByParamRequest) (GetAllResponse, error)
}

type OrderService struct {
	DB             port.DB
	InvoiceService invoice.Provider
	ProductService product.Provider
	tm             transaction_control.TransactionManager
}

func NewOrderService(db port.DB,
	is invoice.Provider,
	ps product.Provider) Provider {

	return &OrderService{
		DB:             db,
		InvoiceService: is,
		ProductService: ps,
		tm:             transaction_control.NewTransactionManager(),
	}
}

func (o *OrderService) Place(ctx context.Context, req *PlaceRequest) (int, error) {

	id, err := o.tm.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		o.tm.Init(ctx)

		// Validate retailer ID
		if err := o.validate_retailerId(ctx, req.RetailerId); err != nil {
			return 0, err
		}

		// Validate items
		if err := o.validate_items(ctx, req.Items); err != nil {
			return 0, err
		}

		// Convert items to port format
		items := make([]port.Item, 0, len(req.Items))
		for _, i := range req.Items {
			items = append(items, port.Item{
				ProductId: i.ProductId,
				Quantity:  i.Quantity,
			})
		}

		// Create order
		id, err := o.DB.Create(ctx, &port.CreateRequest{
			RetailerId: req.RetailerId,
			Items:      items,
		})
		if err != nil {
			return 0, ErrUnknown
		}

		// Set status to pending
		if err := o.DB.UpdateOrderStatus(ctx, &port.UpdateOrderStatusRequest{
			Id:     id,
			Status: PENDING_STATUS,
		}); err != nil {
			return 0, ErrUnknown
		}

		// Create invoice
		lineItems := make([]invoice.Item, 0, len(req.Items))
		for _, i := range req.Items {
			lineItems = append(lineItems, invoice.Item{
				ProductId:       i.ProductId,
				ProductName:     "",
				ProductQuantity: i.Quantity,
				ProductPrice:    1,
			})
		}

		if _, err = o.InvoiceService.Create(ctx, &invoice.CreateRequest{
			Status:    invoice.DRAFT_STATUS,
			OrderId:   id,
			SubTotal:  0,
			LineItems: lineItems,
			TaxAmount: 0,
		}); err != nil {
			return 0, ErrUnknown
		}

		return id, nil
	})

	if err != nil {
		return 0, err
	}
	return id, nil
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
		Id:         resp.Id,
		RetailerId: resp.RetailerId,
		Items:      items,
		Status:     resp.Status,
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
			Id:         i.Id,
			RetailerId: i.RetailerId,
			Items:      items,
			Status:     i.Status,
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
					Id:         i.Id,
					RetailerId: i.RetailerId,
					Items:      items,
					Status:     i.Status,
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
					Id:         i.Id,
					RetailerId: i.RetailerId,
					Items:      items,
					Status:     i.Status,
				})
			}
		}
	}
	if len(return_response.List) == 0 {
		return return_response, ErrEmptyGetResponse
	}

	return return_response, nil
}
