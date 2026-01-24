package domain

import (
	"strings"
	"testing"
)

const (
	errMsgExpectedValidTrue    = "Expected valid=true, got valid=false, error: %v"
	errMsgExpectedNoError      = "Expected no error, got: %v"
	errMsgExpectedValidFalse   = "Expected valid=false, got valid=true"
	errMsgExpectedErrorContain = "Expected error to contain '%s', got '%s'"
)

// assertValidatorResult is a helper function to reduce cognitive complexity in test functions
func assertValidatorResult(t *testing.T, valid bool, err error, wantValid bool, errContains string) {
	t.Helper()
	if wantValid {
		if !valid {
			t.Errorf(errMsgExpectedValidTrue, err)
		}
		if err != nil {
			t.Errorf(errMsgExpectedNoError, err)
		}
	} else {
		if valid {
			t.Errorf(errMsgExpectedValidFalse)
		}
		if err == nil {
			t.Errorf(errMsgExpectedErrorNone)
		}
		if errContains != "" && err != nil && !strings.Contains(err.Error(), errContains) {
			t.Errorf(errMsgExpectedErrorContain, errContains, err.Error())
		}
	}
}

func TestEmailValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantValid   bool
		errContains string
	}{
		{
			name:      "valid email",
			value:     "test@example.com",
			wantValid: true,
		},
		{
			name:      "valid email with subdomain",
			value:     "user@mail.example.com",
			wantValid: true,
		},
		{
			name:        "empty email",
			value:       "",
			wantValid:   false,
			errContains: "empty",
		},
		{
			name:        "invalid email format",
			value:       "not-an-email",
			wantValid:   false,
			errContains: "email",
		},
		{
			name:        "email missing @",
			value:       "testexample.com",
			wantValid:   false,
			errContains: "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := emailValidator(tt.value)
			assertValidatorResult(t, valid, err, tt.wantValid, tt.errContains)
		})
	}
}

func TestStringValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantValid   bool
		errContains string
	}{
		{
			name:      "valid string",
			value:     "test string",
			wantValid: true,
		},
		{
			name:      "valid string at max length",
			value:     strings.Repeat("a", 255),
			wantValid: true,
		},
		{
			name:        "empty string",
			value:       "",
			wantValid:   false,
			errContains: "",
		},
		{
			name:        "string too long",
			value:       strings.Repeat("a", 256),
			wantValid:   false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := stringValidator(tt.value)
			assertValidatorResult(t, valid, err, tt.wantValid, tt.errContains)
		})
	}
}

func TestExternalIDValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantValid   bool
		errContains string
	}{
		{
			name:      "empty external ID is valid",
			value:     "",
			wantValid: true,
		},
		{
			name:      "valid external ID",
			value:     "ext-123",
			wantValid: true,
		},
		{
			name:      "valid external ID at max length",
			value:     strings.Repeat("a", 255),
			wantValid: true,
		},
		{
			name:        "external ID too long",
			value:       strings.Repeat("a", 256),
			wantValid:   false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := externalIDValidator(tt.value)
			assertValidatorResult(t, valid, err, tt.wantValid, tt.errContains)
		})
	}
}

func TestNameValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantValid   bool
		errContains string
	}{
		{
			name:      "valid name",
			value:     "John Doe",
			wantValid: true,
		},
		{
			name:        "empty name",
			value:       "",
			wantValid:   false,
			errContains: "",
		},
		{
			name:        "name too long",
			value:       strings.Repeat("a", 256),
			wantValid:   false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := nameValidator(tt.value)
			assertValidatorResult(t, valid, err, tt.wantValid, tt.errContains)
		})
	}
}

func TestUserStatusValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       UserStatus
		wantValid   bool
		errContains string
	}{
		{
			name:      "valid active status",
			value:     UserStatusActive,
			wantValid: true,
		},
		{
			name:      "valid inactive status",
			value:     UserStatusInactive,
			wantValid: true,
		},
		{
			name:      "valid suspended status",
			value:     UserStatusSuspended,
			wantValid: true,
		},
		{
			name:        "invalid status",
			value:       UserStatus("invalid"),
			wantValid:   false,
			errContains: "",
		},
		{
			name:        "empty status",
			value:       UserStatus(""),
			wantValid:   false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := userStatusValidator(tt.value)
			assertValidatorResult(t, valid, err, tt.wantValid, tt.errContains)
		})
	}
}

func TestUserTypeValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       UserType
		wantValid   bool
		errContains string
	}{
		{
			name:      "valid selfservice type",
			value:     UserTypeSelfService,
			wantValid: true,
		},
		{
			name:      "valid officer type",
			value:     UserTypeOfficer,
			wantValid: true,
		},
		{
			name:        "invalid type",
			value:       UserType("invalid"),
			wantValid:   false,
			errContains: "",
		},
		{
			name:        "empty type",
			value:       UserType(""),
			wantValid:   false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := userTypeValidator(tt.value)
			assertValidatorResult(t, valid, err, tt.wantValid, tt.errContains)
		})
	}
}

func TestPhoneNumberValidator(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantValid   bool
		errContains string
	}{
		{
			name:      "valid phone number at min length",
			value:     "1234567890",
			wantValid: true,
		},
		{
			name:      "valid phone number at max length",
			value:     "12345678901234567890",
			wantValid: true,
		},
		{
			name:      "valid phone number in between",
			value:     "1234567890123",
			wantValid: true,
		},
		{
			name:        "empty phone number",
			value:       "",
			wantValid:   false,
			errContains: "",
		},
		{
			name:        "phone number too short",
			value:       "123456789",
			wantValid:   false,
			errContains: "",
		},
		{
			name:        "phone number too long",
			value:       "123456789012345678901",
			wantValid:   false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := phoneNumberValidator(tt.value)
			assertValidatorResult(t, valid, err, tt.wantValid, tt.errContains)
		})
	}
}
