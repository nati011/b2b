package order

import (
	"context"
	"errors"

	"b2b.nati011.github.com/internal/core/domain/invoice"
	port "b2b.nati011.github.com/internal/port/order"
)

var (
	ErrIdNotFound                         = errors.New("oopsy, id not found")
	ErrRetailerIdNotSupplied              = errors.New("oopsy, retailer id mandatory")
	ErrAtleastOneOrderItemNeeded          = errors.New("oopsy, order items cannot be empty")
	ErrItemMemberProductIdOrQuantityEmpty = errors.New("oopsy, either order item member productId or quantity missing")
	ErrUnknown                            = errors.New("oopsy, unknown error")
	ErrEmptyGetResponse                   = errors.New("oopsy, empty get response")
	ErrAlreadyCanceled                    = errors.New("oopsy, order already canceled")
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
}

func NewOrderService(db port.DB, is invoice.Provider) Provider {
	return &OrderService{
		DB:             db,
		InvoiceService: is,
	}
}

func (o *OrderService) Place(ctx context.Context, req *PlaceRequest) (int, error) {
	//validate
	err := o.validate_retailerId(ctx, req.RetailerId)
	if err != nil {
		return 0, err
	}
	err = o.validate_items(ctx, req.Items)
	if err != nil {
		return 0, err
	}

	//save data
	items := []port.Item{}
	for _, i := range req.Items {
		items = append(items, port.Item{
			ProductId: i.ProductId,
			Quantity:  i.Quantity,
		})
	}
	id, err := o.DB.Create(ctx, &port.CreateRequest{
		RetailerId: req.RetailerId,
		Items:      items,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
	}

	//set status to pending
	err = o.DB.UpdateOrderStatus(ctx, &port.UpdateOrderStatusRequest{
		Id:     id,
		Status: PENDING_STATUS,
	})
	if err != nil {
		switch err {
		default:
			return 0, ErrUnknown
		}
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
