package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCredentialNotFound = errors.New("credential not found")
)

// Credential represents a basic auth credential domain entity
type Credential struct {
	ID        string
	Username  string // Can be username, email, or phone number
	Password  string // Hashed password
	UserID    string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCredential creates a new credential with the given parameters
func NewCredential(username, hashedPassword, userID string) (*Credential, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if hashedPassword == "" {
		return nil, errors.New("password hash is required")
	}
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	now := time.Now()
	return &Credential{
		ID:        uuid.New().String(),
		Username:  username,
		Password:  hashedPassword,
		UserID:    userID,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// IsActive checks if the credential is active
func (c *Credential) IsActive() bool {
	return c.Active
}

// UpdatePassword updates the password hash
func (c *Credential) UpdatePassword(hashedPassword string) error {
	if hashedPassword == "" {
		return errors.New("password hash is required")
	}
	c.Password = hashedPassword
	c.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the credential
func (c *Credential) Deactivate() {
	c.Active = false
	c.UpdatedAt = time.Now()
}

