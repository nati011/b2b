package domain

import (
	"errors"
	"strings"
	"testing"
)

const (
	errMsgFailedToCreatePhoneNumber = "Failed to create phone number: %v"
)

// assertPhoneCreationResult is a helper function to reduce cognitive complexity
func assertPhoneCreationResult(t *testing.T, phone PhoneNumber, err error, wantErr bool, errContains string, expectedValue string) {
	t.Helper()
	if wantErr {
		if err == nil {
			t.Errorf("Expected error but got none")
			return
		}
		if errContains != "" && !strings.Contains(err.Error(), errContains) {
			t.Errorf("Expected error to contain '%s', got '%s'", errContains, err.Error())
		}
		return
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if phone.String() != expectedValue {
		t.Errorf("Expected phone %s, got %s", expectedValue, phone.String())
	}
}

func TestNewPhoneNumber(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		validator   PhoneValidator
		wantErr     bool
		errContains string
	}{
		{
			name:      "successful creation with valid phone number",
			value:     "1234567890",
			validator: nil,
			wantErr:   false,
		},
		{
			name:      "successful creation with phone number at min length",
			value:     "1234567890",
			validator: nil,
			wantErr:   false,
		},
		{
			name:      "successful creation with phone number at max length",
			value:     "12345678901234567890",
			validator: nil,
			wantErr:   false,
		},
		{
			name:        "fails with empty phone number",
			value:       "",
			validator:   nil,
			wantErr:     true,
			errContains: "empty",
		},
		{
			name:        "fails with phone number too short",
			value:       "123456789",
			validator:   nil,
			wantErr:     true,
			errContains: "short",
		},
		{
			name:        "fails with phone number too long",
			value:       "123456789012345678901",
			validator:   nil,
			wantErr:     true,
			errContains: "long",
		},
		{
			name:  "uses custom validator when provided",
			value: "custom-phone",
			validator: func(value string) error {
				if value == "custom-phone" {
					return nil
				}
				return errors.New("validation failed")
			},
			wantErr: false,
		},
		{
			name:  "custom validator can reject valid default format",
			value: "1234567890",
			validator: func(value string) error {
				return errors.New("custom validation failed")
			},
			wantErr:     true,
			errContains: "custom validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhoneNumber(tt.value, tt.validator)
			assertPhoneCreationResult(t, phone, err, tt.wantErr, tt.errContains, tt.value)
		})
	}
}

func TestPhoneNumberString(t *testing.T) {
	phone, err := NewPhoneNumber("1234567890", nil)
	if err != nil {
		t.Fatalf(errMsgFailedToCreatePhoneNumber, err)
	}

	if phone.String() != "1234567890" {
		t.Errorf("Expected String() to return '1234567890', got '%s'", phone.String())
	}
}

func TestPhoneNumberValue(t *testing.T) {
	phone, err := NewPhoneNumber("1234567890", nil)
	if err != nil {
		t.Fatalf(errMsgFailedToCreatePhoneNumber, err)
	}

	if phone.Value() != "1234567890" {
		t.Errorf("Expected Value() to return '1234567890', got '%s'", phone.Value())
	}
}

func TestPhoneNumberEquals(t *testing.T) {
	phone1, _ := NewPhoneNumber("1234567890", nil)
	phone2, _ := NewPhoneNumber("1234567890", nil)
	phone3, _ := NewPhoneNumber("9876543210", nil)

	tests := []struct {
		name     string
		phone1   PhoneNumber
		phone2   PhoneNumber
		expected bool
	}{
		{
			name:     "returns true for equal phone numbers",
			phone1:   phone1,
			phone2:   phone2,
			expected: true,
		},
		{
			name:     "returns false for different phone numbers",
			phone1:   phone1,
			phone2:   phone3,
			expected: false,
		},
		{
			name:     "returns true for same phone number instance",
			phone1:   phone1,
			phone2:   phone1,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.phone1.Equals(tt.phone2)
			if result != tt.expected {
				t.Errorf("Expected Equals() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPhoneNumberIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		phone    PhoneNumber
		expected bool
	}{
		{
			name:     "returns false for non-empty phone number",
			phone:    mustNewPhoneNumberForPhoneTest(t, "1234567890"),
			expected: false,
		},
		{
			name:     "returns true for empty phone number",
			phone:    PhoneNumber{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.phone.IsEmpty()
			if result != tt.expected {
				t.Errorf("Expected IsEmpty() = %v, got %v", tt.expected, result)
			}
		})
	}
}

// Helper function
func mustNewPhoneNumberForPhoneTest(t *testing.T, value string) PhoneNumber {
	phone, err := NewPhoneNumber(value, nil)
	if err != nil {
		t.Fatalf(errMsgFailedToCreatePhoneNumber, err)
	}
	return phone
}
