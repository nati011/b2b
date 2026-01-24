package resource

import (
	"context"
	"marketplace/internal/infra/authz/resource/domain"
	"marketplace/pkg/logger"
	"database/sql"
	"errors"
	"time"
)

type repository interface {
	Create(ctx context.Context, resource *domain.Resource) error
	Update(ctx context.Context, resource *domain.Resource) error
	Deprecate(ctx context.Context, id string, deprecatedAt *time.Time) error
	FindByCode(ctx context.Context, code string) (*domain.Resource, error)
	List(ctx context.Context) ([]*domain.Resource, error)
}

var (
	// ErrResourceNotFound indicates lookup failure for a resource code.
	ErrResourceNotFound = errors.New("resource not found")
	// ErrResourceAlreadyExists indicates duplicate resource code attempts.
	ErrResourceAlreadyExists = errors.New("resource already exists")
	// ErrResourceDeprecated indicates attempts to use deprecated resources.
	ErrResourceDeprecated = errors.New("resource is deprecated")
	// ErrActionNotFound indicates the requested action is not registered.
	ErrActionNotFound = errors.New("resource action not found")
	// ErrActionDeprecated indicates the action was deprecated.
	ErrActionDeprecated = errors.New("resource action is deprecated")
)

// ActionSpec represents an action mutation payload.
type ActionSpec struct {
	Name        string
	Description string
}

// Service orchestrates resource catalog operations.
type Service struct {
	r repository
}

func NewService(repository repository) *Service {
	return &Service{r: repository}
}

// Register creates a new resource with the provided metadata/actions.
func (s *Service) Register(ctx context.Context, code, serviceName, description string, actions []ActionSpec) (*domain.Resource, error) {
	_, err := s.r.FindByCode(ctx, code)
	if err == nil {
		logger.Warn("Resource registration failed: resource already exists", "code", code, "service", serviceName)
		return nil, ErrResourceAlreadyExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		logger.Error("Resource registration failed: repository lookup error", "code", code, "service", serviceName, "error", err)
		return nil, err
	}

	resource, err := domain.NewResource(code, serviceName, description)
	if err != nil {
		logger.Warn("Resource registration failed: validation error", "code", code, "service", serviceName, "error", err)
		return nil, err
	}
	resource.AttachActions(s.toDomainActions(resource.ID, actions))

	if err := s.r.Create(ctx, resource); err != nil {
		logger.Error("Resource registration failed: repository error", "code", code, "service", serviceName, "error", err)
		return nil, err
	}

	logger.Info("Resource registered successfully", "code", code, "service", serviceName, "resource_id", resource.ID)
	return resource, nil
}

// Update modifies a resource and ensures provided actions exist.
func (s *Service) Update(ctx context.Context, code, serviceName, description string, actions []ActionSpec) (*domain.Resource, error) {
	resource, err := s.r.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("Resource update failed: resource not found", "code", code)
			return nil, ErrResourceNotFound
		}
		logger.Error("Resource update failed: repository lookup error", "code", code, "error", err)
		return nil, err
	}

	if err := resource.Update(serviceName, description); err != nil {
		logger.Warn("Resource update failed: validation error", "code", code, "service", serviceName, "error", err)
		return nil, err
	}
	if len(actions) > 0 {
		resource.AttachActions(s.toDomainActions(resource.ID, actions))
	}

	if err := s.r.Update(ctx, resource); err != nil {
		logger.Error("Resource update failed: repository error", "code", code, "service", serviceName, "error", err)
		return nil, err
	}

	logger.Info("Resource updated successfully", "code", code, "service", serviceName, "resource_id", resource.ID)
	return resource, nil
}

// Deprecate marks a resource as deprecated.
func (s *Service) Deprecate(ctx context.Context, code string) error {
	resource, err := s.r.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("Resource deprecation failed: resource not found", "code", code)
			return ErrResourceNotFound
		}
		logger.Error("Resource deprecation failed: repository lookup error", "code", code, "error", err)
		return err
	}

	now := time.Now()
	if err := s.r.Deprecate(ctx, resource.ID, &now); err != nil {
		logger.Error("Resource deprecation failed: repository error", "code", code, "resource_id", resource.ID, "error", err)
		return err
	}

	logger.Info("Resource deprecated successfully", "code", code, "resource_id", resource.ID)
	return nil
}

// List returns all resources and actions.
func (s *Service) List(ctx context.Context) ([]*domain.Resource, error) {
	resources, err := s.r.List(ctx)
	if err != nil {
		logger.Error("Resource list failed: repository error", "error", err)
		return nil, err
	}
	logger.Debug("Resource list completed", "count", len(resources))
	return resources, nil
}

// Get fetches a single resource by code.
func (s *Service) Get(ctx context.Context, code string) (*domain.Resource, error) {
	resource, err := s.r.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	return resource, nil
}

// ResolveAction ensures the (resource, action) tuple is valid and not deprecated.
func (s *Service) ResolveAction(ctx context.Context, code, action string) (*domain.Resource, *domain.Action, error) {
	resource, err := s.r.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrResourceNotFound
		}
		return nil, nil, err
	}
	if resource.IsDeprecated() {
		logger.Warn("Action resolution failed: resource is deprecated", "code", code, "action", action)
		return nil, nil, ErrResourceDeprecated
	}

	act, ok := resource.FindAction(action)
	if !ok {
		logger.Warn("Action resolution failed: action not found", "code", code, "action", action)
		return nil, nil, ErrActionNotFound
	}
	if act.IsDeprecated() {
		logger.Warn("Action resolution failed: action is deprecated", "code", code, "action", action)
		return nil, nil, ErrActionDeprecated
	}

	return resource, act, nil
}

// Bootstrap ensures resources declared in the manifest exist.
func (s *Service) Bootstrap(ctx context.Context, manifest []domain.ManifestResource) error {
	for _, spec := range manifest {
		if err := s.ensureManifestEntry(ctx, spec); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ensureManifestEntry(ctx context.Context, spec domain.ManifestResource) error {
	resource, err := s.r.FindByCode(ctx, spec.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Resource doesn't exist, create it
			actions := fromManifestActions(spec.Actions)
			_, createErr := s.Register(ctx, spec.Code, spec.Service, spec.Description, actions)
			return createErr
		}
		return err
	}
	// Resource exists, update metadata if changed and ensure actions exist.
	resource.Update(spec.Service, spec.Description)
	resource.AttachActions(s.toDomainActions(resource.ID, fromManifestActions(spec.Actions)))
	return s.r.Update(ctx, resource)
}

func (s *Service) toDomainActions(resourceID string, specs []ActionSpec) []domain.Action {
	return toDomainActions(resourceID, specs)
}

// fromManifestActions converts manifest actions to ActionSpecs.
func fromManifestActions(actions []domain.ManifestAction) []ActionSpec {
	specs := make([]ActionSpec, len(actions))
	for i, action := range actions {
		specs[i] = ActionSpec{
			Name:        action.Name,
			Description: action.Description,
		}
	}
	return specs
}

// toDomainActions converts ActionSpecs to domain Actions.
func toDomainActions(resourceID string, specs []ActionSpec) []domain.Action {
	actions := make([]domain.Action, 0, len(specs))
	for _, spec := range specs {
		action, err := domain.NewAction(resourceID, spec.Name, spec.Description)
		if err != nil {
			logger.Warn("Failed to create action from spec", "resource_id", resourceID, "name", spec.Name, "error", err)
			continue
		}
		actions = append(actions, action)
	}
	return actions
}

// RegisterResource allows modules to dynamically register a resource during initialization.
// This method should be called when a module creates its resources, typically during application startup
// or module initialization. If the resource already exists (e.g., from the manifest file), it will be
// updated to ensure actions are current. If it doesn't exist, it will be created.
//
// Example usage in a module:
//
//	err := resourceService.RegisterResource(ctx, "my-module", "my-service", "My module resource", []resource.ActionSpec{
//	    {Name: "view", Description: "View my module data"},
//	    {Name: "manage", Description: "Manage my module data"},
//	})
func (s *Service) RegisterResource(ctx context.Context, code, service, description string, actions []ActionSpec) error {
	// Check if resource already exists
	existing, err := s.r.FindByCode(ctx, code)
	if err == nil && existing != nil {
		// Resource exists, update it to ensure actions are current
		_, updateErr := s.Update(ctx, code, service, description, actions)
		return updateErr
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// Resource doesn't exist, create it
	_, err = s.Register(ctx, code, service, description, actions)
	return err
}
