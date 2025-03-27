package retailer

import (
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/application/retailer/db"
	user "b2b.nati011.github.com/internal/core/application/user"
)

type TestContainer struct {
	UserService     user.Provider
	RetailerService Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.UserService = user.NewTestContainer().UserService
	container.RetailerService = NewRetailerService(container.UserService, db_adapter.NewMock())

	return container
}
