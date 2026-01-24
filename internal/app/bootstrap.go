package app

import (
	"context"
	"errors"
	"os"

	"marketplace/internal/config"
	basicauthservice "marketplace/internal/infra/auth/basic/service"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	resourcedomain "marketplace/internal/infra/authz/resource/domain"
	resourceservice "marketplace/internal/infra/authz/resource/service"
	roleservice "marketplace/internal/infra/authz/role/service"
	userservice "marketplace/internal/infra/user/service"

	"github.com/google/uuid"
)

// Bootstrap seeds resources and roles using the provided configuration.
// It is safe to call multiple times; existing entries are left untouched.
func Bootstrap(
	ctx context.Context,
	cfg *config.Config,
	userService *userservice.Service,
	roleService *roleservice.RoleService,
	permService *permissionservice.PermissionService,
	resourceService *resourceservice.Service,
	basicAuthService *basicauthservice.Service,
) error {
	_ = userService
	_ = basicAuthService

	if cfg == nil {
		return nil
	}

	if resourceService != nil && manifestExists(cfg.Resources.Path) {
		manifest, err := resourcedomain.LoadResources(cfg.Resources.Path)
		if err != nil {
			return err
		}
		if err := resourceService.Bootstrap(ctx, manifest); err != nil {
			return err
		}
	}

	if roleService != nil && permService != nil && rolesFileExists(cfg) {
		roles, err := config.LoadRolesFromFile(cfg.Roles.Path)
		if err != nil {
			return err
		}

		for _, role := range roles {
			permissionIDs, err := resolveRolePermissions(ctx, permService, role.Permissions)
			if err != nil {
				return err
			}
			if _, err := roleService.Create(ctx, role.ID, role.Name, role.Description, permissionIDs); err != nil {
				if errors.Is(err, roleservice.ErrRoleAlreadyExists) || errors.Is(err, roleservice.ErrRoleNameConflict) {
					continue
				}
				return err
			}
		}
	}

	return nil
}

func resolveRolePermissions(ctx context.Context, permService *permissionservice.PermissionService, patterns []config.PermissionPattern) ([]string, error) {
	permissionIDs := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		for _, action := range pattern.Actions {
			perm, err := permService.GetByResourceAndAction(ctx, pattern.Resource, action)
			if err == nil {
				permissionIDs = append(permissionIDs, perm.ID)
				continue
			}

			newID := uuid.NewString()
			created, createErr := permService.Create(ctx, newID, pattern.Resource, action, pattern.Resource+"."+action)
			if createErr != nil {
				if errors.Is(createErr, permissionservice.ErrPermissionAlreadyExists) {
					perm, lookupErr := permService.GetByResourceAndAction(ctx, pattern.Resource, action)
					if lookupErr != nil {
						return nil, lookupErr
					}
					permissionIDs = append(permissionIDs, perm.ID)
					continue
				}
				return nil, createErr
			}
			permissionIDs = append(permissionIDs, created.ID)
		}
	}
	return permissionIDs, nil
}

func manifestExists(path string) bool {
	if path == "" {
		path = "config/resources.yaml"
	}
	_, err := os.Stat(path)
	return err == nil
}

func rolesFileExists(cfg *config.Config) bool {
	if cfg.Roles == nil {
		return false
	}
	path := cfg.Roles.Path
	if path == "" {
		path = "config/roles.yaml"
	}
	_, err := os.Stat(path)
	return err == nil
}
