package test_container

import (
	"database/sql"

	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor/db"
	"b2b.nati011.github.com/internal/core/application/user"
	"b2b.nati011.github.com/internal/core/domain/distributor"
)

type TestContainer struct {
	UserService        user.Provider
	DistributorService distributor.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.UserService = user.NewTestContainer().UserService
	container.DistributorService = distributor.NewDistributorService(container.UserService, db_adapter.NewPostgres(db))

	return container
}

func (t *TestContainer) Teardown(db *sql.DB) {
	t.UserService = user.NewTestContainer().UserService
	t.DistributorService = distributor.NewDistributorService(t.UserService, db_adapter.NewPostgres(db))
}
