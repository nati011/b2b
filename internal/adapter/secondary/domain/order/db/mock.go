package order

import (
	"context"

	port_commons "b2b.nati011.github.com/internal/port/commons/db"
	port "b2b.nati011.github.com/internal/port/domain/order"
)

type Item struct {
	ProductId int
	Quantity  int
}
type MockOrder struct {
	Id             int
	RetailerId     int
	Items          []Item
	Total          float64
	Status         string
	PaymentStatus  string
	DeliveryStatus string
}

type Mock struct {
	orders []MockOrder
}

func NewMock() port.DB {
	return &Mock{}
}

func (m *Mock) GetByID(ctx context.Context, id int) (port.GetResponse, error) {
	for _, i := range m.orders {
		if i.Id == id {
			items := []port.Item{}
			for _, i := range i.Items {
				items = append(items, port.Item{
					ProductId: i.ProductId,
					Quantity:  i.Quantity,
				})
			}
			return port.GetResponse{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				Items:          items,
				Status:         i.Status,
				PaymentStatus:  i.PaymentStatus,
				DeliveryStatus: i.DeliveryStatus,
			}, nil
		}
	}
	return port.GetResponse{}, port_commons.ErrSysNoRows
}

func (m *Mock) GetByRetailerID(ctx context.Context, id int) (port.GetAllResponse, error) {
	var resp []port.GetResponse
	for _, i := range m.orders {
		if i.Id == id {
			items := []port.Item{}
			for _, i := range i.Items {
				items = append(items, port.Item{
					ProductId: i.ProductId,
					Quantity:  i.Quantity,
				})
			}
			resp = append(resp, port.GetResponse{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				Items:          items,
				Status:         i.Status,
				PaymentStatus:  i.PaymentStatus,
				DeliveryStatus: i.DeliveryStatus,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetByStatus(ctx context.Context, status string) (port.GetAllResponse, error) {
	var resp []port.GetResponse
	for _, i := range m.orders {
		if i.Status == status {
			items := []port.Item{}
			for _, i := range i.Items {
				items = append(items, port.Item{
					ProductId: i.ProductId,
					Quantity:  i.Quantity,
				})
			}
			resp = append(resp, port.GetResponse{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				Items:          items,
				Status:         i.Status,
				PaymentStatus:  i.PaymentStatus,
				DeliveryStatus: i.DeliveryStatus,
			})
		}
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) GetAll(context.Context) (port.GetAllResponse, error) {
	var resp []port.GetResponse
	for _, i := range m.orders {
		items := []port.Item{}
		for _, i := range i.Items {
			items = append(items, port.Item{
				ProductId: i.ProductId,
				Quantity:  i.Quantity,
			})
		}
		resp = append(resp, port.GetResponse{
			Id:             i.Id,
			RetailerId:     i.RetailerId,
			Items:          items,
			Status:         i.Status,
			PaymentStatus:  i.PaymentStatus,
			DeliveryStatus: i.DeliveryStatus,
		})
	}
	if len(resp) == 0 {
		return port.GetAllResponse{}, port_commons.ErrSysNoRows
	}
	return port.GetAllResponse{
		List: resp,
	}, nil
}

func (m *Mock) Create(ctx context.Context, req *port.CreateRequest) (int, error) {
	newResourceId := len(m.orders) + 1
	items := []Item{}
	for _, i := range req.Items {
		items = append(items, Item{
			ProductId: i.ProductId,
			Quantity:  i.Quantity,
		})
	}
	m.orders = append(m.orders, MockOrder{
		Id:             newResourceId,
		RetailerId:     req.RetailerId,
		Items:          items,
		Status:         req.Status,
		PaymentStatus:  req.PaymentStatus,
		DeliveryStatus: req.DeliveryStatus,
	})
	return newResourceId, nil
}

func (m *Mock) UpdateOrderStatus(ctx context.Context, req *port.UpdateOrderStatusRequest) error {
	updatedResources := []MockOrder{}
	for _, i := range m.orders {
		if req.Id == i.Id {
			updatedResources = append(updatedResources, MockOrder{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				Items:          i.Items,
				Status:         req.Status,
				PaymentStatus:  i.PaymentStatus,
				DeliveryStatus: i.DeliveryStatus,
			})
		} else {
			updatedResources = append(updatedResources, MockOrder{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				Items:          i.Items,
				Status:         i.Status,
				PaymentStatus:  i.PaymentStatus,
				DeliveryStatus: i.DeliveryStatus,
			})
		}

	}
	m.orders = updatedResources
	return nil
}

func (m *Mock) UpdatePaymentStatus(ctx context.Context, req *port.UpdateOrderPaymentStatusRequest) error {
	updatedResources := []MockOrder{}
	for _, i := range m.orders {
		if req.Id == i.Id {
			updatedResources = append(updatedResources, MockOrder{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				Items:          i.Items,
				Status:         i.Status,
				PaymentStatus:  req.PaymentStatus,
				DeliveryStatus: i.DeliveryStatus,
			})
		} else {
			updatedResources = append(updatedResources, MockOrder{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				Items:          i.Items,
				Status:         i.Status,
				PaymentStatus:  i.PaymentStatus,
				DeliveryStatus: i.DeliveryStatus,
			})
		}

	}
	m.orders = updatedResources
	return nil
}

func (m *Mock) UpdateDeliveryStatus(ctx context.Context, req *port.UpdateOrderDeliveryStatusRequest) error {
	updatedResources := []MockOrder{}
	for _, i := range m.orders {
		if req.Id == i.Id {
			updatedResources = append(updatedResources, MockOrder{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				Items:          i.Items,
				Status:         i.Status,
				PaymentStatus:  i.PaymentStatus,
				DeliveryStatus: req.DeliveryStatus,
			})
		} else {
			updatedResources = append(updatedResources, MockOrder{
				Id:             i.Id,
				RetailerId:     i.RetailerId,
				Items:          i.Items,
				Status:         i.Status,
				PaymentStatus:  i.PaymentStatus,
				DeliveryStatus: i.DeliveryStatus,
			})
		}

	}
	m.orders = updatedResources
	return nil
}
