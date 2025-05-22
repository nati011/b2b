package authorization

import (
	"context"

	role "b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/user"
)

type ResourceAccess struct {
	Roles []string `json:"roles"`
}

type Provider interface {
	IsAuthorizedForResource(ctx context.Context, userId int, resourceId int) (bool, error)
	GetUserAuthorization(ctx context.Context, userId int) (ResourceAccess, error)
}

type Authorization struct {
	RoleService role.Provider
	UserService user.Provider
}

func NewAuthorization(rs role.Provider, u user.Provider) Provider {
	return &Authorization{
		RoleService: rs,
		UserService: u,
	}
}

func (a *Authorization) IsAuthorizedForResource(ctx context.Context, userId int, resourceId int) (bool, error) {
	return false, nil
}

func (a *Authorization) GetUserAuthorization(ctx context.Context, userId int) (ResourceAccess, error) {
	return ResourceAccess{}, nil
}
