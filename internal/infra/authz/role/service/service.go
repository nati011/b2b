package role

import (
	"context"
	permissionDomain "marketplace/internal/infra/authz/permission/domain"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	"marketplace/internal/infra/authz/role/domain"
	"marketplace/pkg/logger"
	"marketplace/pkg/pagination"
	"errors"
)

var (
	// ErrRoleNotFound indicates that a role was not found in the repository.
	ErrRoleNotFound = errors.New("role not found")

	// ErrRoleAlreadyExists indicates that a role with the same identifier already exists.
	ErrRoleAlreadyExists = errors.New("role already exists")

	// ErrRoleNameConflict indicates a duplicate role name.
	ErrRoleNameConflict = errors.New("role name already exists")

	// ErrInvalidPermissionIDs indicates that at least one provided permission ID is invalid.
	ErrInvalidPermissionIDs = errors.New("invalid permission identifiers supplied")
)

type roleRepository interface {
	Create(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, role *domain.Role) error
	FindByID(ctx context.Context, id string) (*domain.Role, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Role], error)
	Exists(ctx context.Context, id string) (bool, error)
	FindByName(ctx context.Context, name string) (*domain.Role, error)
}

type permissionLookup interface {
	FindByID(ctx context.Context, id string) (*permissionDomain.Permission, error)
}

type RoleService struct {
	repository  roleRepository
	permissions permissionLookup
}

func NewRoleService(repository roleRepository, permissions permissionLookup) *RoleService {
	return &RoleService{
		repository:  repository,
		permissions: permissions,
	}
}

// Create registers a new role after ensuring uniqueness.
func (s *RoleService) Create(ctx context.Context, id, name, description string, permissionIDs []string) (*domain.Role, error) {
	exists, err := s.repository.Exists(ctx, id)
	if err != nil {
		logger.Error("Role creation failed: existence check error", "role_id", id, "name", name, "error", err)
		return nil, err
	}
	if exists {
		logger.Warn("Role creation failed: role already exists", "role_id", id, "name", name)
		return nil, ErrRoleAlreadyExists
	}

	if err := s.ensureNameUnique(ctx, name, ""); err != nil {
		if errors.Is(err, ErrRoleNameConflict) {
			logger.Warn("Role creation failed: name conflict", "role_id", id, "name", name)
		} else {
			logger.Error("Role creation failed: name uniqueness check error", "role_id", id, "name", name, "error", err)
		}
		return nil, err
	}

	// Validate permission IDs exist
	if err := s.validatePermissionIDs(ctx, permissionIDs); err != nil {
		if errors.Is(err, ErrInvalidPermissionIDs) {
			logger.Warn("Role creation failed: invalid permission IDs", "role_id", id, "name", name, "permission_ids", permissionIDs)
		} else {
			logger.Error("Role creation failed: permission validation error", "role_id", id, "name", name, "error", err)
		}
		return nil, err
	}

	role, err := domain.NewRole(id, name, description, permissionIDs)
	if err != nil {
		logger.Warn("Role creation failed: validation error", "role_id", id, "name", name, "error", err)
		return nil, err
	}
	if err := s.repository.Create(ctx, role); err != nil {
		logger.Error("Role creation failed: repository error", "role_id", id, "name", name, "error", err)
		return nil, err
	}

	logger.Info("Role created successfully", "role_id", id, "name", name)
	return role, nil
}

// Get fetches a role by identifier.
func (s *RoleService) Get(ctx context.Context, id string) (*domain.Role, error) {
	role, err := s.repository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			logger.Debug("Role retrieval failed: role not found", "role_id", id)
		} else {
			logger.Error("Role retrieval failed: repository error", "role_id", id, "error", err)
		}
		return nil, err
	}
	logger.Debug("Role retrieved successfully", "role_id", id)
	return role, nil
}

// Update modifies an existing role's attributes.
func (s *RoleService) Update(ctx context.Context, id, name, description string, permissionIDs []string) (*domain.Role, error) {
	// Name uniqueness/permission checks are handled at service/repository level,
	// but structural validation of the aggregate lives in the domain.
	if err := domain.ValidateRoleInput(id, name, description); err != nil {
		logger.Warn("Role update failed: validation error", "role_id", id, "name", name, "error", err)
		return nil, err
	}

	role, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("Role update failed: role not found", "role_id", id)
		return nil, err
	}

	if err := s.ensureNameUnique(ctx, name, id); err != nil {
		if errors.Is(err, ErrRoleNameConflict) {
			logger.Warn("Role update failed: name conflict", "role_id", id, "name", name)
		} else {
			logger.Error("Role update failed: name uniqueness check error", "role_id", id, "name", name, "error", err)
		}
		return nil, err
	}

	// Validate permission IDs exist
	if err := s.validatePermissionIDs(ctx, permissionIDs); err != nil {
		if errors.Is(err, ErrInvalidPermissionIDs) {
			logger.Warn("Role update failed: invalid permission IDs", "role_id", id, "name", name, "permission_ids", permissionIDs)
		} else {
			logger.Error("Role update failed: permission validation error", "role_id", id, "name", name, "error", err)
		}
		return nil, err
	}

	if err := role.Update(name, description, permissionIDs); err != nil {
		logger.Warn("Role update failed: validation error", "role_id", id, "name", name, "error", err)
		return nil, err
	}

	if err := s.repository.Update(ctx, role); err != nil {
		logger.Error("Role update failed: repository error", "role_id", id, "name", name, "error", err)
		return nil, err
	}

	logger.Info("Role updated successfully", "role_id", id, "name", name)
	return role, nil
}

// Delete removes an existing role.
func (s *RoleService) Delete(ctx context.Context, id string) error {
	_, err := s.repository.FindByID(ctx, id)
	if err != nil {
		logger.Debug("Role deletion failed: role not found", "role_id", id)
		return err
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		logger.Error("Role deletion failed: repository error", "role_id", id, "error", err)
		return err
	}

	logger.Info("Role deleted successfully", "role_id", id)
	return nil
}

// List retrieves a paginated list of roles.
func (s *RoleService) List(ctx context.Context, pageReq pagination.PageRequest) (pagination.PageResult[*domain.Role], error) {
	result, err := s.repository.FindAll(ctx, pageReq)
	if err != nil {
		logger.Error("Role pagination failed: repository error", "page", pageReq.Page, "limit", pageReq.Limit, "error", err)
		return pagination.PageResult[*domain.Role]{}, err
	}
	logger.Debug("Role pagination completed", "page", result.Page, "total", result.Total, "items", len(result.Items))
	return result, nil
}

// GetByName fetches a role by its name (case-insensitive).
func (s *RoleService) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	role, err := s.repository.FindByName(ctx, name)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			logger.Debug("Role retrieval failed: role not found", "name", name)
		} else {
			logger.Error("Role retrieval failed: repository error", "name", name, "error", err)
		}
		return nil, err
	}
	logger.Debug("Role retrieved successfully", "name", name, "role_id", role.ID)
	return role, nil
}

func (s *RoleService) ensureNameUnique(ctx context.Context, name, excludeID string) error {
	existing, err := s.repository.FindByName(ctx, name)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			return nil
		}
		return err
	}
	if existing.ID != excludeID {
		return ErrRoleNameConflict
	}
	return nil
}

func (s *RoleService) validatePermissionIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	for _, id := range ids {
		_, err := s.permissions.FindByID(ctx, id)
		if err != nil {
			if errors.Is(err, permissionservice.ErrPermissionNotFound) {
				return ErrInvalidPermissionIDs
			}
			return err
		}
	}
	return nil
}

// HasPermission checks if a role has a specific permission by resource and action.
// This method loads the role's permissions and checks if any match the given resource/action.
func (s *RoleService) HasPermission(ctx context.Context, roleID, resource, action string) (bool, error) {
	role, err := s.repository.FindByID(ctx, roleID)
	if err != nil {
		return false, err
	}

	// Load permissions for the role
	for _, permissionID := range role.PermissionIDs {
		perm, err := s.permissions.FindByID(ctx, permissionID)
		if err != nil {
			// If permission not found, skip it (might have been deleted)
			continue
		}
		if perm.Resource.Code == resource && perm.Action == action {
			return true, nil
		}
	}

	return false, nil
}

// GetPermissions returns the full permission objects for a role.
// This is used when we need to return permissions in DTOs.
func (s *RoleService) GetPermissions(ctx context.Context, roleID string) ([]*permissionDomain.Permission, error) {
	role, err := s.repository.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	permissions := make([]*permissionDomain.Permission, 0, len(role.PermissionIDs))
	for _, permissionID := range role.PermissionIDs {
		perm, err := s.permissions.FindByID(ctx, permissionID)
		if err != nil {
			// If permission not found, skip it (might have been deleted)
			continue
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}
