package role

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	db_resource_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	db_role_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"

	resource "b2b.nati011.github.com/internal/core/application/resource"
	"b2b.nati011.github.com/internal/core/application/role"
)

type TestContainer struct {
	RoleService role.Provider

	ResourceService resource.Provider
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
	c := TestContainer{}
	c.ResourceService = resource.NewResource(
		db_resource_adapter.NewPostgres(db, &config.Pagination{
			Limit:  10,
			Offset: 0,
		}),
	)
	c.RoleService = role.NewRole(
		db_role_adapter.NewPostgres(db, &config.Pagination{
			Limit:  10,
			Offset: 0,
		}),
		c.ResourceService,
	)
	return c
}

func (t *TestContainer) TeardownIntegrationTestContainer(db *sql.DB) {
}
