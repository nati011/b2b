package domain

import (
	"strings"

	validator "github.com/nucleus-proj/validate/v2"
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

	return "", validator.NewJSONError([]string{"invalid order status"})
}

// ValidateOrderInput validates order metadata and items.
func ValidateOrderInput(customerID int64, paymentStatus, deliveryStatus, confirmationStatus, currency string, items []OrderItem) error {
	if customerID <= 0 {
		return validator.NewJSONError([]string{"customer_id must be provided"})
	}
	if len(items) == 0 {
		return validator.NewJSONError([]string{"at least one order item is required"})
	}
	if err := ValidateOrderMetadata(paymentStatus, deliveryStatus, confirmationStatus, currency); err != nil {
		return err
	}
	for _, item := range items {
		if item.ProductID <= 0 {
			return validator.NewJSONError([]string{"order item product_id must be provided"})
		}
		if item.Quantity <= 0 {
			return validator.NewJSONError([]string{"order item quantity must be greater than zero"})
		}
	}
	return nil
}

// ValidateOrderMetadata validates order status fields and currency.
func ValidateOrderMetadata(paymentStatus, deliveryStatus, confirmationStatus, currency string) error {
	result := validator.New().
		And(validator.MaxLen(paymentStatus, 255)).
		And(validator.MaxLen(deliveryStatus, 255)).
		And(validator.MaxLen(confirmationStatus, 255)).
		And(validator.MaxLen(currency, 10))

	validation := result.Validate()
	if !validation.IsValid {
		return validator.NewJSONError(validation.Message)
	}
	return nil
}
