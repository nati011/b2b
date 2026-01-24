package domain

// Email represents an email address value object
type Email struct {
	value string
}

// NewEmail creates a new email value object
func NewEmail(value string) (Email, error) {
	ok, err := emailValidator(value)
	if !ok {
		return Email{}, err
	}

	return Email{value: value}, nil
}

// String returns the email as a string
func (e Email) String() string {
	return e.value
}

// Value returns the email value
func (e Email) Value() string {
	return e.value
}

// Equals checks if two emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsEmpty checks if the email is empty
func (e Email) IsEmpty() bool {
	return e.value == ""
}
