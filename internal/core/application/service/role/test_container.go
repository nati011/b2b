package role

import (
	"database/sql"

	db_resource_mock "b2b.nati011.github.com/internal/adapter/secondary/resource/db"
	db_role_mock "b2b.nati011.github.com/internal/adapter/secondary/role/db"

	resource "b2b.nati011.github.com/internal/core/application/service/resource"
)

type TestContainer struct {
	RoleService Provider

	ResourceService resource.Provider
}

func NewTestContainer() TestContainer {
	c := TestContainer{}
	c.ResourceService = initResourceService()
	c.RoleService = NewRole(
		db_role_mock.NewMock(),
		c.ResourceService,
	)
	return c
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
	c := TestContainer{}
	c.ResourceService = resource.NewResource(
		db_resource_mock.NewPostgres(
			db,
		),
	)
	c.RoleService = NewRole(
		db_role_mock.NewPostgres(
			db,
		),
		c.ResourceService,
	)
	return c
}

func initResourceService() resource.Provider {
	return resource.NewResource(
		db_resource_mock.NewMock(),
	)
}

//update user info
//remove
