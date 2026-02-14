package domain

import (
	"strings"
	"testing"
)

const (
	testEmailValue            = "test@example.com"
	errMsgFailedToCreateEmail = "Failed to create email: %v"
)

// assertEmailCreationResult is a helper function to reduce cognitive complexity
func assertEmailCreationResult(t *testing.T, email Email, err error, wantErr bool, errContains string, expectedValue string) {
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

	if email.String() != expectedValue {
		t.Errorf("Expected email %s, got %s", expectedValue, email.String())
	}
}

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantErr     bool
		errContains string
	}{
		{
			name:    "successful creation with valid email",
			value:   testEmailValue,
			wantErr: false,
		},
		{
			name:    "successful creation with valid email with subdomain",
			value:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "successful creation with valid email with plus",
			value:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:        "fails with empty email",
			value:       "",
			wantErr:     true,
			errContains: "empty",
		},
		{
			name:        "fails with invalid email format",
			value:       "not-an-email",
			wantErr:     true,
			errContains: "email",
		},
		{
			name:        "fails with email missing @",
			value:       "testexample.com",
			wantErr:     true,
			errContains: "email",
		},
		{
			name:        "fails with email missing domain",
			value:       "test@",
			wantErr:     true,
			errContains: "email",
		},
		{
			name:        "fails with email missing local part",
			value:       "@example.com",
			wantErr:     true,
			errContains: "email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.value)
			assertEmailCreationResult(t, email, err, tt.wantErr, tt.errContains, tt.value)
		})
	}
}

func TestEmailString(t *testing.T) {
	email, err := NewEmail(testEmailValue)
	if err != nil {
		t.Fatalf(errMsgFailedToCreateEmail, err)
	}

	if email.String() != testEmailValue {
		t.Errorf("Expected String() to return '%s', got '%s'", testEmailValue, email.String())
	}
}

func TestEmailValue(t *testing.T) {
	email, err := NewEmail(testEmailValue)
	if err != nil {
		t.Fatalf(errMsgFailedToCreateEmail, err)
	}

	if email.Value() != testEmailValue {
		t.Errorf("Expected Value() to return '%s', got '%s'", testEmailValue, email.Value())
	}
}

func TestEmailEquals(t *testing.T) {
	email1, _ := NewEmail(testEmailValue)
	email2, _ := NewEmail(testEmailValue)
	email3, _ := NewEmail("other@example.com")

	tests := []struct {
		name     string
		email1   Email
		email2   Email
		expected bool
	}{
		{
			name:     "returns true for equal emails",
			email1:   email1,
			email2:   email2,
			expected: true,
		},
		{
			name:     "returns false for different emails",
			email1:   email1,
			email2:   email3,
			expected: false,
		},
		{
			name:     "returns true for same email instance",
			email1:   email1,
			email2:   email1,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.email1.Equals(tt.email2)
			if result != tt.expected {
				t.Errorf("Expected Equals() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEmailIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		email    Email
		expected bool
	}{
		{
			name:     "returns false for non-empty email",
			email:    mustNewEmailForEmailTest(t, testEmailValue),
			expected: false,
		},
		{
			name:     "returns true for empty email",
			email:    Email{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.email.IsEmpty()
			if result != tt.expected {
				t.Errorf("Expected IsEmpty() = %v, got %v", tt.expected, result)
			}
		})
	}
}

// Helper function
func mustNewEmailForEmailTest(t *testing.T, value string) Email {
	email, err := NewEmail(value)
	if err != nil {
		t.Fatalf(errMsgFailedToCreateEmail, err)
	}
	return email
}
