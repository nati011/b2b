package permission

import (
	"context"
	"marketplace/internal/infra/authz/permission/domain"
	resourceDomain "marketplace/internal/infra/authz/resource/domain"
	"marketplace/pkg/logger"
	"marketplace/pkg/pagination"
	"errors"
)

var (
	// ErrPermissionNotFound indicates a permission lookup failure.
	ErrPermissionNotFound = errors.New("permission not found")
	// ErrPermissionAlreadyExists indicates duplicate IDs or resource/action combos.
	ErrPermissionAlreadyExists = errors.New("permission already exists")
)

type permissionRepository interface {
	Create(ctx context.Context, permission *domain.Permission) error
	Update(ctx context.Context, permission *domain.Permission) error
	FindByID(ctx context.Context, id string) (*domain.Permission, error)
	FindByResourceAndAction(ctx context.Context, resourceCode, action string) (*domain.Permission, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Permission], error)
	Exists(ctx context.Context, id string) (bool, error)
}

// PermissionService coordinates permission business logic.
type resourceCatalog interface {
	ResolveAction(ctx context.Context, code, action string) (*resourceDomain.Resource, *resourceDomain.Action, error)
}

// PermissionService coordinates permission business logic.
type PermissionService struct {
	repository permissionRepository
	resources  resourceCatalog
}

func NewPermissionService(repository permissionRepository, resources resourceCatalog) *PermissionService {
	return &PermissionService{
		repository: repository,
		resources:  resources,
	}
}

// Create registers a new permission ensuring uniqueness.
func (s *PermissionService) Create(ctx context.Context, id, resourceCode, action, description string) (*domain.Permission, error) {
	exists, err := s.repository.Exists(ctx, id)
	if err != nil {
		logger.Error("Permission creation failed: existence check error", "permission_id", id, "resource_code", resourceCode, "action", action, "error", err)
		return nil, err
	}
	if exists {
		logger.Warn("Permission creation failed: permission already exists", "permission_id", id, "resource_code", resourceCode, "action", action)
		return nil, ErrPermissionAlreadyExists
	}

	res, actionModel, err := s.resources.ResolveAction(ctx, resourceCode, action)
	if err != nil {
		logger.Error("Permission creation failed: resource/action resolution error", "permission_id", id, "resource_code", resourceCode, "action", action, "error", err)
		return nil, err
	}

	permission, err := domain.NewPermissionEntity(id, toResourceRef(res), actionModel.Name, description)
	if err != nil {
		logger.Warn("Permission creation failed: validation error", "permission_id", id, "resource_code", resourceCode, "action", action, "error", err)
		return nil, err
	}
	if err := s.repository.Create(ctx, permission); err != nil {
		logger.Error("Permission creation failed: repository error", "permission_id", id, "resource_code", resourceCode, "action", action, "error", err)
		return nil, err
	}

	logger.Info("Permission created successfully", "permission_id", id, "resource_code", resourceCode, "action", action)
	return permission, nil
}

// Get fetches a permission by identifier.
func (s *PermissionService) Get(ctx context.Context, id string) (*domain.Permission, error) {
	permission, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			logger.Debug("Permission retrieval failed: permission not found", "permission_id", id)
		} else {
			logger.Error("Permission retrieval failed: repository error", "permission_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("Permission retrieved successfully", "permission_id", id)
	return permission, nil
}

// Update modifies an existing permission's fields.
func (s *PermissionService) Update(ctx context.Context, id, resourceCode, action, description string) (*domain.Permission, error) {
	permission, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("Permission update failed: permission not found", "permission_id", id)
		return nil, err
	}

	res, actionModel, err := s.resources.ResolveAction(ctx, resourceCode, action)
	if err != nil {
		logger.Error("Permission update failed: resource/action resolution error", "permission_id", id, "resource_code", resourceCode, "action", action, "error", err)
		return nil, err
	}

	if err := permission.Update(toResourceRef(res), actionModel.Name, description); err != nil {
		logger.Warn("Permission update failed: validation error", "permission_id", id, "resource_code", resourceCode, "action", action, "error", err)
		return nil, err
	}
	if err := s.repository.Update(ctx, permission); err != nil {
		logger.Error("Permission update failed: repository error", "permission_id", id, "resource_code", resourceCode, "action", action, "error", err)
		return nil, err
	}

	logger.Info("Permission updated successfully", "permission_id", id, "resource_code", resourceCode, "action", action)
	return permission, nil
}

// Delete removes a permission.
func (s *PermissionService) Delete(ctx context.Context, id string) error {
	_, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("Permission deletion failed: permission not found", "permission_id", id)
		return err
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		logger.Error("Permission deletion failed: repository error", "permission_id", id, "error", err)
		return err
	}

	logger.Info("Permission deleted successfully", "permission_id", id)
	return nil
}

// List retrieves a paginated list of permissions.
func (s *PermissionService) List(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Permission], error) {
	result, err := s.repository.FindAll(ctx, pageReq)
	if err != nil {
		logger.Error("Permission pagination failed: repository error", "page", pageReq.Page, "limit", pageReq.Limit, "error", err)
		return pagination.PageResult[*domain.Permission]{}, err
	}
	logger.Debug("Permission pagination completed", "page", result.Page, "total", result.Total, "items", len(result.Items))
	return result, nil
}

// GetByResourceAndAction fetches a permission by resource code and action.
func (s *PermissionService) GetByResourceAndAction(ctx context.Context, resourceCode, action string) (*domain.Permission, error) {
	permission, err := s.repository.FindByResourceAndAction(ctx, resourceCode, action)
	if err != nil {
		// Check if the error is ErrPermissionNotFound (from repository or service)
		// Since repository and service are in different packages but both define ErrPermissionNotFound,
		// we need to check the error message as well
		if errors.Is(err, ErrPermissionNotFound) || err.Error() == "permission not found" {
			logger.Debug("Permission retrieval failed: permission not found", "resource_code", resourceCode, "action", action)
			return nil, ErrPermissionNotFound
		}
		logger.Error("Permission retrieval failed: repository error", "resource_code", resourceCode, "action", action, "error", err)
		return nil, err
	}
	logger.Debug("Permission retrieved successfully", "resource_code", resourceCode, "action", action)
	return permission, nil
}

func toResourceRef(res *resourceDomain.Resource) domain.ResourceRef {
	return domain.ResourceRef{
		ID:          res.ID,
		Code:        res.Code,
		Service:     res.Service,
		Description: res.Description,
	}
}
