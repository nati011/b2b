package domain

import (
	fluent_validator "marketplace/pkg/validate"
)

// ValidateResourceInput validates the core resource attributes.
func ValidateResourceInput(code, service, description string) error {
	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(code)).
		And(fluent_validator.MaxLen(code, 255)).
		And(fluent_validator.NonEmpty(service)).
		And(fluent_validator.MaxLen(service, 255)).
		And(fluent_validator.MaxLen(description, 1024)).
		Validate()

	if !result.IsValid {
		return fluent_validator.NewJSONError(result.Message)
	}

	return nil
}

// ValidateActionInput validates an individual action specification.
func ValidateActionInput(name, description string) error {
	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(name)).
		And(fluent_validator.MaxLen(name, 255)).
		And(fluent_validator.MaxLen(description, 1024)).
		Validate()

	if !result.IsValid {
		return fluent_validator.NewJSONError(result.Message)
	}

	return nil
}
