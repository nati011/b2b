package distributor

import (
	db_adapter "b2b.nati011.github.com/internal/adapter/secondary/domain/distributor/db"
	"b2b.nati011.github.com/internal/core/application/event"
	user "b2b.nati011.github.com/internal/core/application/user"
	distributorApproval "b2b.nati011.github.com/internal/core/domain/distributor_approval"
)

type TestContainer struct {
	UserService                user.Provider
	DistributorService         Provider
	DistributorApprovalService distributorApproval.Provider
	Event                      *event.Broker
}

func NewPackageIntegrationTestContainer() TestContainer {
	container := TestContainer{}
	container.Event = event.NewBroker()
	container.UserService = user.NewTestContainer().UserService
	container.DistributorApprovalService = distributorApproval.NewTestContainer().DistributorApprovalService
	container.DistributorService = NewDistributorService(container.UserService, db_adapter.NewMock(), container.DistributorApprovalService, container.Event)

	return container
}

func (t *TestContainer) Teardown() {
	t.UserService = user.NewTestContainer().UserService
	t.DistributorService = NewDistributorService(t.UserService, db_adapter.NewMock(), t.DistributorApprovalService, t.Event)
}
