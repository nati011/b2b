package user

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	db_resource_mock "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	db_role_mock "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	db_provider "b2b.nati011.github.com/internal/adapter/secondary/application/user/db"
	auth "b2b.nati011.github.com/internal/core/application/auth"
	auth_test "b2b.nati011.github.com/internal/core/application/auth/test"
	resource "b2b.nati011.github.com/internal/core/application/resource"
	role "b2b.nati011.github.com/internal/core/application/role"
)

var db_global *sql.DB

type TestContainer struct {
	UserService     Provider
	RoleService     role.Provider
	ResourceService resource.Provider
	AuthService     auth.Provider
}

func NewTestContainer() TestContainer {
	c := TestContainer{}

	c.ResourceService = resource.NewTestContainer().ResourceService
	c.RoleService = role.NewRole(db_role_mock.NewMock(), c.ResourceService)
	c.AuthService = auth.NewTestContainer().Service
	c.UserService = NewUser(
		db_provider.NewMock(),
		c.RoleService,
		c.AuthService,
	)
	return c
}

func (t *TestContainer) Teardown() {
	t.AuthService = auth.NewTestContainer().Service
	t.UserService = NewUser(
		db_provider.NewMock(),
		t.RoleService,
		t.AuthService,
	)
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
	db_global = db
	c := TestContainer{}
	c.RoleService = role.NewRole(
		db_role_mock.NewPostgres(db, config.DefaultPaginationBuilder().Build()),
		resource.NewResource(
			db_resource_mock.NewPostgres(db, config.DefaultPaginationBuilder().Build()),
		),
	)
	c.AuthService = auth_test.NewIntegrationTestContainer().Service
	c.UserService = NewUser(
		db_provider.NewPostgres(
			db,
			config.DefaultPaginationBuilder().Build(),
		),
		c.RoleService,
		c.AuthService,
	)
	return c
}

func (t *TestContainer) TeardownIntegrationTestContainer() {
	t.AuthService = auth_test.NewIntegrationTestContainer().Service
	t.UserService = NewUser(
		db_provider.NewPostgres(
			db_global,
			config.DefaultPaginationBuilder().Build(),
		),
		t.RoleService,
		t.AuthService,
	)
}
