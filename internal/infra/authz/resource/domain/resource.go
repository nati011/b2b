package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Resource represents a permissionable surface (service + code) along with its allowed actions.
type Resource struct {
	ID           string
	Code         string
	Service      string
	Description  string
	DeprecatedAt *time.Time
	Actions      []Action
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Action describes an allowed operation for a resource.
type Action struct {
	ID           string
	ResourceID   string
	Name         string
	Description  string
	DeprecatedAt *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ManifestResource provides a declarative specification used during bootstrap.
type ManifestResource struct {
	Code        string
	Service     string
	Description string
	Actions     []ManifestAction
}

// ManifestAction is the manifest representation for allowed actions.
type ManifestAction struct {
	Name        string
	Description string
}

// NewResource constructs a new Resource aggregate.
func NewResource(code, service, description string) (*Resource, error) {
	if err := ValidateResourceInput(code, service, description); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Resource{
		ID:          uuid.NewString(),
		Code:        normalizeCode(code),
		Service:     strings.TrimSpace(service),
		Description: strings.TrimSpace(description),
		Actions:     []Action{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update mutates the resource metadata and bumps UpdatedAt.
func (r *Resource) Update(service, description string) error {
	if err := ValidateResourceInput(r.Code, service, description); err != nil {
		return err
	}

	r.Service = strings.TrimSpace(service)
	r.Description = strings.TrimSpace(description)
	r.UpdatedAt = time.Now()
	return nil
}

// IsDeprecated indicates whether the resource is deprecated.
func (r *Resource) IsDeprecated() bool {
	return r.DeprecatedAt != nil
}

// AttachActions replaces the in-memory action list.
func (r *Resource) AttachActions(actions []Action) {
	r.Actions = actions
}

// NewAction constructs a new action for the provided resource.
func NewAction(resourceID, name, description string) (Action, error) {
	if err := ValidateActionInput(name, description); err != nil {
		return Action{}, err
	}

	now := time.Now()
	return Action{
		ID:          uuid.NewString(),
		ResourceID:  resourceID,
		Name:        normalizeAction(name),
		Description: strings.TrimSpace(description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// IsDeprecated indicates whether the action is deprecated.
func (a Action) IsDeprecated() bool {
	return a.DeprecatedAt != nil
}

// Normalize helpers ---------------------------------------------------------
func normalizeCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func normalizeAction(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// FindAction returns an action by name if present.
func (r Resource) FindAction(name string) (*Action, bool) {
	target := normalizeAction(name)
	for i := range r.Actions {
		if r.Actions[i].Name == target {
			return &r.Actions[i], true
		}
	}
	return nil, false
}

// NormalizeCode exposes the normalization helper for consumers.
func NormalizeCode(code string) string {
	return normalizeCode(code)
}

// NormalizeAction exposes the normalization helper for consumers.
func NormalizeAction(name string) string {
	return normalizeAction(name)
}
