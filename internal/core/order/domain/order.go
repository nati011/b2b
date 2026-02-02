package domain

import (
	"encoding/json"
	"time"
)

// OrderStatus captures the lifecycle state of an order.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

// OrderItem represents an item within an order snapshot.
type OrderItem struct {
	OrderID   int64     `json:"order_id,omitempty"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     *float64  `json:"price,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Order represents an order snapshot for checkout and fulfillment.
type Order struct {
	ID                      int64
	CustomerID              int64
	SupplierID              int64
	Status                  OrderStatus
	PaymentStatus           string
	DeliveryStatus          string
	ConfirmationStatus      string
	Total                   *float64
	CustomerSnapshot        json.RawMessage
	CartSnapshot            json.RawMessage
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Items                   []OrderItem
}

// OrderMetadata holds the mutable order fields beyond customer and items.
type OrderMetadata struct {
	PaymentStatus           string
	DeliveryStatus          string
	ConfirmationStatus      string
	Total                   *float64
	CustomerSnapshot        json.RawMessage
	CartSnapshot            json.RawMessage
}

// NewOrder creates a new order entity after validation.
func NewOrder(customerID int64, status string, metadata OrderMetadata, items []OrderItem) (*Order, error) {
	parsedStatus, err := ParseOrderStatus(status)
	if err != nil {
		return nil, err
	}
	if err := ValidateOrderInput(customerID, metadata.PaymentStatus, metadata.DeliveryStatus, metadata.ConfirmationStatus, items); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Order{
		CustomerID:              customerID,
		SupplierID:              0, // Will be set by service layer
		Status:                  parsedStatus,
		PaymentStatus:           metadata.PaymentStatus,
		DeliveryStatus:          metadata.DeliveryStatus,
		ConfirmationStatus:      metadata.ConfirmationStatus,
		Total:                   metadata.Total,
		CustomerSnapshot:        metadata.CustomerSnapshot,
		CreatedAt:               now,
		UpdatedAt:               now,
		Items:                   items,
	}, nil
}

// Update modifies order metadata fields in-place.
func (o *Order) Update(status string, metadata OrderMetadata) error {
	parsedStatus, err := ParseOrderStatus(status)
	if err != nil {
		return err
	}
	if err := ValidateOrderMetadata(metadata.PaymentStatus, metadata.DeliveryStatus, metadata.ConfirmationStatus); err != nil {
		return err
	}

	o.Status = parsedStatus
	o.PaymentStatus = metadata.PaymentStatus
	o.DeliveryStatus = metadata.DeliveryStatus
	o.ConfirmationStatus = metadata.ConfirmationStatus
	o.Total = metadata.Total
	o.CustomerSnapshot = metadata.CustomerSnapshot
	o.UpdatedAt = time.Now()
	return nil
}
