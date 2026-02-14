package domain

import (
	fluent_validator "marketplace/pkg/validate"
)

// ValidatePermissionInput validates core permission fields.
func ValidatePermissionInput(id, resourceCode, action, description string) error {
	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(id)).
		And(fluent_validator.MaxLen(id, 255)).
		And(fluent_validator.NonEmpty(resourceCode)).
		And(fluent_validator.MaxLen(resourceCode, 255)).
		And(fluent_validator.NonEmpty(action)).
		And(fluent_validator.MaxLen(action, 255)).
		And(fluent_validator.MaxLen(description, 1024)).
		Validate()

	if !result.IsValid {
		return fluent_validator.NewJSONError(result.Message)
	}

	return nil
}


