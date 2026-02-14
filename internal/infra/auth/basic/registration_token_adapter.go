package basic

import (
	"context"
	userDomain "marketplace/internal/infra/user/domain"
	userservice "marketplace/internal/infra/user/service"
	"errors"
)

// Registration token validation errors - auth module agnostic
var (
	ErrRegistrationTokenInvalid = errors.New("invalid registration token")
	ErrRegistrationTokenExpired = errors.New("registration token has expired")
	ErrRegistrationTokenUsed    = errors.New("registration token has already been used")
)

// RegistrationTokenAdapter adapts the user module's RegistrationTokenService
// to the auth module's RegistrationTokenValidator interface, translating
// user domain errors to auth module errors to maintain module independence.
type RegistrationTokenAdapter struct {
	service *userservice.RegistrationTokenService
}

// NewRegistrationTokenAdapter creates a new adapter that wraps a RegistrationTokenService
func NewRegistrationTokenAdapter(service *userservice.RegistrationTokenService) *RegistrationTokenAdapter {
	return &RegistrationTokenAdapter{
		service: service,
	}
}

// ValidateAndConsume validates a registration token and marks it as used.
// It translates user domain errors to auth module errors to maintain decoupling.
func (a *RegistrationTokenAdapter) ValidateAndConsume(ctx context.Context, token string) (string, error) {
	userID, err := a.service.ValidateAndConsume(ctx, token)
	if err != nil {
		// Translate user domain errors to auth module errors
		if errors.Is(err, userDomain.ErrRegistrationTokenInvalid) {
			return "", ErrRegistrationTokenInvalid
		}
		if errors.Is(err, userDomain.ErrRegistrationTokenExpired) {
			return "", ErrRegistrationTokenExpired
		}
		if errors.Is(err, userDomain.ErrRegistrationTokenUsed) {
			return "", ErrRegistrationTokenUsed
		}
		// For other errors, return as-is (e.g., database errors)
		return "", err
	}
	return userID, nil
}
