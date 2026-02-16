package domain

import (
	"strings"

	"marketplace/pkg/validate"
)

// ParseOrderStatus validates and normalizes order status values.
func ParseOrderStatus(value string) (OrderStatus, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return OrderStatusPending, nil
	}

	validStatuses := []string{
		string(OrderStatusPending),
		string(OrderStatusConfirmed),
		string(OrderStatusShipped),
		string(OrderStatusDelivered),
		string(OrderStatusCancelled),
	}
	for _, status := range validStatuses {
		if normalized == status {
			return OrderStatus(normalized), nil
		}
	}

	return "", validate.NewJSONError([]string{"invalid order status"})
}

// ValidateOrderInput validates order metadata and items.
func ValidateOrderInput(customerID int64, paymentStatus, deliveryStatus, confirmationStatus, deliveryAddress string, items []OrderItem) error {
	if customerID <= 0 {
		return validate.NewJSONError([]string{"customer_id must be provided"})
	}
	if strings.TrimSpace(deliveryAddress) == "" {
		return validate.NewJSONError([]string{"delivery_address is required"})
	}
	if len(items) == 0 {
		return validate.NewJSONError([]string{"at least one order item is required"})
	}
	if err := ValidateOrderMetadata(paymentStatus, deliveryStatus, confirmationStatus); err != nil {
		return err
	}
	if res := validate.New().And(validate.MaxLen(deliveryAddress, 500)).Validate(); !res.IsValid {
		return validate.NewJSONError(res.Message)
	}
	for _, item := range items {
		if item.ProductID <= 0 {
			return validate.NewJSONError([]string{"order item product_id must be provided"})
		}
		if item.Quantity <= 0 {
			return validate.NewJSONError([]string{"order item quantity must be greater than zero"})
		}
	}
	return nil
}

// ValidateOrderMetadata validates order status fields.
func ValidateOrderMetadata(paymentStatus, deliveryStatus, confirmationStatus string) error {
	result := validate.New().
		And(validate.MaxLen(paymentStatus, 255)).
		And(validate.MaxLen(deliveryStatus, 255)).
		And(validate.MaxLen(confirmationStatus, 255))

	validation := result.Validate()
	if !validation.IsValid {
		return validate.NewJSONError(validation.Message)
	}
	return nil
}
