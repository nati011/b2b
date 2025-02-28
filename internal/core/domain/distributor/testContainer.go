package distributor

import (
	user "b2b.nati011.github.com/internal/core/application/service/user"
)

type TestContainer struct {
	UserService        user.Provider
	DistributorService Provider
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	user_testContainer := user.NewTestContainer()
	container.UserService = user_testContainer.UserService

	return container
}
