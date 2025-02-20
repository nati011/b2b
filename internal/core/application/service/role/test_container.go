package role

import (
	db_resource_mock "b2b.nati011.github.com/internal/adapter/secondary/resource/db"
	db_role_mock "b2b.nati011.github.com/internal/adapter/secondary/role/db"

	resource "b2b.nati011.github.com/internal/core/application/service/resource"
)

type TestContainer struct {
	RoleService Provider

	resourceService resource.Provider
}

func NewTestContainer() TestContainer {
	c := TestContainer{}
	c.resourceService = initResourceService()
	c.RoleService = NewRole(
		db_role_mock.NewMock(),
		c.resourceService,
	)
	return c
}

func initResourceService() resource.Provider {
	return resource.NewResource(
		db_resource_mock.NewMock(),
	)
}
