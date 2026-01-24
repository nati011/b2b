package domain

import (
	validator "github.com/nucleus-proj/validate/v2"
)

// ValidateProductInput validates product fields and metadata.
func ValidateProductInput(supplierID int64, name string, metadata ProductMetadata) error {
	if supplierID <= 0 {
		return validator.NewJSONError([]string{"supplier_id must be provided"})
	}

	result := validator.New().
		And(validator.NonEmpty(name)).
		And(validator.MaxLen(name, 255)).
		And(validator.MaxLen(metadata.Description, 255)).
		And(validator.MaxLen(metadata.ExternalID, 255)).
		And(validator.MaxLen(metadata.Unit, 50))

	if metadata.Price != nil && *metadata.Price < 0 {
		return validator.NewJSONError([]string{"price must be non-negative"})
	}

	if metadata.TotalQuantity < 0 {
		return validator.NewJSONError([]string{"total_quantity must be non-negative"})
	}
	if metadata.ReservedQuantity < 0 {
		return validator.NewJSONError([]string{"reserved_quantity must be non-negative"})
	}
	if metadata.TotalQuantity > 0 && metadata.ReservedQuantity > metadata.TotalQuantity {
		return validator.NewJSONError([]string{"reserved_quantity must not exceed total_quantity"})
	}

	for _, categoryID := range metadata.CategoryIDs {
		if categoryID <= 0 {
			return validator.NewJSONError([]string{"category_id must be positive"})
		}
	}

	validation := result.Validate()
	if !validation.IsValid {
		return validator.NewJSONError(validation.Message)
	}

	return nil
}
