package domain

import (
	"time"
)

// Role represents a role entity
type Role struct {
	ID            string
	Name          string
	Description   string
	PermissionIDs []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewRole creates a new role
func NewRole(id, name, description string, permissionIDs []string) (*Role, error) {
	if err := ValidateRoleInput(id, name, description); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Role{
		ID:            id,
		Name:          name,
		Description:   description,
		PermissionIDs: permissionIDs,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// UpdatePermissions updates the role's permission IDs
func (r *Role) UpdatePermissions(permissionIDs []string) {
	r.PermissionIDs = permissionIDs
	r.UpdatedAt = time.Now()
}

// UpdateDescription updates the role's description
func (r *Role) UpdateDescription(description string) error {
	if err := ValidateRoleInput(r.ID, r.Name, description); err != nil {
		return err
	}

	r.Description = description
	r.UpdatedAt = time.Now()
	return nil
}

// Update applies the provided properties to the role.
func (r *Role) Update(name, description string, permissionIDs []string) error {
	if err := ValidateRoleInput(r.ID, name, description); err != nil {
		return err
	}

	r.Name = name
	r.Description = description
	r.PermissionIDs = permissionIDs
	r.UpdatedAt = time.Now()
	return nil
}
