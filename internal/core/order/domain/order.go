package domain

import (
	"encoding/json"
	"time"

	goodmoney "github.com/the-nucleus-project/good_money"
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
	OrderID   int64
	ProductID int64
	Quantity  int
	Price     *goodmoney.Money
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Order represents an order snapshot for checkout and fulfillment.
type Order struct {
	ID                      int64
	CustomerID              int64
	Status                  OrderStatus
	PaymentStatus           string
	DeliveryStatus          string
	ConfirmationStatus      string
	Total                   *goodmoney.Money
	CustomerSnapshot        json.RawMessage
	ShippingAddressSnapshot json.RawMessage
	BillingAddressSnapshot  json.RawMessage
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Items                   []OrderItem
}

// OrderMetadata holds the mutable order fields beyond customer and items.
type OrderMetadata struct {
	PaymentStatus           string
	DeliveryStatus          string
	ConfirmationStatus      string
	Total                   *goodmoney.Money
	CustomerSnapshot        json.RawMessage
	ShippingAddressSnapshot json.RawMessage
	BillingAddressSnapshot  json.RawMessage
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
		Status:                  parsedStatus,
		PaymentStatus:           metadata.PaymentStatus,
		DeliveryStatus:          metadata.DeliveryStatus,
		ConfirmationStatus:      metadata.ConfirmationStatus,
		Total:                   metadata.Total,
		CustomerSnapshot:        metadata.CustomerSnapshot,
		ShippingAddressSnapshot: metadata.ShippingAddressSnapshot,
		BillingAddressSnapshot:  metadata.BillingAddressSnapshot,
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
	o.ShippingAddressSnapshot = metadata.ShippingAddressSnapshot
	o.BillingAddressSnapshot = metadata.BillingAddressSnapshot
	o.UpdatedAt = time.Now()
	return nil
}
