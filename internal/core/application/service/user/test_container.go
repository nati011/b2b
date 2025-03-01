package user

import (
	"database/sql"

	db_resource_mock "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	db_role_mock "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	db_user_mock "b2b.nati011.github.com/internal/adapter/secondary/application/user/db"
	resource "b2b.nati011.github.com/internal/core/application/service/resource"
	role "b2b.nati011.github.com/internal/core/application/service/role"
)

type TestContainer struct {
	UserService Provider

	RoleService role.Provider
}

func NewTestContainer() TestContainer {
	c := TestContainer{}

	role_testContainer := role.NewTestContainer()
	c.RoleService = role_testContainer.RoleService
	c.UserService = NewUser(
		db_user_mock.NewMock(),
		c.RoleService,
	)
	return c
}

func initRoleService() role.Provider {
	return role.NewTestContainer().RoleService
}

func initResourceService() resource.Provider {
	return resource.NewResource(db_resource_mock.NewMock())
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
	c := TestContainer{}
	c.UserService = NewUser(
		db_user_mock.NewPostgres(
			db,
		),
		c.RoleService,
	)
	c.RoleService = role.NewRole(
		db_role_mock.NewPostgres(
			db,
		),
		resource.NewResource(
			db_resource_mock.NewPostgres(
				db,
			),
		),
	)
	return c
}
