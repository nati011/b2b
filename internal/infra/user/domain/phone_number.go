package domain

// PhoneValidator is a function type for validating phone numbers
type PhoneValidator func(value string) error

// PhoneNumber represents a phone number value object
type PhoneNumber struct {
	value string
}

// NewPhoneNumber creates a new phone number value object
// If validator is nil, uses the default validator (basic length check)
func NewPhoneNumber(value string, validator PhoneValidator) (PhoneNumber, error) {
	if validator != nil {
		if err := validator(value); err != nil {
			return PhoneNumber{}, err
		}
	} else {
		// Fallback to basic validation if no validator provided
		if ok, err := phoneNumberValidator(value); !ok {
			return PhoneNumber{}, err
		}
	}

	return PhoneNumber{value: value}, nil
}

// String returns the phone number as a string
func (p PhoneNumber) String() string {
	return p.value
}

// Value returns the phone number value
func (p PhoneNumber) Value() string {
	return p.value
}

// Equals checks if two phone numbers are equal
func (p PhoneNumber) Equals(other PhoneNumber) bool {
	return p.value == other.value
}

// IsEmpty checks if the phone number is empty
func (p PhoneNumber) IsEmpty() bool {
	return p.value == ""
}
