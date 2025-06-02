package test_container

import (
	"database/sql"

	"b2b.nati011.github.com/config"
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor/db"
	"b2b.nati011.github.com/internal/core/application/user"
	"b2b.nati011.github.com/internal/core/domain/distributor"
	distributorApproval "b2b.nati011.github.com/internal/core/domain/distributor_approval"
)

type TestContainer struct {
	UserService                user.Provider
	DistributorService         distributor.Provider
	DistributorApprovalService distributorApproval.Provider
}

func NewDBIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.UserService = user.NewTestContainer().UserService
	container.DistributorService = distributor.NewDistributorService(container.UserService, db_adapter.NewPostgres(db,
		config.DefaultPaginationBuilder().Build()), container.DistributorApprovalService)

	return container
}

func (t *TestContainer) Teardown(db *sql.DB) {
	t.UserService = user.NewTestContainer().UserService
	t.DistributorService = distributor.NewDistributorService(t.UserService,
		db_adapter.NewPostgres(db, config.DefaultPaginationBuilder().Build()), t.DistributorApprovalService)
}
