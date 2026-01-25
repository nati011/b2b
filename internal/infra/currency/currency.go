package currency

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	goodmoney "github.com/nucleus-proj/goodmoney"
	fluent_validator "github.com/nucleus-proj/validate"
)

const (
	// MaxCurrencyCodeLength is the maximum allowed length for a currency code (ISO 4217 is 3 chars, but allow some flexibility)
	MaxCurrencyCodeLength = 10
)

// ErrorCode represents error codes for currency errors
type ErrorCode string

const (
	// ErrorCodeInvalidEnumValue indicates an invalid enum value
	ErrorCodeInvalidEnumValue ErrorCode = "invalid.enum.value"
	// ErrorCodeValidationFailed indicates general validation failure
	ErrorCodeValidationFailed ErrorCode = "validation.failed"
	// ErrorCodeCurrencyNotFound indicates that a currency was not found
	ErrorCodeCurrencyNotFound ErrorCode = "currency.not.found"
	// ErrorCodeCurrencyAlreadyExists indicates that a currency already exists
	ErrorCodeCurrencyAlreadyExists ErrorCode = "currency.already.exists"
	// ErrorCodeCurrencyCodeConflict indicates a duplicate currency code
	ErrorCodeCurrencyCodeConflict ErrorCode = "currency.code.conflict"
	// ErrorCodeCurrencyNotSupported indicates that a currency is not supported by goodmoney
	ErrorCodeCurrencyNotSupported ErrorCode = "currency.not.supported"
	// ErrorCodeCurrencyAlreadyConfigured indicates that a currency has already been configured
	ErrorCodeCurrencyAlreadyConfigured ErrorCode = "currency.already.configured"
)

// CurrencyError represents a currency error
type CurrencyError struct {
	Code    ErrorCode
	Message string
	Details map[string]interface{}
}

func (e *CurrencyError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// NewCurrencyError creates a new currency error
func NewCurrencyError(code ErrorCode, message string) *CurrencyError {
	return &CurrencyError{
		Code:    code,
		Message: message,
		Details: nil,
	}
}

// NewCurrencyErrorWithDetails creates a new currency error with details
func NewCurrencyErrorWithDetails(code ErrorCode, message string, details map[string]interface{}) *CurrencyError {
	return &CurrencyError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

var (
	// ErrCurrencyNotFound indicates that a currency was not found
	ErrCurrencyNotFound = errors.New("currency not found")

	// ErrCurrencyCodeRequired indicates that currency code is required
	ErrCurrencyCodeRequired = NewCurrencyError(
		ErrorCodeValidationFailed,
		"currency code is required",
	)

	// ErrCurrencyCodeTooLong indicates that currency code exceeds maximum length
	ErrCurrencyCodeTooLong = NewCurrencyError(
		ErrorCodeValidationFailed,
		fmt.Sprintf("currency code cannot exceed %d characters", MaxCurrencyCodeLength),
	)

	// ErrCurrencyCodeConflict indicates a duplicate currency code
	ErrCurrencyCodeConflict = NewCurrencyError(
		ErrorCodeCurrencyCodeConflict,
		"currency code already exists",
	)

	// ErrCurrencyNotSupported indicates that a currency is not supported by goodmoney
	ErrCurrencyNotSupported = NewCurrencyError(
		ErrorCodeCurrencyNotSupported,
		"currency code is not supported by goodmoney library",
	)

	// ErrCurrencyAlreadyConfigured indicates that a currency has already been configured and cannot be changed
	ErrCurrencyAlreadyConfigured = NewCurrencyError(
		ErrorCodeCurrencyAlreadyConfigured,
		"currency has already been configured during registration and cannot be changed",
	)
)

// NewInvalidEnumError creates an error for invalid enum values
func NewInvalidEnumError(enumType, value string) *CurrencyError {
	return NewCurrencyErrorWithDetails(
		ErrorCodeInvalidEnumValue,
		fmt.Sprintf("Invalid %s value: %s", enumType, value),
		map[string]interface{}{
			"enum_type": enumType,
			"value":     value,
		},
	)
}

// Currency represents a currency entity
type Currency struct {
	// Base fields
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time

	// Currency identification
	Code string

	// Audit fields
	CreatedBy string
	UpdatedBy string
}

// NewCurrencyParams contains parameters for creating a currency
type NewCurrencyParams struct {
	Code string
}

// NewCurrency creates a new currency with validation
func NewCurrency(params NewCurrencyParams) (*Currency, error) {
	// Normalize currency code (trim and uppercase)
	params.Code = strings.TrimSpace(strings.ToUpper(params.Code))

	// Validate currency code
	if params.Code == "" {
		return nil, ErrCurrencyCodeRequired
	}

	if len(params.Code) > MaxCurrencyCodeLength {
		return nil, ErrCurrencyCodeTooLong
	}

	// Validate currency code is supported by goodmoney
	if !goodmoney.ValidateCurrency(params.Code) {
		return nil, ErrCurrencyNotSupported
	}

	now := time.Now()

	currency := &Currency{
		ID:        uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
		Code:      params.Code,
	}

	return currency, nil
}

// GetID returns the currency ID
func (c *Currency) GetID() string {
	return c.ID
}

// IsValid performs basic validation on the currency
func (c *Currency) IsValid() error {
	validator := fluent_validator.New()

	if c.Code == "" {
		validator = validator.And(fluent_validator.ValidatorFunc(func() fluent_validator.ValidationResult {
			return fluent_validator.ValidationResult{
				IsValid: false,
				Message: []string{"currency code is required"},
			}
		}))
	}

	if len(c.Code) > MaxCurrencyCodeLength {
		validator = validator.And(fluent_validator.ValidatorFunc(func() fluent_validator.ValidationResult {
			return fluent_validator.ValidationResult{
				IsValid: false,
				Message: []string{fmt.Sprintf("currency code cannot exceed %d characters", MaxCurrencyCodeLength)},
			}
		}))
	}

	result := validator.Validate()
	if !result.IsValid {
		return fluent_validator.NewJSONError(result.Message)
	}

	return nil
}
