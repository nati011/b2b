package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user
type UserStatus string

const (
	// UserStatusActive represents an active user
	UserStatusActive UserStatus = "active"
	// UserStatusInactive represents an inactive user
	UserStatusInactive UserStatus = "inactive"
	// UserStatusSuspended represents a suspended user
	UserStatusSuspended UserStatus = "suspended"
)

// UserStatus represents the type of a user
type UserType string

const (
	// UserStatusActive represents self service user
	UserTypeSelfService UserType = "selfservice"
	// UserTypeOfficer represents officer user
	UserTypeOfficer UserType = "officer"
)

// UserTypeFromString converts a string to UserType, defaulting to UserTypeSelfService if invalid.
func UserTypeFromString(s string) UserType {
	normalized := strings.ToLower(strings.TrimSpace(s))
	// Handle both "selfservice" and "self_service" formats
	if normalized == "selfservice" || normalized == "self_service" {
		return UserTypeSelfService
	}
	switch normalized {
	case string(UserTypeOfficer):
		return UserTypeOfficer
	default:
		return UserTypeSelfService
	}
}

// User represents a user entity
// Note: Roles are stored as IDs only to maintain domain boundaries.
// Role objects are loaded separately when needed via the service layer.
type User struct {
	ID          string
	ExternalID  string
	Email       Email
	PhoneNumber PhoneNumber
	Name        string
	Status      UserStatus
	UserType    UserType
	RoleIDs     []string // Role IDs only - no direct domain dependency
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewUser creates a new user with either email or phone number (at least one required)
func NewUser(externalID string, email Email, phoneNumber PhoneNumber, name string, userType UserType) (*User, error) {
	if ok, err := externalIDValidator(externalID); !ok {
		return nil, err
	}

	if ok, err := nameValidator(name); !ok {
		return nil, err
	}

	if ok, err := userTypeValidator(userType); !ok {
		return nil, err
	}

	// At least one identifier (email or phone number) must be provided
	if email.IsEmpty() && phoneNumber.IsEmpty() {
		return nil, errors.New("either email or phone number must be provided")
	}

	now := time.Now()

	user := &User{
		ID:          uuid.New().String(),
		ExternalID:  externalID,
		Email:       email,
		PhoneNumber: phoneNumber,
		Name:        name,
		Status:      UserStatusInactive,
		UserType:    userType,
		RoleIDs:     []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return user, nil
}

// GetID returns the user ID
func (u *User) GetID() string {
	return u.ID
}

// Update updates user information
func (u *User) Update(externalID string, email Email, name string) error {
	if ok, err := externalIDValidator(externalID); !ok {
		return err
	}
	if ok, err := nameValidator(name); !ok {
		return err
	}

	updatedAt := time.Now()

	u.ExternalID = externalID
	u.Email = email
	u.Name = name
	u.UpdatedAt = updatedAt
	return nil
}

// Activate activates the user
func (u *User) Activate() error {
	newStatus := UserStatusActive

	if ok, err := userStatusValidator(newStatus); !ok {
		return err
	}

	if u.Status == UserStatusActive {
		return nil
	}

	updatedAt := time.Now()

	u.Status = newStatus
	u.UpdatedAt = updatedAt
	return nil
}

// Deactivate deactivates the user
func (u *User) Deactivate() error {
	newStatus := UserStatusInactive

	if ok, err := userStatusValidator(newStatus); !ok {
		return err
	}

	if u.Status == UserStatusInactive {
		return nil
	}

	updatedAt := time.Now()

	u.Status = newStatus
	u.UpdatedAt = updatedAt
	return nil
}

// Suspend suspends the user
func (u *User) Suspend() error {
	newStatus := UserStatusSuspended

	if ok, err := userStatusValidator(newStatus); !ok {
		return err
	}

	updatedAt := time.Now()

	u.Status = newStatus
	u.UpdatedAt = updatedAt
	return nil
}

// AssignRole assigns a role ID to the user
func (u *User) AssignRole(roleID string) {
	for _, rID := range u.RoleIDs {
		if rID == roleID {
			return // Role already assigned
		}
	}

	updatedAt := time.Now()
	u.RoleIDs = append(u.RoleIDs, roleID)
	u.UpdatedAt = updatedAt
}

// RevokeRole revokes a role ID from the user
func (u *User) RevokeRole(roleID string) {
	for i, rID := range u.RoleIDs {
		if rID == roleID {
			updatedAt := time.Now()
			u.RoleIDs = append(u.RoleIDs[:i], u.RoleIDs[i+1:]...)
			u.UpdatedAt = updatedAt
			return
		}
	}
}

// HasRole checks if the user has a specific role ID
func (u *User) HasRole(roleID string) bool {
	for _, rID := range u.RoleIDs {
		if rID == roleID {
			return true
		}
	}
	return false
}

// CanLogin checks if the user can log in
func (u *User) CanLogin() bool {
	return u.Status == UserStatusActive
}
