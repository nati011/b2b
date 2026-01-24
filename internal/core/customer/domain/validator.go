package domain

import (
	"strings"

	validator "github.com/nucleus-proj/validate/v2"
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

	return "", validator.NewJSONError([]string{"invalid status"})
}

// ValidateCustomerInput validates core customer fields.
func ValidateCustomerInput(fullName, city, region, woreda, phoneNumber, email string) error {
	if strings.TrimSpace(email) == "" && strings.TrimSpace(phoneNumber) == "" {
		return validator.NewJSONError([]string{"either email or phone number must be provided"})
	}

	result := validator.New().
		And(validator.NonEmpty(fullName)).
		And(validator.MaxLen(fullName, 255)).
		And(validator.MaxLen(city, 100)).
		And(validator.MaxLen(region, 100)).
		And(validator.MaxLen(woreda, 100))

	if email != "" {
		result.And(validator.EmailValid(email)).
			And(validator.MaxLen(email, 255))
	}

	if phoneNumber != "" {
		result.And(validator.MinLen(phoneNumber, 10)).
			And(validator.MaxLen(phoneNumber, 20))
	}

	validation := result.Validate()
	if !validation.IsValid {
		return validator.NewJSONError(validation.Message)
	}

	return nil
}
