package service

import (
	"context"
)

// UserPermissionChecker implements PermissionChecker using role lookup
type UserPermissionChecker struct {
	userRepo              Repository
	rolePermissionChecker RolePermissionChecker
}

// NewUserPermissionChecker creates a new permission checker
func NewUserPermissionChecker(userRepo Repository, rolePermissionChecker RolePermissionChecker) *UserPermissionChecker {
	return &UserPermissionChecker{
		userRepo:              userRepo,
		rolePermissionChecker: rolePermissionChecker,
	}
}

// HasPermission checks if a user has a specific permission by:
// 1. Loading the user's role IDs
// 2. For each role, checking if it has the permission using the role permission checker
func (c *UserPermissionChecker) HasPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	user, err := c.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check each role for the permission
	for _, roleID := range user.RoleIDs {
		hasPerm, err := c.rolePermissionChecker.HasPermission(ctx, roleID, resource, action)
		if err != nil {
			// If role not found or error, skip it (might have been deleted)
			continue
		}
		if hasPerm {
			return true, nil
		}
	}

	return false, nil
}
