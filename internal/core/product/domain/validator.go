package domain

import (
	"marketplace/pkg/validate"
)

// ValidateProductInput validates product fields and metadata.
func ValidateProductInput(supplierID int64, name string, metadata ProductMetadata) error {
	if supplierID <= 0 {
		return validate.NewJSONError([]string{"supplier_id must be provided"})
	}

	result := validate.New().
		And(validate.NonEmpty(name)).
		And(validate.MaxLen(name, 255)).
		And(validate.MaxLen(metadata.Description, 255)).
		And(validate.MaxLen(metadata.ExternalID, 255)).
		And(validate.MaxLen(metadata.Unit, 50))

	if metadata.Price != nil && metadata.Price.IsNegative() {
		return validate.NewJSONError([]string{"price must be non-negative"})
	}

	if metadata.TotalQuantity < 0 {
		return validate.NewJSONError([]string{"total_quantity must be non-negative"})
	}
	if metadata.ReservedQuantity < 0 {
		return validate.NewJSONError([]string{"reserved_quantity must be non-negative"})
	}
	if metadata.TotalQuantity > 0 && metadata.ReservedQuantity > metadata.TotalQuantity {
		return validate.NewJSONError([]string{"reserved_quantity must not exceed total_quantity"})
	}

	for _, categoryID := range metadata.CategoryIDs {
		if categoryID <= 0 {
			return validate.NewJSONError([]string{"category_id must be positive"})
		}
	}

	validation := result.Validate()
	if !validation.IsValid {
		return validate.NewJSONError(validation.Message)
	}

	return nil
}
