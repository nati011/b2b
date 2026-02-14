package service

import "context"

// PhoneValidationScheme represents a phone number validation scheme
type PhoneValidationScheme struct {
	ID           string
	Code         string
	Name         string
	Description  string
	RegexPattern string
	IsDefault    bool
	Active       bool
}

// PhoneValidationRepository defines the interface for phone validation scheme operations
type PhoneValidationRepository interface {
	FindByCode(ctx context.Context, code string) (*PhoneValidationScheme, error)
	FindDefault(ctx context.Context) (*PhoneValidationScheme, error)
}
