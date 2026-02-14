package domain

import (
	"errors"
	"time"
)

var (
	ErrRegistrationTokenExpired = errors.New("registration token has expired")
	ErrRegistrationTokenUsed    = errors.New("registration token has already been used")
	ErrRegistrationTokenInvalid = errors.New("invalid registration token")
)

// RegistrationToken represents a one-time, time-limited token for credential creation
type RegistrationToken struct {
	ID        string
	UserID    string
	Token     string    // Secure random token
	Used      bool      // Whether the token has been used
	ExpiresAt time.Time // Token expiration time
	CreatedAt time.Time
}

// IsValid checks if the token is valid (not used and not expired)
func (t *RegistrationToken) IsValid() bool {
	if t.Used {
		return false
	}
	if time.Now().After(t.ExpiresAt) {
		return false
	}
	return true
}

// MarkAsUsed marks the token as used
func (t *RegistrationToken) MarkAsUsed() {
	t.Used = true
}
