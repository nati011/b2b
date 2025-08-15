package role

import (
	db_resource_mock "b2b.nati011.github.com/internal/adapter/secondary/application/resource/db"
	db_role_mock "b2b.nati011.github.com/internal/adapter/secondary/application/role/db"

	resource "b2b.nati011.github.com/internal/core/application/resource"
)

type TestContainer struct {
	RoleService     Provider
	ResourceService resource.Provider
}

func NewTestContainer() TestContainer {
	c := TestContainer{}
	c.ResourceService = resource.NewResource(
		db_resource_mock.NewMock(),
	)
	c.RoleService = NewRole(
		db_role_mock.NewMock(),
		c.ResourceService,
	)
	return c
}

func (c *TestContainer) Teardown() {
	c.ResourceService = resource.NewResource(
		db_resource_mock.NewMock(),
	)
	c.RoleService = NewRole(
		db_role_mock.NewMock(),
		c.ResourceService,
	)
}
