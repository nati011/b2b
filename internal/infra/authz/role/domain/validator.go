package domain

import (
	fluent_validator "marketplace/pkg/validate"
)

// ValidateRoleInput validates the core role fields.
func ValidateRoleInput(id, name, description string) error {
	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(id)).
		And(fluent_validator.MaxLen(id, 255)).
		And(fluent_validator.NonEmpty(name)).
		And(fluent_validator.MaxLen(name, 255)).
		And(fluent_validator.MaxLen(description, 1024)).
		Validate()

	if !result.IsValid {
		return fluent_validator.NewJSONError(result.Message)
	}

	return nil
}


