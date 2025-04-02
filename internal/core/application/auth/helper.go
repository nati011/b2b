package auth

import (
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func validateName(firstName string, lastName string) error {
	if firstName == "" {
		return ErrFirstNameNotSupplied
	}

	if lastName == "" {
		return ErrLastNameNotSupplied
	}

	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return ErrEmailNotSupplied
	}

	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

func validatePasswords(password string) error {
	if password == "" {
		return ErrPasswordNotSupplied
	}
	return nil
}
