package distributorApproval

import (
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor_approval/db"
)

type TestContainer struct {
	DistributorApprovalService Provider
}

func NewTestContainer() TestContainer {
	container := TestContainer{}
	container.DistributorApprovalService = NewDistributorApprovalService(db_adapter.NewMock())

	return container
}

func (t *TestContainer) Teardown() {
	t.DistributorApprovalService = NewDistributorApprovalService(db_adapter.NewMock())
}
