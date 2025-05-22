package authorization

import (
	"context"

	role "b2b.nati011.github.com/internal/core/application/role"
)

type ResourceAccess struct {
	Roles []string `json:"roles"`
}

type Provider interface {
	IsAuthorizedForResource(ctx context.Context, role string) (bool, error)
	GetUserAuthorization(ctx context.Context) (ResourceAccess, error)
}

type Authorization struct {
	RoleService role.Provider
}

func NewAuthorization() Provider {
	return &Authorization{}
}

func (a *Authorization) IsAuthorizedForResource(ctx context.Context, role string) (bool, error) {
	return false, nil
}

func (a *Authorization) GetUserAuthorization(ctx context.Context) (ResourceAccess, error) {
	return ResourceAccess{}, nil
}
