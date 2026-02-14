package test

import (
	"database/sql"

	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor_approval/db"
	distributorApproval "b2b.nati011.github.com/internal/core/domain/distributor_approval"
)

type TestContainer struct {
	DistributorApprovalService distributorApproval.Provider
}

func NewIntegrationTestContainer(db *sql.DB) TestContainer {
	container := TestContainer{}
	container.DistributorApprovalService = distributorApproval.NewDistributorApprovalService(db_adapter.NewPostgres(db))

	return container
}

func (t *TestContainer) Teardown(db *sql.DB) {
	t.DistributorApprovalService = distributorApproval.NewDistributorApprovalService(db_adapter.NewPostgres(db))
}
