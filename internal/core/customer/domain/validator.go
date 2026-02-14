package domain

import (
	"strings"

	"marketplace/pkg/validate"
)

// ParseCustomerStatus normalizes and validates a customer status value.
func ParseCustomerStatus(value string) (CustomerStatus, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return CustomerStatusInactive, nil
	}

	validStatuses := []string{
		string(CustomerStatusActive),
		string(CustomerStatusInactive),
		string(CustomerStatusSuspended),
	}

	for _, status := range validStatuses {
		if normalized == status {
			return CustomerStatus(normalized), nil
		}
	}

	return "", validate.NewJSONError([]string{"invalid status"})
}

// ValidateCustomerInput validates core customer fields.
func ValidateCustomerInput(fullName, city, region, woreda, phoneNumber, email string) error {
	if strings.TrimSpace(email) == "" && strings.TrimSpace(phoneNumber) == "" {
		return validate.NewJSONError([]string{"either email or phone number must be provided"})
	}

	result := validate.New().
		And(validate.NonEmpty(fullName)).
		And(validate.MaxLen(fullName, 255)).
		And(validate.MaxLen(city, 100)).
		And(validate.MaxLen(region, 100)).
		And(validate.MaxLen(woreda, 100))

	if email != "" {
		result.And(validate.EmailValid(email)).
			And(validate.MaxLen(email, 255))
	}

	if phoneNumber != "" {
		result.And(validate.MinLen(phoneNumber, 10)).
			And(validate.MaxLen(phoneNumber, 20))
	}

	validation := result.Validate()
	if !validation.IsValid {
		return validate.NewJSONError(validation.Message)
	}

	return nil
}
