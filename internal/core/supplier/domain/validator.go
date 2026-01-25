package domain

import (
	"strings"

	"marketplace/pkg/validate"
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

	return "", validate.NewJSONError([]string{"invalid status"})
}

// ValidateSupplierInput validates core supplier fields.
func ValidateSupplierInput(businessName, supportEmail, supportPhone string) error {
	if strings.TrimSpace(supportEmail) == "" && strings.TrimSpace(supportPhone) == "" {
		return validate.NewJSONError([]string{"either support email or support phone must be provided"})
	}

	result := validate.New().
		And(validate.NonEmpty(businessName)).
		And(validate.MaxLen(businessName, 255))

	if supportEmail != "" {
		result.And(validate.EmailValid(supportEmail)).
			And(validate.MaxLen(supportEmail, 255))
	}

	if supportPhone != "" {
		result.And(validate.MinLen(supportPhone, 10)).
			And(validate.MaxLen(supportPhone, 50))
	}

	validation := result.Validate()
	if !validation.IsValid {
		return validate.NewJSONError(validation.Message)
	}

	return nil
}

