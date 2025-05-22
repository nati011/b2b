package user

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	db_resource_mock "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	db_role_mock "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"
	db_provider "b2b.nati011.github.com/internal/adapter/secondary/application/user/db"
	auth "b2b.nati011.github.com/internal/core/application/authentication"
	resource "b2b.nati011.github.com/internal/core/application/resource"
	role "b2b.nati011.github.com/internal/core/application/role"
)

var db_global *sql.DB

type TestContainer struct {
	UserService Provider
	RoleService role.Provider
	AuthService auth.Provider
}

func NewTestContainer() TestContainer {
	c := TestContainer{}

	role_testContainer := role.NewTestContainer()
	c.RoleService = role_testContainer.RoleService
	c.AuthService = auth.NewIntegrationAuthContainer()
	c.UserService = NewUser(
		db_provider.NewMock(),
		c.RoleService,
		c.AuthService,
	)
	return c
}

func (t *TestContainer) Teardown() {
	role_testContainer := role.NewTestContainer()
	t.RoleService = role_testContainer.RoleService
	t.AuthService = auth.NewIntegrationAuthContainer()
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
		db_role_mock.NewPostgres(db, &config.Pagination{
			Limit:  10,
			Offset: 0,
		}),
		resource.NewResource(
			db_resource_mock.NewPostgres(db, &config.Pagination{
				Limit:  10,
				Offset: 0,
			}),
		),
	)
	c.AuthService = auth.NewIntegrationAuthContainer()
	c.UserService = NewUser(
		db_provider.NewPostgres(
			db,
			&config.Pagination{
				Limit:  10,
				Offset: 0,
			},
		),
		c.RoleService,
		c.AuthService,
	)
	return c
}

func (t *TestContainer) TeardownIntegrationTestContainer() {
	t.AuthService = auth.NewIntegrationAuthContainer()
	t.UserService = NewUser(
		db_provider.NewPostgres(
			db_global,
			&config.Pagination{
				Limit:  10,
				Offset: 0,
			},
		),
		t.RoleService,
		t.AuthService,
	)
}
