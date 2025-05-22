package authorization

import (
	"b2b.nati011.github.com/internal/core/application/resource"
	"b2b.nati011.github.com/internal/core/application/role"
	"b2b.nati011.github.com/internal/core/application/user"

	db_resource_provider "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	db_role_mock "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	db_user_provider "b2b.nati011.github.com/internal/adapter/secondary/application/user/db"
	auth "b2b.nati011.github.com/internal/core/application/authentication"
)

type TestContainer struct {
	UserService          user.Provider
	RoleService          role.Provider
	ResourceService      resource.Provider
	AuthorizationService Provider
}

func NewTestContainer() TestContainer {
	container := TestContainer{}
	container.RoleService = role.NewRole(db_role_mock.NewMock(), container.ResourceService)
	container.ResourceService = resource.NewResource(db_resource_provider.NewMock())
	container.UserService = user.NewUser(
		db_user_provider.NewMock(),
		container.RoleService,
		auth.NewIntegrationAuthContainer(),
	)
	container.AuthorizationService = NewAuthorization(container.RoleService, container.UserService)
	return container
}

func (t *TestContainer) teardown() TestContainer {
	t.AuthorizationService = NewAuthorization(t.RoleService, t.UserService)
	return *t
}
