package app

import (
	"context"
	"errors"
	"fmt"
	"os"

	"marketplace/internal/config"
	basicauthservice "marketplace/internal/infra/auth/basic/service"
	permissionservice "marketplace/internal/infra/authz/permission/service"
	resourcedomain "marketplace/internal/infra/authz/resource/domain"
	resourceservice "marketplace/internal/infra/authz/resource/service"
	roleservice "marketplace/internal/infra/authz/role/service"
	userdomain "marketplace/internal/infra/user/domain"
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
			// Generate deterministic UUID if role.ID is not a valid UUID
			roleID := role.ID
			if _, err := uuid.Parse(roleID); err != nil {
				// Generate a deterministic UUID from the role name using SHA1
				// This ensures the same role name always gets the same UUID
				namespace := uuid.NameSpaceDNS
				roleID = uuid.NewSHA1(namespace, []byte(role.Name)).String()
			}
			if _, err := roleService.Create(ctx, roleID, role.Name, role.Description, permissionIDs); err != nil {
				if errors.Is(err, roleservice.ErrRoleAlreadyExists) || errors.Is(err, roleservice.ErrRoleNameConflict) {
					continue
				}
				return err
			}
		}
	}

	// Bootstrap superadmin users from config
	if cfg.Auth.Basic != nil && userService != nil && basicAuthService != nil {
		if err := bootstrapSuperadminUsers(ctx, cfg, userService, roleService, basicAuthService); err != nil {
			return err
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

// bootstrapSuperadminUsers creates superadmin users from config if they don't exist
func bootstrapSuperadminUsers(
	ctx context.Context,
	cfg *config.Config,
	userService *userservice.Service,
	roleService *roleservice.RoleService,
	basicAuthService *basicauthservice.Service,
) error {
	basicCfg := cfg.Auth.Basic
	if len(basicCfg.Users) == 0 {
		return nil
	}

	// Get admin role (prefer "admin" role, fallback to "Super Administrator" for backward compatibility)
	adminRole, err := roleService.GetByName(ctx, "Admin")
	if err != nil {
		// Fallback to Super Administrator role if Admin role doesn't exist
		adminRole, err = roleService.GetByName(ctx, "Super Administrator")
		if err != nil {
			// If role doesn't exist yet, skip user creation (roles must be bootstrapped first)
			return nil
		}
	}

	for _, authUser := range basicCfg.Users {
		// Construct email: user_prefix + username + email_domain
		emailStr := basicCfg.UserPrefix + authUser.Username + basicCfg.EmailDomain
		
		// Check if user already exists with the new email
		existingUser, err := userService.FindByEmail(ctx, emailStr)
		if err == nil && existingUser != nil {
			// User exists, ensure they have admin role
			if !existingUser.HasRole(adminRole.ID) {
				if _, err := userService.AssignRole(ctx, existingUser.ID, adminRole.ID); err != nil {
					return fmt.Errorf("failed to assign admin role to existing user %s: %w", emailStr, err)
				}
			}
			// Check if credentials exist
			_, err := basicAuthService.Repository().FindByUsername(ctx, emailStr)
			if err != nil {
				// Credentials don't exist, create them
				if _, err := basicAuthService.CreateCredential(ctx, emailStr, authUser.Password, existingUser.ID); err != nil {
					return fmt.Errorf("failed to create credentials for user %s: %w", emailStr, err)
				}
			}
			continue
		}

		// Also check for old email format (backward compatibility)
		// If username was "superadmin", check for "admin-superadmin@admin.local"
		oldEmailStr := basicCfg.UserPrefix + "superadmin" + basicCfg.EmailDomain
		if oldEmailStr != emailStr {
			oldUser, err := userService.FindByEmail(ctx, oldEmailStr)
			if err == nil && oldUser != nil {
				// Old user exists, ensure they have admin role
				if !oldUser.HasRole(adminRole.ID) {
					if _, err := userService.AssignRole(ctx, oldUser.ID, adminRole.ID); err != nil {
						return fmt.Errorf("failed to assign admin role to existing user %s: %w", oldEmailStr, err)
					}
				}
				// Update credentials if they exist with old email
				_, err := basicAuthService.Repository().FindByUsername(ctx, oldEmailStr)
				if err == nil {
					// Credentials exist with old email, that's fine - user can still login
					continue
				}
			}
		}

		// User doesn't exist, create it
		email, err := userdomain.NewEmail(emailStr)
		if err != nil {
			return fmt.Errorf("invalid email %s: %w", emailStr, err)
		}

		phoneNumber := userdomain.PhoneNumber{} // Empty phone number
		externalID := fmt.Sprintf("EXT-ADMIN-%s", authUser.Username)
		
		result, err := userService.Create(ctx, externalID, email, phoneNumber, basicCfg.UserName, userdomain.UserTypeOfficer)
		if err != nil {
			if errors.Is(err, userservice.ErrEmailAlreadyExists) {
				// User was created between check and create, skip
				continue
			}
			return fmt.Errorf("failed to create superadmin user %s: %w", emailStr, err)
		}

		// Assign admin role
		if _, err := userService.AssignRole(ctx, result.User.ID, adminRole.ID); err != nil {
			return fmt.Errorf("failed to assign admin role to user %s: %w", emailStr, err)
		}

		// Activate the user
		if result.User.Status != userdomain.UserStatusActive {
			if _, err := userService.Activate(ctx, result.User.ID); err != nil {
				return fmt.Errorf("failed to activate superadmin user %s: %w", emailStr, err)
			}
		}

		// Create credentials
		if _, err := basicAuthService.CreateCredential(ctx, emailStr, authUser.Password, result.User.ID); err != nil {
			return fmt.Errorf("failed to create credentials for user %s: %w", emailStr, err)
		}
	}

	return nil
}
