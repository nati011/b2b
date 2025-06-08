package order

import (
	"context"
	"time"
)

type Item struct {
	ProductId    int
	ProductName  string
	ProductPrice float64
	Quantity     int
	Price        float64
}
type CreateRequest struct {
	RetailerId     int
	Items          []Item
	Status         string
	PaymentStatus  string
	DeliveryStatus string
	Total          float64
}

type GetResponse struct {
	Id             int
	RetailerId     int
	RetailerName   string
	Items          []Item
	Total          float64
	Status         string
	PaymentStatus  string
	DeliveryStatus string
	CreatedAt      time.Time
}

type GetAllResponse struct {
	List       []GetResponse
	TotalCount int64
}

type UpdateOrderStatusRequest struct {
	Id     int
	Status string
}

type UpdateOrderDeliveryStatusRequest struct {
	Id             int
	DeliveryStatus string
}

type UpdateOrderPaymentStatusRequest struct {
	Id            int
	PaymentStatus string
}

type Reader interface {
	GetByID(context.Context, int) (GetResponse, error)
	GetByRetailerID(context.Context, int) (GetAllResponse, error)
	GetByStatus(context.Context, string) (GetAllResponse, error)
	GetAll(context.Context) (GetAllResponse, error)
}

type Writer interface {
	Create(context.Context, *CreateRequest) (int, error)
	UpdateOrderStatus(context.Context, *UpdateOrderStatusRequest) error
	UpdatePaymentStatus(context.Context, *UpdateOrderPaymentStatusRequest) error
	UpdateDeliveryStatus(context.Context, *UpdateOrderDeliveryStatusRequest) error
}

type DB interface {
	Reader
	Writer
}
