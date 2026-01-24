package domain

import "time"

// ResourceRef carries lightweight resource metadata for permission scoping.
type ResourceRef struct {
	ID          string
	Code        string
	Service     string
	Description string
}

// Permission represents a permission entity/value object.
// For role assignments we care about the resource metadata + action tuple.
type Permission struct {
	ID          string
	Resource    ResourceRef
	Action      string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPermission creates a lightweight permission value object used by roles.
func NewPermission(resource ResourceRef, action string) Permission {
	return Permission{
		Resource: resource,
		Action:   action,
	}
}

// NewPermissionEntity constructs a persisted permission with metadata.
func NewPermissionEntity(id string, resource ResourceRef, action, description string) (*Permission, error) {
	if err := ValidatePermissionInput(id, resource.Code, action, description); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Permission{
		ID:          id,
		Resource:    resource,
		Action:      action,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update modifies the permission fields in-place.
func (p *Permission) Update(resource ResourceRef, action, description string) error {
	if err := ValidatePermissionInput(p.ID, resource.Code, action, description); err != nil {
		return err
	}

	p.Resource = resource
	p.Action = action
	p.Description = description
	p.UpdatedAt = time.Now()
	return nil
}

// Equals checks if two permissions are equal (resource/action match).
func (p Permission) Equals(other Permission) bool {
	return p.Resource.Code == other.Resource.Code && p.Action == other.Action
}

// String returns a string representation of the permission.
func (p Permission) String() string {
	return p.Resource.Code + ":" + p.Action
}
