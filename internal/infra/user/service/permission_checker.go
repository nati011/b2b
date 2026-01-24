package service

import (
	"context"
)

// PermissionChecker defines the interface for checking user permissions.
// This interface decouples the user module from authz modules.
type PermissionChecker interface {
	// HasPermission checks if a user has a specific permission.
	// It takes userID, resource, and action as parameters.
	HasPermission(ctx context.Context, userID, resource, action string) (bool, error)
}

