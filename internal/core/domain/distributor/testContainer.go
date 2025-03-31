package distributor

import (
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/distributor/db"
	user "b2b.nati011.github.com/internal/core/application/user"
)

type TestContainer struct {
	UserService        user.Provider
	DistributorService Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.UserService = user.NewTestContainer().UserService
	container.DistributorService = NewDistributorService(container.UserService, db_adapter.NewMock())

	return container
}

func (t *TestContainer) Cleanup() {
	t.UserService = user.NewTestContainer().UserService
	t.DistributorService = NewDistributorService(t.UserService, db_adapter.NewMock())
}
