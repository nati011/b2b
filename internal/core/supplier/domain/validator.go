package domain

import (
	"strings"

	validator "github.com/nucleus-proj/validate/v2"
)

// ParseSupplierStatus normalizes and validates a supplier status value.
func ParseSupplierStatus(value string) (SupplierStatus, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return SupplierStatusInactive, nil
	}

	validStatuses := []string{
		string(SupplierStatusActive),
		string(SupplierStatusInactive),
		string(SupplierStatusSuspended),
	}

	for _, status := range validStatuses {
		if normalized == status {
			return SupplierStatus(normalized), nil
		}
	}

	return "", validator.NewJSONError([]string{"invalid status"})
}

// ValidateSupplierInput validates core supplier fields.
func ValidateSupplierInput(businessName, supportEmail, supportPhone string) error {
	if strings.TrimSpace(supportEmail) == "" && strings.TrimSpace(supportPhone) == "" {
		return validator.NewJSONError([]string{"either support email or support phone must be provided"})
	}

	result := validator.New().
		And(validator.NonEmpty(businessName)).
		And(validator.MaxLen(businessName, 255))

	if supportEmail != "" {
		result.And(validator.EmailValid(supportEmail)).
			And(validator.MaxLen(supportEmail, 255))
	}

	if supportPhone != "" {
		result.And(validator.MinLen(supportPhone, 10)).
			And(validator.MaxLen(supportPhone, 50))
	}

	validation := result.Validate()
	if !validation.IsValid {
		return validator.NewJSONError(validation.Message)
	}

	return nil
}

