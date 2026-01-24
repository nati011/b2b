package domain

import (
	fluent_validator "marketplace/pkg/validate"
)

func emailValidator(value string) (bool, error) {
	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(value)).
		And(fluent_validator.EmailValid(value)).
		Validate()

	if !result.IsValid {
		return false, fluent_validator.NewJSONError(result.Message)
	}
	return true, nil
}

// stringValidator validates a non-empty string with a maximum length of 255 characters.
// Used for externalID and name fields.
func stringValidator(value string) (bool, error) {
	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(value)).
		And(fluent_validator.MaxLen(value, 255)).
		Validate()

	if !result.IsValid {
		return false, fluent_validator.NewJSONError(result.Message)
	}

	return true, nil
}

// externalIDValidator validates an external ID if provided (optional).
// If empty, validation passes. If provided, must be max 255 characters.
func externalIDValidator(value string) (bool, error) {
	// External ID is optional - empty string is valid
	if value == "" {
		return true, nil
	}
	// If provided, validate it's not too long
	result := fluent_validator.New().
		And(fluent_validator.MaxLen(value, 255)).
		Validate()

	if !result.IsValid {
		return false, fluent_validator.NewJSONError(result.Message)
	}
	return true, nil
}

func nameValidator(value string) (bool, error) {
	return stringValidator(value)
}

func userStatusValidator(value UserStatus) (bool, error) {
	validStatuses := []string{
		string(UserStatusActive),
		string(UserStatusInactive),
		string(UserStatusSuspended),
	}

	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(string(value))).
		And(fluent_validator.ContainsString(validStatuses, string(value))).
		Validate()

	if !result.IsValid {
		return false, fluent_validator.NewJSONError(result.Message)
	}

	return true, nil
}

func userTypeValidator(value UserType) (bool, error) {
	validTypes := []string{
		string(UserTypeSelfService),
		string(UserTypeOfficer),
	}

	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(string(value))).
		And(fluent_validator.ContainsString(validTypes, string(value))).
		Validate()

	if !result.IsValid {
		return false, fluent_validator.NewJSONError(result.Message)
	}

	return true, nil
}

func phoneNumberValidator(value string) (bool, error) {
	result := fluent_validator.New().
		And(fluent_validator.NonEmpty(value)).
		And(fluent_validator.MinLen(value, 10)).
		And(fluent_validator.MaxLen(value, 20)).
		Validate()

	if !result.IsValid {
		return false, fluent_validator.NewJSONError(result.Message)
	}
	return true, nil
}
